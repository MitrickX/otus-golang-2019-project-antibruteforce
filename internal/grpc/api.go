package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/memory/bucket"
)

type API struct {
	blackList              entities.IPList
	whiteList              entities.IPList
	loginBucketsStorage    entities.BucketStorage
	passwordBucketsStorage entities.BucketStorage
	ipBucketsStorage       entities.BucketStorage
	loginBucketLimit       uint
	passwordBucketLimit    uint
	ipBucketLimit          uint
	nowTimeFn              func() time.Time
}

func NewAPI() *API {
	return &API{
		nowTimeFn: func() time.Time {
			return time.Now()
		},
		loginBucketLimit:    10,
		passwordBucketLimit: 100,
		ipBucketLimit:       1000,
	}
}

func (a *API) AddInBlackList(ctx context.Context, request *IPRequest) (*None, error) {
	ip := entities.IP(request.Ip)
	err := a.blackList.Add(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &None{}, nil
}

func (a *API) AddInWhiteList(ctx context.Context, request *IPRequest) (*None, error) {
	ip := entities.IP(request.Ip)
	err := a.whiteList.Add(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &None{}, nil
}

func (a *API) DeleteFromBlackList(ctx context.Context, request *IPRequest) (*None, error) {
	ip := entities.IP(request.Ip)
	err := a.blackList.Delete(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &None{}, nil
}

func (a *API) DeleteFromWhiteList(ctx context.Context, request *IPRequest) (*None, error) {
	ip := entities.IP(request.Ip)
	err := a.whiteList.Delete(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &None{}, nil
}

func (a *API) ClearBucket(ctx context.Context, request *BucketRequest) (*None, error) {
	var err error

	err = deleteFromBucketStorage(ctx, a.loginBucketsStorage, request.Login, "login")
	if err != nil {
		return nil, err
	}

	err = deleteFromBucketStorage(ctx, a.passwordBucketsStorage, request.Password, "password")
	if err != nil {
		return nil, err
	}

	err = deleteFromBucketStorage(ctx, a.ipBucketsStorage, entities.IP(request.Ip), "ip")
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
	return a.whiteList.IsConform(ctx, ip)
}

func (a *API) isConformByBlackList(ctx context.Context, ip entities.IP) (bool, error) {
	return a.blackList.IsConform(ctx, ip)
}

func (a *API) isConformByIPBucket(ctx context.Context, ip entities.IP) (bool, error) {
	var b entities.Bucket
	var err error

	b, err = a.getIPBucket(ctx, ip)
	if err != nil {
		return false, err
	}

	return b.IsConform(a.nowTimeFn()), nil
}

func (a *API) isConformByPasswordBucket(ctx context.Context, password string) (bool, error) {
	var b entities.Bucket
	var err error

	b, err = a.getPasswordBucket(ctx, password)
	if err != nil {
		return false, err
	}

	return b.IsConform(a.nowTimeFn()), nil
}

func (a *API) isConformByLoginBucket(ctx context.Context, login string) (bool, error) {
	var b entities.Bucket
	var err error

	b, err = a.getLoginBucket(ctx, login)
	if err != nil {
		return false, err
	}

	return b.IsConform(a.nowTimeFn()), nil
}

func (a *API) getIPBucket(ctx context.Context, ip entities.IP) (entities.Bucket, error) {
	return getBucketFromStorage(ctx, a.ipBucketsStorage, ip, a.ipBucketLimit)
}

func (a *API) getPasswordBucket(ctx context.Context, password string) (entities.Bucket, error) {
	return getBucketFromStorage(ctx, a.passwordBucketsStorage, password, a.passwordBucketLimit)
}

func (a *API) getLoginBucket(ctx context.Context, login string) (entities.Bucket, error) {
	return getBucketFromStorage(ctx, a.loginBucketsStorage, login, a.loginBucketLimit)
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
