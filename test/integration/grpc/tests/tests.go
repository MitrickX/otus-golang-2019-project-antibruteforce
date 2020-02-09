// +build !unit

package tests

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/DATA-DOG/godog/gherkin"
	grpcAPI "github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"
	"google.golang.org/grpc"
)

const (
	AddInBlackListMethodName      = "AddInBlackList"
	AddInWhiteListMethodName      = "AddInWhiteList"
	DeleteFromBlackListMethodName = "DeleteFromBlackList"
	DeleteFromWhiteListMethodName = "DeleteFromWhiteList"
	ClearBucketMethodName         = "ClearBucket"
	AuthMethodName                = "Auth"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ipListMethod func(context.Context, *grpcAPI.IPRequest, ...grpc.CallOption) (*grpcAPI.None, error)

func docStringToString(data *gherkin.DocString) string {
	replacer := strings.NewReplacer("\n", "", "\t", "")
	return replacer.Replace(data.Content)
}

func docStringToAuthRequest(params *gherkin.DocString) (*grpcAPI.AuthRequest, error) {
	query := docStringToString(params)

	p, err := url.ParseQuery(query)
	if err != nil {
		return nil, fmt.Errorf("parse params failed: %s", err)
	}

	request := &grpcAPI.AuthRequest{
		Login:    p.Get("login"),
		Password: p.Get("password"),
		Ip:       p.Get("ip"),
	}

	if request.Login == "random" {
		/* #nosec */
		request.Login = "l" + strconv.Itoa(rand.Int())
	}

	if request.Password == "random" {
		/* #nosec */
		request.Password = "p" + strconv.Itoa(rand.Int())
	}

	return request, nil
}

func docStringToIPRequest(param *gherkin.DocString) (*grpcAPI.IPRequest, error) {
	query := docStringToString(param)

	p, err := url.ParseQuery(query)
	if err != nil {
		return nil, fmt.Errorf("parse params failed: %s", err)
	}

	return &grpcAPI.IPRequest{Ip: p.Get("ip")}, nil
}

func docStringToBucketRequest(params *gherkin.DocString) (*grpcAPI.BucketRequest, error) {
	query := docStringToString(params)

	p, err := url.ParseQuery(query)
	if err != nil {
		return nil, fmt.Errorf("parse params failed: %s", err)
	}

	return &grpcAPI.BucketRequest{
		Login:    p.Get("login"),
		Password: p.Get("password"),
		Ip:       p.Get("ip"),
	}, nil
}

func stringTimesToInt(val string) (int, error) {
	if val == "" {
		return 1, nil
	}

	if val == "loginLimit" {
		return int(GetConfig().LoginLimit), nil
	}

	if val == "passwordLimit" {
		return int(GetConfig().PasswordLimit), nil
	}

	if val == "ipLimit" {
		return int(GetConfig().IPLimit), nil
	}

	times, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("times param is not int: %s", err)
	}

	return times, nil
}

func isMethod(method string) bool {
	switch method {
	case
		AddInBlackListMethodName,
		AddInWhiteListMethodName,
		DeleteFromBlackListMethodName,
		DeleteFromWhiteListMethodName,
		ClearBucketMethodName,
		AuthMethodName:
		return true
	default:
		return false
	}
}

func isIPListMethod(method string) bool {
	switch method {
	case
		AddInBlackListMethodName,
		AddInWhiteListMethodName,
		DeleteFromBlackListMethodName,
		DeleteFromWhiteListMethodName:
		return true
	default:
		return false
	}
}

func getIPListMethodByName(name string) ipListMethod {
	client := GetConfig().apiClient

	switch name {
	case AddInBlackListMethodName:
		return client.AddInBlackList
	case AddInWhiteListMethodName:
		return client.AddInWhiteList
	case DeleteFromBlackListMethodName:
		return client.DeleteFromBlackList
	case DeleteFromWhiteListMethodName:
		return client.DeleteFromWhiteList
	default:
		return nil
	}
}
