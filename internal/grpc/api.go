package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/memory/bucket"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/memory/ip"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	DefaultLoginBucketLimit    = 10
	DefaultPasswordBucketLimit = 100
	DefaultIPBucketLimit       = 1000
)

// Bucket limits config
type LimitsConfig struct {
	LoginLimit    uint
	PasswordLimit uint
	IPLimit       uint
}

func NewLimitsConfigByViper(v *viper.Viper) LimitsConfig {
	limits := v.GetStringMapString("limits")
	return LimitsConfig{
		LoginLimit:    getUintFromStringMap(limits, "login", DefaultLoginBucketLimit),
		PasswordLimit: getUintFromStringMap(limits, "password", DefaultPasswordBucketLimit),
		IPLimit:       getUintFromStringMap(limits, "ip", DefaultIPBucketLimit),
	}
}

// Set of buckets storages
type StorageSet struct {
	LoginStorage    entities.BucketStorage
	PasswordStorage entities.BucketStorage
	IPStorage       entities.BucketStorage
}

type ListSet struct {
	BlackList entities.IPList
	WhiteList entities.IPList
}

// GRPC API struct
type API struct {
	LimitsConfig
	StorageSet
	ListSet
	nowTimeFn func() time.Time
}

func NewAPIByViper(v *viper.Viper) *API {

	api := &API{
		LimitsConfig: NewLimitsConfigByViper(v),
		StorageSet: StorageSet{
			LoginStorage:    bucket.NewStorage(),
			PasswordStorage: bucket.NewStorage(),
			IPStorage:       bucket.NewStorage(),
		},
		ListSet: ListSet{
			BlackList: ip.NewList(),
			WhiteList: ip.NewList(),
		},
	}

	return api
}

func (a *API) Run(port string) error {
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

func (a *API) ClearBlackList(ctx context.Context, _ *None) (*None, error) {
	err := a.BlackList.Clear(ctx)
	if err != nil {
		return nil, err
	}
	return &None{}, nil
}

func (a *API) ClearWhiteList(ctx context.Context, _ *None) (*None, error) {
	err := a.WhiteList.Clear(ctx)
	if err != nil {
		return nil, err
	}
	return &None{}, nil
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

func (a *API) now() time.Time {
	if a.nowTimeFn == nil {
		a.nowTimeFn = time.Now
	}
	return a.nowTimeFn()
}

func getBucketFromStorage(ctx context.Context, storage entities.BucketStorage, key interface{}, limit uint) (entities.Bucket, error) {
	var b entities.Bucket
	var err error

	b, err = storage.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	// bucket not exist yet
	if b == nil {
		b = bucket.NewTokenBucketByLimitInMinute(limit)
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
