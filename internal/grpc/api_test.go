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
}

// Run test grpc server
func runTestAPI(listener net.Listener) (a *API, resultCh chan error) {
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

// Run test grpc client
func runTestClient(listener *bufconn.Listener) (client ApiClient, resultCh chan error) {
	resultCh = make(chan error, 1)

	bufDialer := func(_ context.Context, _ string) (net.Conn, error) {
		return listener.Dial()
	}
	conn, err := grpc.Dial("bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		resultCh <- fmt.Errorf("grpc Dial with bufconn connection return error %s", err)
		return
	}

	client = NewApiClient(conn)
	return
}

// Run Server and Client for grpc that bound with pipe in memory
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

// Test that ip correct added in black list
func TestAPI_AddInBlackList(t *testing.T) {
	api, client := runTestPipe(t)

	_, err := client.AddInBlackList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "add in blacklist ip `127.0.0.1`")

	cnt, err := api.BlackList.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count in blacklist after add ip `127.0.0.1`")
}

// Test that ip correct deleted from black list
func TestAPI_DeleteFromBlackList(t *testing.T) {
	api, client := runTestPipe(t)

	_, err := client.AddInBlackList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "add in blacklist ip `127.0.0.1`")

	_, err = client.AddInBlackList(context.Background(), &IPRequest{Ip: "127.0.0.0/24"})
	assertNotErrorResult(t, err, "add in blacklist ip `127.0.0.1/24`")

	_, err = client.DeleteFromBlackList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "delete from blacklist ip `127.0.0.1`")

	cnt, err := api.BlackList.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count in blacklist after one 2 and delete 1 IPs")
}

// Test that ip correct added in white list
func TestAPI_AddInWhiteList(t *testing.T) {
	api, client := runTestPipe(t)

	_, err := client.AddInWhiteList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "add in whitelist ip `127.0.0.1`")

	cnt, err := api.WhiteList.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count in whitelist after add ip `127.0.0.1`")
}

// Test that ip correct deleted from white list
func TestAPI_DeleteFromWhiteList(t *testing.T) {
	api, client := runTestPipe(t)

	_, err := client.AddInWhiteList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "add in whitelist ip `127.0.0.1`")

	_, err = client.AddInWhiteList(context.Background(), &IPRequest{Ip: "127.0.0.0/24"})
	assertNotErrorResult(t, err, "add in whitelist ip `127.0.0.1/24`")

	_, err = client.DeleteFromWhiteList(context.Background(), &IPRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "delete from whitelist ip `127.0.0.1`")

	cnt, err := api.WhiteList.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count in whitelist after one 2 and delete 1 IPs")
}

// Test that bucket correct deleted for specified login
func TestAPI_ClearBucketForLogin(t *testing.T) {
	api, client := runTestPipe(t)

	err := api.LoginStorage.Add(context.Background(), bucket.NewTokenBucketByLimitInMinute(10), "test")
	assertNotErrorResult(t, err, "add new bucket for login `test`")

	_, err = client.ClearBucket(context.Background(), &BucketRequest{Login: "test"})
	assertNotErrorResult(t, err, "delete bucket for login `test`")

	cnt, err := api.LoginStorage.Count(context.Background())
	assertCountResult(t, 0, cnt, err, "count after delete bucket for login `test`")
}

// Test that bucket correct deleted for specified password
func TestAPI_ClearBucketForPassword(t *testing.T) {
	api, client := runTestPipe(t)

	err := api.PasswordStorage.Add(context.Background(), bucket.NewTokenBucketByLimitInMinute(10), "1234")
	assertNotErrorResult(t, err, "add new bucket for password `1234`")

	_, err = client.ClearBucket(context.Background(), &BucketRequest{Password: "1234"})
	assertNotErrorResult(t, err, "delete bucket for password `1234`")

	cnt, err := api.PasswordStorage.Count(context.Background())
	assertCountResult(t, 0, cnt, err, "count after delete bucket for password `1234`")
}

