package tests

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"

	"github.com/DATA-DOG/godog"
)

type featureTest struct {
	responseError error
}

func newFeatureTest() *featureTest {
	return new(featureTest)
}

func (t *featureTest) iCallMethodWithParams(methodName string, params *gherkin.DocString) error {

	cfg := GetConfig()
	apiClient := cfg.apiClient

	ctx, _ := context.WithTimeout(context.Background(), cfg.timeout)

	replacer := strings.NewReplacer("\n", "", "\t", "")
	query := replacer.Replace(params.Content)

	switch methodName {
	case "AddInBlackList":
		_, t.responseError = apiClient.AddInBlackList(ctx, &grpc.IPRequest{Ip: query})
	case "AddInWhiteList":
		_, t.responseError = apiClient.AddInWhiteList(ctx, &grpc.IPRequest{Ip: query})
	case "DeleteFromBlackList":
		_, t.responseError = apiClient.DeleteFromBlackList(ctx, &grpc.IPRequest{Ip: query})
	case "DeleteFromWhiteList":
		_, t.responseError = apiClient.DeleteFromWhiteList(ctx, &grpc.IPRequest{Ip: query})
	case "ClearBucket":
		//_, t.responseError = apiClient.ClearBucket(ctx, &grpc.IPRequest{Ip: query})
	case "Auth":
	default:
		return fmt.Errorf("unexpected method %s", methodName)
	}

	return nil
}

func (t *featureTest) theErrorMustBe(expected string) error {
	var expectedErr error
	if expected != "nil" {
		expectedErr = errors.New(expected)
	}
	if t.responseError != expectedErr {
		return fmt.Errorf("unexpected response error `%s` instreadof `%s`", t.responseError, expectedErr)
	}
	return nil
}

func FeatureContext(s *godog.Suite, t *featureTest) {
	s.Step(`^I call method "([^"]*)" with params:$`, t.iCallMethodWithParams)
	s.Step(`^The error must be "([^"]*)"$`, t.theErrorMustBe)
}
