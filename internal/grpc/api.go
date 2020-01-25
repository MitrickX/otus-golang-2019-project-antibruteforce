package grpc

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/memory/ip"
	"github.com/spf13/viper"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/memory/bucket"
)

const (
	DefaultLoginBucketLimit    = 10
	DefaultPasswordBucketLimit = 100
	DefaultIPBucketLimit       = 1000
)

type API struct {
	BlackList entities.IPList
	WhiteList entities.IPList

	LoginBucketsStorage    entities.BucketStorage
	PasswordBucketsStorage entities.BucketStorage
	IPBucketsStorage       entities.BucketStorage

	LoginBucketLimit    uint
	PasswordBucketLimit uint
	IPBucketLimit       uint

	nowTimeFn func() time.Time
}

func NewAPIByViper(v *viper.Viper) *API {
	limits := v.GetStringMapString("limits")

	loginBucketLimit := getIntFromStringMap(limits, "login", DefaultLoginBucketLimit)
	passwordBucketLimit := getIntFromStringMap(limits, "password", DefaultPasswordBucketLimit)
	ipBucketLimit := getIntFromStringMap(limits, "ip", DefaultIPBucketLimit)

	api := &API{
		BlackList:              ip.NewList(),
		WhiteList:              ip.NewList(),
		LoginBucketsStorage:    bucket.NewStorage(),
		PasswordBucketsStorage: bucket.NewStorage(),
		IPBucketsStorage:       bucket.NewStorage(),
		LoginBucketLimit:       uint(loginBucketLimit),
		PasswordBucketLimit:    uint(passwordBucketLimit),
		IPBucketLimit:          uint(ipBucketLimit),
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
	ip := entities.IP(request.Ip)
	err := a.BlackList.Add(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &None{}, nil
}

func (a *API) AddInWhiteList(ctx context.Context, request *IPRequest) (*None, error) {
	ip := entities.IP(request.Ip)
	err := a.WhiteList.Add(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &None{}, nil
}

func (a *API) DeleteFromBlackList(ctx context.Context, request *IPRequest) (*None, error) {
	ip := entities.IP(request.Ip)
	err := a.BlackList.Delete(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &None{}, nil
}

func (a *API) DeleteFromWhiteList(ctx context.Context, request *IPRequest) (*None, error) {
	ip := entities.IP(request.Ip)
	err := a.WhiteList.Delete(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &None{}, nil
}

func (a *API) ClearBucket(ctx context.Context, request *BucketRequest) (*None, error) {
	var err error

	err = deleteFromBucketStorage(ctx, a.LoginBucketsStorage, request.Login, "login")
	if err != nil {
		return nil, err
	}

	err = deleteFromBucketStorage(ctx, a.PasswordBucketsStorage, request.Password, "password")
	if err != nil {
		return nil, err
	}

	err = deleteFromBucketStorage(ctx, a.IPBucketsStorage, entities.IP(request.Ip), "ip")
	if err != nil {
		return nil, err
	}

	return &None{}, nil
}

func (a *API) Auth(ctx context.Context, request *AuthRequest) (*OkResponse, error) {
	ip := entities.IP(request.Ip)

	var err error
	var conform bool

	conform, err = a.isConformByWhiteList(ctx, ip)
	if err != nil {
		return nil, err
	}
	if conform {
		return &OkResponse{Ok: true}, nil
	}

	conform, err = a.isConformByBlackList(ctx, ip)
	if err != nil {
		return nil, err
	}
	if conform {
		return &OkResponse{Ok: false}, nil
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
	if !conform {
		return &OkResponse{Ok: false}, nil
	}

	return &OkResponse{Ok: true}, nil

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
	return getBucketFromStorage(ctx, a.IPBucketsStorage, ip, a.IPBucketLimit)
}

func (a *API) getPasswordBucket(ctx context.Context, password string) (entities.Bucket, error) {
	return getBucketFromStorage(ctx, a.PasswordBucketsStorage, password, a.PasswordBucketLimit)
}

func (a *API) getLoginBucket(ctx context.Context, login string) (entities.Bucket, error) {
	return getBucketFromStorage(ctx, a.LoginBucketsStorage, login, a.LoginBucketLimit)
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

func getIntFromStringMap(m map[string]string, key string, defaultVal int) int {
	val, ok := m[key]
	if !ok {
		return defaultVal
	}
	valInt, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return valInt
}
