package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/memory/bucket"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/sql/ip"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// DefaultLoginBucketLimit is default value of rate (in minute) for login buckets
	DefaultLoginBucketLimit = 10
	// DefaultPasswordBucketLimit is default value of rate (in minute) for password buckets
	DefaultPasswordBucketLimit = 100
	// DefaultIPBucketLimit is default value of rate (in minute) for IP buckets
	DefaultIPBucketLimit = 1000
	// DefaultBucketActiveTimeout is default value of bucket active timeout
	DefaultBucketActiveTimeout = 2 * time.Minute
)

// LimitsConfig is set of limits of all types of buckets
type LimitsConfig struct {
	LoginLimit    uint
	PasswordLimit uint
	IPLimit       uint
}

// NewLimitsConfigByViper constructs new LimitConfig by viper (get from "limits" section of app config)
func NewLimitsConfigByViper(v *viper.Viper) LimitsConfig {
	limits := v.GetStringMapString("limits")

	return LimitsConfig{
		LoginLimit:    getUintFromStringMap(limits, "login", DefaultLoginBucketLimit),
		PasswordLimit: getUintFromStringMap(limits, "password", DefaultPasswordBucketLimit),
		IPLimit:       getUintFromStringMap(limits, "ip", DefaultIPBucketLimit),
	}
}

// StorageSet struct represents set of buckets storages
type StorageSet struct {
	LoginStorage    entities.BucketStorage
	PasswordStorage entities.BucketStorage
	IPStorage       entities.BucketStorage
}

// ListSet struct represents set of IP lists (black and white)
type ListSet struct {
	BlackList entities.IPList
	WhiteList entities.IPList
}

// API is GRPC API server
type API struct {
	LimitsConfig
	StorageSet
	ListSet
	bucketActiveTimeout time.Duration
	logger              *zap.SugaredLogger
	nowTimeFn           func() time.Time
}

// NewAPIByViper constructs new API data type by viper app config and connection to DB
func NewAPIByViper(v *viper.Viper, db *sqlx.DB, logger *zap.SugaredLogger) *API {
	timeouts := v.GetStringMapString("timeouts")
	bucketActiveTimeout := getDurationFromStringMap(timeouts, "bucket_active", DefaultBucketActiveTimeout)

	api := &API{
		LimitsConfig: NewLimitsConfigByViper(v),
		StorageSet: StorageSet{
			LoginStorage:    bucket.NewStorage(),
			PasswordStorage: bucket.NewStorage(),
			IPStorage:       bucket.NewStorage(),
		},
		ListSet: ListSet{
			BlackList: ip.NewList(db, "black"),
			WhiteList: ip.NewList(db, "white"),
		},
		bucketActiveTimeout: bucketActiveTimeout,
		logger:              logger,
	}

	return api
}

// Run runs GRPC API server on port
func (a *API) Run(port string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.runBucketStorageCleaner(ctx)

	s := grpc.NewServer()

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("API.Run, net.Listener error %w", err)
	}

	reflection.Register(s)
	RegisterApiServer(s, a)

	err = s.Serve(l)
	if err != nil {
		return fmt.Errorf("API.Run, grpc.Serive error %w", err)
	}

	return nil
}

// AddInBlackList adds IP in black list
func (a *API) AddInBlackList(ctx context.Context, request *IPRequest) (*None, error) {
	ip, err := entities.New(request.Ip)
	if err != nil {
		return nil, err
	}

	err = a.BlackList.Add(ctx, ip)
	if err != nil {
		return nil, err
	}

	return &None{}, nil
}

// AddInWhiteList adds IP in white list
func (a *API) AddInWhiteList(ctx context.Context, request *IPRequest) (*None, error) {
	ip, err := entities.New(request.Ip)
	if err != nil {
		return nil, err
	}

	err = a.WhiteList.Add(ctx, ip)
	if err != nil {
		return nil, err
	}

	return &None{}, nil
}