// Test that bucket correct deleted for specified ip
func TestAPI_ClearBucketForIP(t *testing.T) {
	api, client := runTestPipe(t)

	err := api.IPStorage.Add(
		context.Background(),
		bucket.NewTokenBucketByLimitInMinute(10),
		entities.IP("127.0.0.1"),
	)
	assertNotErrorResult(t, err, "add new bucket for IP `127.0.0.1`")

	_, err = client.ClearBucket(context.Background(), &BucketRequest{Ip: "127.0.0.1"})
	assertNotErrorResult(t, err, "delete bucket for IP `127.0.0.1`")

	cnt, err := api.IPStorage.Count(context.Background())
	assertCountResult(t, 0, cnt, err, "count after delete bucket for IP `127.0.0.1`")
}

// Test auth when ip conform white list, because ip in white list
func TestAPI_AuthIPConformWhiteList(t *testing.T) {
	api, client := runTestPipe(t)

	// limit for ip bucket
	api.IPLimit = 1

	ip := entities.IP("127.0.0.1")

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

// Test auth when ip conform white list, because there is subnet ip in white list that conform this ip
func TestAPI_AuthIPConformWhiteList2(t *testing.T) {
	api, client := runTestPipe(t)

	api.IPLimit = 1

	ip := entities.IP("127.0.0.1")
	subnetIP := entities.IP("127.0.0.0/24")

	// now we have for ip bucket that overflowing after 1 try
	// but we add this ip in white list, so actually doesn't matter what status of bucket is
	_, _ = client.AddInWhiteList(context.Background(), &IPRequest{Ip: string(subnetIP)})

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

// Test auth when ip conform black list, because ip in black list
func TestApi_AuthIPConformBlackList1(t *testing.T) {
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

// Test auth when ip conform black list, because there is subnet ip in black list that conform this ip
func TestApi_AuthIPConformBlackList2(t *testing.T) {
	_, client := runTestPipe(t)

	ip := entities.IP("127.0.0.1")
	subnetIP := entities.IP("127.0.0.0/24")

	_, _ = client.AddInBlackList(context.Background(), &IPRequest{Ip: string(subnetIP)})

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

// Test auth when login bucket is overflowing
func TestAPI_AuthOverflowLoginBucket(t *testing.T) {
	api, client := runTestPipe(t)

	api.LoginLimit = 1
	api.PasswordLimit = 100
	api.IPLimit = 1000

	ip := entities.IP("127.0.0.1")
	login := "test"

	// first try for this login must return true
	response, err := client.Auth(context.Background(), &AuthRequest{
		Login:    login,
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("1st auth for login `%s`", login))

	// second try for this login must return false, cause of overflowing
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    login,
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, false, response, err, fmt.Sprintf("2d auth for same login `%s`", login))

	login = "test2"

	// first try for different login must return true - new bucket
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    login,
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("1st auth for different login `%s`", login))
}

// Test auth when login bucket is not overflowing
func TestAPI_AuthNotOverflowLoginBucket(t *testing.T) {
	api, client := runTestPipe(t)

	// limit is 1 try in minute for each login
	api.LoginLimit = 1
	api.PasswordLimit = 100
	api.IPLimit = 1000

	// deterministic timing
	nowTime := time.Now()
	api.nowTimeFn = func() time.Time {
		return nowTime
	}

	ip := entities.IP("127.0.0.1")
	login := "test"

	// first try for this login must return true
	response, err := client.Auth(context.Background(), &AuthRequest{
		Login:    login,
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("1st auth for login `%s`", login))

	// "wait" 1 minute
	api.nowTimeFn = func() time.Time {
		return nowTime.Add(time.Minute)
	}

	// second try for this login must return true, cause we wait and not exceed limit 1 time in minute
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    login,
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("2d auth for same login but after 1 minute wait`%s`", login))
}

// Test auth when password bucket is overflowing
func TestAPI_AuthOverflowPasswordBucket(t *testing.T) {
	api, client := runTestPipe(t)

	api.LoginLimit = 10
	api.PasswordLimit = 1 // limit 1 try in minute for each password
	api.IPLimit = 1000

	ip := entities.IP("127.0.0.1")
	password := "1234"

	// first try for this password must return true
	response, err := client.Auth(context.Background(), &AuthRequest{
		Login:    "test",
		Password: password,
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("1st auth for password `%s`", password))

	// second try for this password must return false, cause of overflowing
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    "test2", // even if different login
		Password: password,
		Ip:       string(ip),
	})
	assertOkResponse(t, false, response, err, fmt.Sprintf("2d auth for same password `%s`", password))

	// try different password
	password = "4567"

	// first try for different password must return true - new bucket
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    "test3",
		Password: password,
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("1st auth for different password `%s`", password))
}

// Test auth when password bucket is overflowing
func TestAPI_AuthNotOverflowPasswordBucket(t *testing.T) {
	api, client := runTestPipe(t)

	api.LoginLimit = 10
	api.PasswordLimit = 1 // limit 1 try in minute for each password
	api.IPLimit = 1000

	// deterministic timing
	nowTime := time.Now()
	api.nowTimeFn = func() time.Time {
		return nowTime
	}

	ip := entities.IP("127.0.0.1")
	password := "1234"

	// first try for this password must return true
	response, err := client.Auth(context.Background(), &AuthRequest{
		Login:    "test",
		Password: password,
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("1st auth for password `%s`", password))

	// "wait" 1 minute
	api.nowTimeFn = func() time.Time {
		return nowTime.Add(time.Minute)
	}

	// second try for same password must return true, cause we wait and not exceed limit 1 time in minute
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    "test", // even if for the same login
		Password: password,
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err,
		fmt.Sprintf("2d auth for same password but after 1 minute wait `%s`", password))
}

// Test auth when ip bucket is overflowing
func TestAPI_AuthOverflowIPBucket(t *testing.T) {
	api, client := runTestPipe(t)

	api.LoginLimit = 10
	api.PasswordLimit = 10
	api.IPLimit = 1 // limit 1 try in minute for each IP

	ip := entities.IP("127.0.0.1")

	// first try for this ip must return true
	response, err := client.Auth(context.Background(), &AuthRequest{
		Login:    "test",
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("1st auth for ip `%s`", ip))

	// second try for this ip must return false, cause of overflowing
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    "test2", // even if different login and password
		Password: "5678",
		Ip:       string(ip),
	})
	assertOkResponse(t, false, response, err, fmt.Sprintf("2d auth for same ip `%s`", ip))

	// try different ip
	ip = entities.IP("127.0.0.2")

	// first try for different ip must return true - new bucket
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    "test",
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("1st auth for different ip `%s`", ip))
}

// Test auth when ip bucket is overflowing
func TestAPI_AuthNotOverflowIPBucket(t *testing.T) {
	api, client := runTestPipe(t)

	api.LoginLimit = 10
	api.PasswordLimit = 10
	api.IPLimit = 1 // limit 1 try in minute for each IP

	// deterministic timing
	nowTime := time.Now()
	api.nowTimeFn = func() time.Time {
		return nowTime
	}

	ip := entities.IP("127.0.0.1")

	// first try for this ip must return true
	response, err := client.Auth(context.Background(), &AuthRequest{
		Login:    "test",
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("1st auth for ip `%s`", ip))

	// "wait" 1 minute
	api.nowTimeFn = func() time.Time {
		return nowTime.Add(time.Minute)
	}

	// second try for same ip must return true, cause we wait and not exceed limit 1 time in minute
	response, err = client.Auth(context.Background(), &AuthRequest{
		Login:    "test", // even if for the same login
		Password: "1234",
		Ip:       string(ip),
	})
	assertOkResponse(t, true, response, err, fmt.Sprintf("2d auth for same ip but after 1 minute wait `%s`", ip))
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
