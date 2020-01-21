package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/memory/bucket"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/memory/ip"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var bufConnSize = 4096

func newTestAPI() *API {
	return &API{
		blackList:              ip.NewList(),
		whiteList:              ip.NewList(),
		loginBucketsStorage:    bucket.NewStorage(),
		passwordBucketsStorage: bucket.NewStorage(),
		ipBucketsStorage:       bucket.NewStorage(),
		nowTimeFn: func() time.Time {
			return time.Now()
		},
	}
}

func runTestAPI(listener *bufconn.Listener) (a *API, resultCh chan error) {

	resultCh = make(chan error, 1)

	a = newTestAPI()

	s := grpc.NewServer()
	RegisterApiServer(s, a)

	go func() {
		err := s.Serve(listener)
		if err != nil {
			resultCh <- fmt.Errorf("test server exited with error %s", err)
		}
	}()

	return
}

func runTestClient(listener *bufconn.Listener) (client ApiClient, resultCh chan error) {
	resultCh = make(chan error, 1)

	bufDialer := func(_ context.Context, _ string) (net.Conn, error) {
		return listener.Dial()
	}
	conn, err := grpc.Dial("bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		resultCh <- fmt.Errorf("grpc Dial with bufconn connection return error %s\n", err)
		return
	}

	client = NewApiClient(conn)
	return
}

func runTestPipe(t *testing.T) (*API, ApiClient) {

	listener := bufconn.Listen(bufConnSize)

	var client ApiClient
	var clientResCh chan error

	service, serverResCh := runTestAPI(listener)

	// If error return error right away not run client and close listener
	select {
	case err := <-serverResCh:
		if err != nil {
			t.Error(err)
		}
		_ = listener.Close()
		return nil, nil
	default:
		client, clientResCh = runTestClient(listener)
	}

	// If server or client return error close listener
	go func() {
		select {
		case err := <-serverResCh:
			if err != nil {
				t.Error(err)
			}
			_ = listener.Close()
		case err := <-clientResCh:
			if err != nil {
				t.Error(err)
			}
			_ = listener.Close()
		}
	}()

	return service, client
}

func TestAPI_AddInBlackList(t *testing.T) {
	api, client := runTestPipe(t)

	_, err := client.AddInBlackList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "add in blacklist ip `127.0.0.1`")

	cnt, err := api.blackList.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count in blacklist after add ip `127.0.0.1`")
}

func TestAPI_DeleteFromBlackList(t *testing.T) {
	api, client := runTestPipe(t)

	_, err := client.AddInBlackList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "add in blacklist ip `127.0.0.1`")

	_, err = client.AddInBlackList(context.Background(), &IPRequest{Ip: "127.0.0.0/24"})
	assertNotErrorResult(t, err, "add in blacklist ip `127.0.0.1/24`")

	_, err = client.DeleteFromBlackList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "delete from blacklist ip `127.0.0.1`")

	cnt, err := api.blackList.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count in blacklist after one 2 and delete 1 IPs")
}

func TestAPI_AddInWhiteList(t *testing.T) {
	api, client := runTestPipe(t)

	_, err := client.AddInWhiteList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "add in whitelist ip `127.0.0.1`")

	cnt, err := api.whiteList.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count in whitelist after add ip `127.0.0.1`")
}

func TestAPI_DeleteFromWhiteList(t *testing.T) {
	api, client := runTestPipe(t)

	_, err := client.AddInWhiteList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "add in whitelist ip `127.0.0.1`")

	_, err = client.AddInWhiteList(context.Background(), &IPRequest{Ip: "127.0.0.0/24"})
	assertNotErrorResult(t, err, "add in whitelist ip `127.0.0.1/24`")

	_, err = client.DeleteFromWhiteList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "delete from whitelist ip `127.0.0.1`")

	cnt, err := api.whiteList.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count in whitelist after one 2 and delete 1 IPs")
}