// DeleteFromBlackList deletes IP from black list
func (a *API) DeleteFromBlackList(ctx context.Context, request *IPRequest) (*None, error) {
	ip, err := entities.New(request.Ip)
	if err != nil {
		return nil, err
	}

	err = a.BlackList.Delete(ctx, ip)
	if err != nil {
		return nil, err
	}

	return &None{}, nil
}

// DeleteFromWhiteList deletes IP from white list
func (a *API) DeleteFromWhiteList(ctx context.Context, request *IPRequest) (*None, error) {
	ip, err := entities.New(request.Ip)
	if err != nil {
		return nil, err
	}

	err = a.WhiteList.Delete(ctx, ip)
	if err != nil {
		return nil, err
	}

	return &None{}, nil
}

// ClearBucket clear bucket for login/password/IP
func (a *API) ClearBucket(ctx context.Context, request *BucketRequest) (*None, error) {
	var err error

	if request.Login != "" {
		err = deleteFromBucketStorage(ctx, a.LoginStorage, request.Login, "login")
		if err != nil {
			return nil, err
		}
	}

	if request.Password != "" {
		err = deleteFromBucketStorage(ctx, a.PasswordStorage, request.Password, "password")
		if err != nil {
			return nil, err
		}
	}

	if request.Ip != "" {
		ip, err := entities.New(request.Ip)
		if err != nil {
			return nil, err
		}

		err = deleteFromBucketStorage(ctx, a.IPStorage, ip, "ip")
		if err != nil {
			return nil, err
		}
	}

	return &None{}, nil
}

// Auth checks it is allowed to auth by this params (login, password, IP)
func (a *API) Auth(ctx context.Context, request *AuthRequest) (*OkResponse, error) {
	var err error

	if request.Ip == "" {
		return nil, errors.New("ip is required for Auth method")
	}

	ip, err := entities.NewWithoutMaskPart(request.Ip)
	if err != nil {
		return nil, err
	}

	var conform bool

	// if ip conform black list - no auth (even if ip conform white list)
	conform, err = a.isConformByBlackList(ctx, ip)
	if err != nil {
		return nil, err
	}

	if conform {
		return &OkResponse{Ok: false}, nil
	}

	// if ip conform white list - auth is ok
	conform, err = a.isConformByWhiteList(ctx, ip)
	if err != nil {
		return nil, err
	}

	if conform {
		return &OkResponse{Ok: true}, nil
	}

	conform, err = a.isConformByIPBucket(ctx, ip)
	if err != nil {
		return nil, err
	}

	if !conform {
		return &OkResponse{Ok: false}, nil
	}

	conform, err = a.isConformByPasswordBucket(ctx, request.Password)
	if err != nil {
		return nil, err
	}

	if !conform {
		return &OkResponse{Ok: false}, nil
	}

	conform, err = a.isConformByLoginBucket(ctx, request.Login)
	if err != nil {
		return nil, err
	}

	return &OkResponse{Ok: conform}, nil
}

// ClearBlackList clear black list of IPs
func (a *API) ClearBlackList(ctx context.Context, _ *None) (*None, error) {
	err := a.BlackList.Clear(ctx)
	if err != nil {
		return nil, err
	}

	return &None{}, nil
}

// ClearWhiteList clear white list of IPs
func (a *API) ClearWhiteList(ctx context.Context, _ *None) (*None, error) {
	err := a.WhiteList.Clear(ctx)
	if err != nil {
		return nil, err
	}

	return &None{}, nil
}

// CountBuckets get all counts of all type of buckets
func (a *API) CountBuckets(ctx context.Context, _ *None) (*BucketCountsResponse, error) {
	var err error

	var loginCount, passwordCount, ipCount int

	loginCount, err = a.LoginStorage.Count(ctx)
	if err != nil {
		return nil, err
	}

	passwordCount, err = a.PasswordStorage.Count(ctx)
	if err != nil {
		return nil, err
	}

	ipCount, err = a.IPStorage.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &BucketCountsResponse{
		Login:    uint32(loginCount),
		Password: uint32(passwordCount),
		Ip:       uint32(ipCount),
	}, nil
}

