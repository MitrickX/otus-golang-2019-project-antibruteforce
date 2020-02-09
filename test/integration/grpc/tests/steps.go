// +build !unit

package tests

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"

	"github.com/DATA-DOG/godog"
)

const (
	BlackKind = "black"
	WhiteKind = "white"
)

type featureTest struct {
	responseErrors []error
	okResponses    []*grpc.OkResponse
}

func newFeatureTest() *featureTest {
	return &featureTest{}
}

func (t *featureTest) iCallMethodWithParams(methodName string, params *gherkin.DocString) error {
	return t.iCallIntTimesMethodWithParams(1, methodName, params)
}

func (t *featureTest) iCallTimesMethodWithParams(times, methodName string, params *gherkin.DocString) error {
	n, err := stringTimesToInt(times)
	if err != nil {
		return fmt.Errorf("couldn't convert string `times` to int %s", err)
	}

	return t.iCallIntTimesMethodWithParams(n, methodName, params)
}

func (t *featureTest) iCallIntTimesMethodWithParams(times int, methodName string, params *gherkin.DocString) error {
	if !isMethod(methodName) {
		return fmt.Errorf("unexpected method %s", methodName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	if isIPListMethod(methodName) {
		method := getIPListMethodByName(methodName)
		if method == nil {
			return fmt.Errorf("coudn't find grpc method by name %s", methodName)
		}

		request, err := docStringToIPRequest(params)
		if err != nil {
			return fmt.Errorf("couldn't convert input params to ip request %s", err)
		}

		_, err = method(ctx, request)
		t.responseErrors = []error{err}

		return nil
	}

	cfg := GetConfig()
	apiClient := cfg.apiClient

	if methodName == "ClearBucket" {
		request, err := docStringToBucketRequest(params)
		if err != nil {
			return fmt.Errorf("couldn't convert input params to bucket request %s", err)
		}

		_, err = apiClient.ClearBucket(ctx, request)
		t.responseErrors = []error{err}

		return nil
	}

	t.responseErrors = nil
	t.okResponses = nil

	for i := 0; i < times; i++ {
		request, err := docStringToAuthRequest(params)
		if err != nil {
			return fmt.Errorf("couldn't convert input params to auth request %s", err)
		}

		okResponse, err := apiClient.Auth(ctx, request)
		t.responseErrors = append(t.responseErrors, err)
		t.okResponses = append(t.okResponses, okResponse)
	}

	return nil
}

func (t *featureTest) listWithIP(kind string, ip string) error {
	if kind != BlackKind && kind != WhiteKind {
		return fmt.Errorf("unexpected kind of list `%s`", kind)
	}

	ip = strings.TrimSpace(ip)

	cfg := GetConfig()
	apiClient := cfg.apiClient

	request := &grpc.IPRequest{Ip: ip}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	var err error
	if kind == BlackKind {
		_, err = apiClient.AddInBlackList(ctx, request)
	} else {
		_, err = apiClient.AddInWhiteList(ctx, request)
	}

	if err != nil {
		return fmt.Errorf("unexpected error when add ip %s in list", ip)
	}

	return nil
}

func (t *featureTest) bucketFor(params *gherkin.DocString) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	cfg := GetConfig()
	apiClient := cfg.apiClient

	request, err := docStringToAuthRequest(params)
	if err != nil {
		return fmt.Errorf("couldn't convert input params to auth request %s", err)
	}

	// we call auth so bucket will be created

	// Ip must not be empty, cause it is required for auth
	if request.Ip == "" {
		request.Ip = "127.0.0.1"
	}

	response, err := apiClient.Auth(ctx, request)

	if err != nil {
		return fmt.Errorf("unexpected error when call auth %s", err)
	}

	if !response.Ok {
		return fmt.Errorf("bucket not created for %+v", request)
	}

	return nil
}

func (t *featureTest) theErrorMustBe(expected string) error {
	l := len(t.responseErrors)

	if l == 0 {
		return errors.New("expected some method be called")
	}

	var expectedErr error
	if expected != "nil" {
		expectedErr = errors.New(expected)
	}

	for index, err := range t.responseErrors {
		if err != expectedErr {
			if l == 1 {
				return fmt.Errorf("unexpected response error `%s` instreadof `%s`", err, expectedErr)
			}

			return fmt.Errorf("unexpected response error (# %d) `%s` instreadof `%s`", index, err, expectedErr)
		}
	}

	return nil
}

func (t *featureTest) theResultMustBe(expected string) error {
	l := len(t.responseErrors)

	if l == 0 {
		return errors.New("expected Auth method be called")
	}

	expectedBool := false
	if strings.ToLower(expected) == "true" {
		expectedBool = true
	}

	for index, res := range t.okResponses {
		if res.Ok != expectedBool {
			if l == 1 {
				return fmt.Errorf("unexpected response `%t` instreadof `%t`", res.Ok, expectedBool)
			}

			return fmt.Errorf("unexpected response error (# %d) `%t` instreadof `%t`", index, res.Ok, expectedBool)
		}
	}

	return nil
}

func (t *featureTest) cleanBucketFor(params *gherkin.DocString) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	request, err := docStringToBucketRequest(params)
	if err != nil {
		return fmt.Errorf("couldn't convert input params to bucket request %s", err)
	}

	_, err = GetConfig().apiClient.ClearBucket(ctx, request)

	if err != nil {
		return fmt.Errorf("unexpected error when call clear bucket %s", err)
	}

	return nil
}

func (t *featureTest) cleanList(kind string) error {
	if kind != BlackKind && kind != WhiteKind {
		return fmt.Errorf("unexpected kind of list `%s`", kind)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	var err error
	if kind == BlackKind {
		_, err = GetConfig().apiClient.ClearBlackList(ctx, &grpc.None{})
	} else {
		_, err = GetConfig().apiClient.ClearWhiteList(ctx, &grpc.None{})
	}

	if err != nil {
		return fmt.Errorf("unexpected erorr when clear %s list %s", kind, err)
	}

	return nil
}

func (t *featureTest) waitMinute(n int) error {
	time.Sleep(time.Duration(n) * time.Minute)
	return nil
}

func FeatureContext(s *godog.Suite, t *featureTest) {
	s.Step(`^Clean bucket for$`, t.cleanBucketFor)
	s.Step(`^Clean "([^"]*)" list$`, t.cleanList)
	s.Step(`^I call method "([^"]*)" with params:$`, t.iCallMethodWithParams)
	s.Step(`^I call "([^"]*)" times method "([^"]*)" with params:$`, t.iCallTimesMethodWithParams)
	s.Step(`^The error must be "([^"]*)"$`, t.theErrorMustBe)
	s.Step(`^"([^"]*)" list with ip="([^"]*)"$`, t.listWithIP)
	s.Step(`^bucket for$`, t.bucketFor)
	s.Step(`^The result must be "([^"]*)"$`, t.theResultMustBe)
	s.Step(`^Wait (\d+) minute$`, t.waitMinute)
}
