package tests

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"

	"github.com/DATA-DOG/godog"
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
		_, err := method(ctx, docStringToIPRequest(params))
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

func (t *featureTest) listWithIp(kind string, ip string) error {

	if kind != "black" && kind != "white" {
		return fmt.Errorf("unexpected kind of list `%s`", kind)
	}

	ip = strings.TrimSpace(ip)

	cfg := GetConfig()
	apiClient := cfg.apiClient

	request := &grpc.IPRequest{Ip: ip}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	var err error
	if kind == "black" {
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

	query := docStringToString(params)

	p, err := url.ParseQuery(query)
	if err != nil {
		return fmt.Errorf("parse params failed: %s", err)
	}

	cfg := GetConfig()
	apiClient := cfg.apiClient

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	// we call auth so bucket will be created
	request := &grpc.AuthRequest{
		Login:    p.Get("login"),
		Password: p.Get("password"),
		Ip:       p.Get("ip"),
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
			} else {
				return fmt.Errorf("unexpected response error (# %d) `%s` instreadof `%s`", index, err, expectedErr)
			}
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
			} else {
				return fmt.Errorf("unexpected response error (# %d) `%t` instreadof `%t`", index, res.Ok, expectedBool)
			}
		}
	}

	return nil
}

func FeatureContext(s *godog.Suite, t *featureTest) {
	s.Step(`^I call method "([^"]*)" with params:$`, t.iCallMethodWithParams)
	s.Step(`^I call "([^"]*)" times method "([^"]*)" with params:$`, t.iCallTimesMethodWithParams)
	s.Step(`^The error must be "([^"]*)"$`, t.theErrorMustBe)
	s.Step(`^"([^"]*)" list with ip="([^"]*)"$`, t.listWithIp)
	s.Step(`^bucket for$`, t.bucketFor)
	s.Step(`^The result must be "([^"]*)"$`, t.theResultMustBe)
}