func TestAPI_ClearBucketForLogin(t *testing.T) {
	api, client := runTestPipe(t)

	err := api.loginBucketsStorage.Add(context.Background(), bucket.NewTokenBucketByLimitInMinute(10), "test")
	assertNotErrorResult(t, err, "add new bucket for login `test`")

	_, err = client.ClearBucket(context.Background(), &BucketRequest{Login: "test"})
	assertNotErrorResult(t, err, "delete bucket for login `test`")

	cnt, err := api.loginBucketsStorage.Count(context.Background())
	assertCountResult(t, 0, cnt, err, "count after delete bucket for login `test`")
}

func TestAPI_ClearBucketForPassword(t *testing.T) {
	api, client := runTestPipe(t)

	err := api.passwordBucketsStorage.Add(context.Background(), bucket.NewTokenBucketByLimitInMinute(10), "1234")
	assertNotErrorResult(t, err, "add new bucket for password `1234`")

	_, err = client.ClearBucket(context.Background(), &BucketRequest{Password: "1234"})
	assertNotErrorResult(t, err, "delete bucket for password `1234`")

	cnt, err := api.passwordBucketsStorage.Count(context.Background())
	assertCountResult(t, 0, cnt, err, "count after delete bucket for password `1234`")
}

func TestAPI_ClearBucketForIP(t *testing.T) {
	api, client := runTestPipe(t)

	err := api.ipBucketsStorage.Add(
		context.Background(),
		bucket.NewTokenBucketByLimitInMinute(10),
		entities.IP("127.0.0.1"),
	)
	assertNotErrorResult(t, err, "add new bucket for IP `127.0.0.1`")

	_, err = client.ClearBucket(context.Background(), &BucketRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "delete bucket for IP `127.0.0.1`")

	cnt, err := api.ipBucketsStorage.Count(context.Background())
	assertCountResult(t, 0, cnt, err, "count after delete bucket for IP `127.0.0.1`")

}

func TestAPI_AuthIPInWhiteList(t *testing.T) {
	api, client := runTestPipe(t)

	ip := entities.IP("127.0.0.1")
	b := bucket.NewTokenBucketByLimitInMinute(1)

	err := api.ipBucketsStorage.Add(context.Background(), b, ip)
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.1`")

	// now we have for ip bucket that overflowing after 1 try
	// but we add this ip in white list, so actually doesn't matter what status of bucket is
	_, _ = client.AddInWhiteList(context.Background(), &IPRequest{Ip: string(ip)})

	// first auth try, must be ok
	response, err := client.Auth(context.Background(), &AuthRequest{
		Login:    "test",
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, "1st auth of ip `127.0.0.1`")

	// second auth try, must be ok cause of white list
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    "test",
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, "2nd auth of ip `127.0.0.1`")
}

func TestApi_AuthIPInBlackList(t *testing.T) {
	_, client := runTestPipe(t)

	ip := entities.IP("127.0.0.1")

	_, _ = client.AddInBlackList(context.Background(), &IPRequest{Ip: string(ip)})

	// auth try must return false, cause ip in black list
	response, err := client.Auth(context.Background(), &AuthRequest{
		Login:    "test",
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, false, response, err, "1st auth of ip `127.0.0.1`")

	// auth try must return false, cause ip in black list
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    "test",
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, false, response, err, "2nd auth of ip `127.0.0.1`")

}

func assertNotErrorResult(t *testing.T, err error, prefix string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %s", prefix, err)
	}
}

func assertCountResult(t *testing.T, expected int, count int, err error, prefix string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %s", prefix, err)
	}
	if count != expected {
		t.Fatalf("%s: unexpected count %d instreadof %d", prefix, count, expected)
	}
}

func assertOkResponse(t *testing.T, expected bool, response *OkResponse, err error, prefix string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %s", prefix, err)
	}
	if response.Ok != expected {
		t.Fatalf("%s: unexpected %t instreadof %t", prefix, response.Ok, expected)
	}
}