func (a *API) isConformByWhiteList(ctx context.Context, ip entities.IP) (bool, error) {
	return a.WhiteList.IsConform(ctx, ip)
}

func (a *API) isConformByBlackList(ctx context.Context, ip entities.IP) (bool, error) {
	return a.BlackList.IsConform(ctx, ip)
}

func (a *API) isConformByIPBucket(ctx context.Context, ip entities.IP) (bool, error) {
	var b entities.Bucket

	var err error

	b, err = a.getIPBucket(ctx, ip)
	if err != nil {
		return false, err
	}

	return b.IsConform(a.now()), nil
}

func (a *API) isConformByPasswordBucket(ctx context.Context, password string) (bool, error) {
	var b entities.Bucket

	var err error

	b, err = a.getPasswordBucket(ctx, password)
	if err != nil {
		return false, err
	}

	return b.IsConform(a.now()), nil
}

func (a *API) isConformByLoginBucket(ctx context.Context, login string) (bool, error) {
	var b entities.Bucket

	var err error

	b, err = a.getLoginBucket(ctx, login)
	if err != nil {
		return false, err
	}

	return b.IsConform(a.now()), nil
}

func (a *API) getIPBucket(ctx context.Context, ip entities.IP) (entities.Bucket, error) {
	return getBucketFromStorage(ctx, a.IPStorage, ip, a.IPLimit)
}

func (a *API) getPasswordBucket(ctx context.Context, password string) (entities.Bucket, error) {
	return getBucketFromStorage(ctx, a.PasswordStorage, password, a.PasswordLimit)
}

func (a *API) getLoginBucket(ctx context.Context, login string) (entities.Bucket, error) {
	return getBucketFromStorage(ctx, a.LoginStorage, login, a.LoginLimit)
}

func (a *API) runBucketStorageCleaner(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(a.bucketActiveTimeout)

	OUTER:
		for {
			select {
			case <-ticker.C:
				loginBucketsCount, _ := a.LoginStorage.ClearNotActive(ctx, a.now())
				passwordBucketsCount, _ := a.PasswordStorage.ClearNotActive(ctx, a.now())
				ipBucketsCount, _ := a.IPStorage.ClearNotActive(ctx, a.now())

				a.logDebugF("API.StorageCleaner: %d/%d/%d", loginBucketsCount, passwordBucketsCount, ipBucketsCount)

			case <-ctx.Done():
				ticker.Stop()
				break OUTER
			}
		}
	}()
}

func (a *API) logDebugF(template string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Debugf(template, args...)
	}
}

func (a *API) now() time.Time {
	if a.nowTimeFn == nil {
		a.nowTimeFn = time.Now
	}

	return a.nowTimeFn()
}

func getBucketFromStorage(ctx context.Context, storage entities.BucketStorage,
	key interface{}, limit uint) (entities.Bucket, error) {
	var b entities.Bucket

	var err error

	b, err = storage.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	// bucket not exist yet
	if b == nil {
		b = bucket.NewTokenBucketByLimitInMinute(time.Now(), limit, time.Minute)

		err = storage.Add(ctx, b, key)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func deleteFromBucketStorage(ctx context.Context, storage entities.BucketStorage, key interface{}, name string) error {
	err := storage.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("error while deleting from %s bucket storage %s", name, err)
	}

	return nil
}

func getUintFromStringMap(m map[string]string, key string, defaultVal uint) uint {
	val, ok := m[key]
	if !ok {
		return defaultVal
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return uint(valInt)
}

func getDurationFromStringMap(m map[string]string, key string, defaultVal time.Duration) time.Duration {
	val, ok := m[key]
	if !ok {
		return defaultVal
	}

	valDuration, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}

	return valDuration
}
