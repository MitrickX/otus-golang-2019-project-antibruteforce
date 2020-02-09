// +build !unit

package ip

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/logger"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/sql"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/sql/ip"
	"github.com/spf13/viper"
)

var db *sqlx.DB
var cfg sql.Config

const DefaultContextTimeout = 3 * time.Second

func init() {
	v := viper.GetViper()

	logger.InitLogger(v)
	l := logger.GetLogger()

	var err error

	cfg = sql.NewConfigByEnv()

	db, err = sql.Connect(cfg)
	if err != nil {
		l.Fatalf("Connect to DB failed: %s", err)
	}

	WaitForTable()
}

func WaitForTable() {
	var err error

	// Wait for external service create ip_list table

	// we will run with timeout in loop, so if any cancel not stop yet will called in defer
	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	for i := 0; i < cfg.ConnectRetries; i++ {
		var ok bool
		var ctx context.Context

		ctx, cancel = context.WithTimeout(context.Background(), DefaultContextTimeout)

		ok, err = sql.IsTableExists(ctx, db, cfg.DbName, "ip_list")

		// cancel after method called
		cancel()

		if ok {
			break
		}

		time.Sleep(time.Second)
	}

	if err != nil {
		log.Fatal(err)
	}
}

// New clean list
func NewList() *ip.List {
	list := ip.NewList(db, "black")

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	err := list.Clear(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return list
}

func TestList_Add(t *testing.T) {
	list := NewList()

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	err = list.Add(ctx, entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1`")

	err = list.Add(ctx, entities.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "add ip `127.0.0.2`")

	err = list.Add(ctx, entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1` (2nd time)")

	cnt, err := list.Count(ctx)
	assertCountResult(t, 2, cnt, err, "count after adds")
}

func TestList_Delete(t *testing.T) {
	list := NewList()

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	err = list.Add(ctx, entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1`")

	err = list.Add(ctx, entities.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "add ip `127.0.0.2`")

	// Delete not existing
	err = list.Delete(ctx, entities.IP("127.0.0.3"))
	assertNotErrorResult(t, err, "delete `127.0.0.3`")

	err = list.Delete(ctx, entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "delete `127.0.0.1`")

	cnt, err := list.Count(ctx)
	assertCountResult(t, 1, cnt, err, "count")
}

func TestList_Has(t *testing.T) {
	list := NewList()

	var ok bool
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	err = list.Add(ctx, entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1`")

	ok, err = list.Has(ctx, entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "has `127.0.0.1`")

	err = list.Delete(ctx, entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "delete ip `127.0.0.1`")

	ok, err = list.Has(ctx, entities.IP("127.0.0.1"))
	if err != nil {
		t.Fatalf("has `127.0.0.1`: unexpected error %s", err)
	}
	if ok {
		t.Fatalf("expected there is not ip `127.0.0.1` in list")
	}
}

func TestList_IsConform_IPv4(t *testing.T) {
	list := NewList()

	var ok bool
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	err = list.Add(ctx, entities.IP("127.0.0.0/24"))
	assertNotErrorResult(t, err, "add ip `127.0.0.0/24`")

	ok, err = list.IsConform(ctx, entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "is conform 127.0.0.1")

	ok, err = list.IsConform(ctx, entities.IP("127.0.0.30"))
	assertOkResult(t, ok, err, "is conform 127.0.0.30")

	ok, err = list.IsConform(ctx, entities.IP("127.0.0.255"))
	assertOkResult(t, ok, err, "is conform 127.0.0.255")

	ok, err = list.IsConform(ctx, entities.IP("128.0.0.4"))
	assertNotOkResult(t, ok, err, "is not conform 128.0.0.4")
}

func TestList_IsConform_IPv6(t *testing.T) {
	list := NewList()

	var ok bool
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	// Full IP of subnet: 2001:DB8:0000:1234:0000:0000:0000:0000
	// Where prefix /64 is 2001:DB8:0000:1234
	subnetIP := entities.IP("2001:DB8:0:1234::/64")

	// Full host IP
	// 2001:DB8:0:1234::2 =>
	//   2001:0DB8:0000:1234 (/64 subnet IP)
	//   0000:0000:0000:0002 (host IP)
	hostIP := entities.IP("2001:DB8:0:1234::2")

	err = list.Add(ctx, subnetIP)
	assertNotErrorResult(t, err, fmt.Sprintf("add ip `%s`", subnetIP))

	ok, err = list.IsConform(ctx, hostIP)
	assertOkResult(t, ok, err, fmt.Sprintf("is conform `%s`", hostIP))

	// Full host IP
	// 2001:DB8:0:1234::1:20 =>
	//   2001:0DB8:0000:1234 (/64 subnet IP)
	//   0000:0000:0001:0020 (host IP)
	hostIP = entities.IP("2001:DB8:0:1234::1:20")

	ok, err = list.IsConform(ctx, hostIP)
	assertOkResult(t, ok, err, fmt.Sprintf("is conform `%s`", hostIP))

	// Full host IP
	// 2001:DB8:0:1234::ffff =>
	//   2001:0DB8:0000:1234 (/64 subnet IP)
	//   0000:0000:0000:ffff (host IP)
	hostIP = entities.IP("2001:DB8:0:1234::ffff")

	ok, err = list.IsConform(ctx, hostIP)
	assertOkResult(t, ok, err, fmt.Sprintf("is conform `%s`", hostIP))

	// Full host IP
	// 2001:DB8:0:1235::2 =>
	//   2001:0DB8:0000:1235 (/64 subnet IP), not prefix of current subnetIP
	//   0000:0000:0000:0002 (host IP)
	hostIP = entities.IP("2001:DB8:0:1235::2")

	ok, err = list.IsConform(ctx, hostIP)
	assertNotOkResult(t, ok, err, fmt.Sprintf("is not conform `%s`", hostIP))
}

func TestList_Count(t *testing.T) {
	list := NewList()

	var cnt int
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	cnt, err = list.Count(ctx)
	assertCountResult(t, 0, cnt, err, "count")

	err = list.Add(ctx, entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1`")

	cnt, err = list.Count(ctx)
	assertCountResult(t, 1, cnt, err, "count")

	err = list.Add(ctx, entities.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "add ip `127.0.0.2`")

	cnt, err = list.Count(ctx)
	assertCountResult(t, 2, cnt, err, "count")

	err = list.Delete(ctx, entities.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "delete `127.0.0.1`")

	cnt, err = list.Count(ctx)
	assertCountResult(t, 1, cnt, err, "count")

	err = list.Delete(ctx, entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "delete `127.0.0.1`")

	cnt, err = list.Count(ctx)
	assertCountResult(t, 0, cnt, err, "count")
}

func TestList_Clear(t *testing.T) {
	list := NewList()

	var cnt int
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	cnt, err = list.Count(ctx)
	assertCountResult(t, 0, cnt, err, "count")

	err = list.Add(ctx, entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1`")

	err = list.Clear(ctx)
	assertNotErrorResult(t, err, "clear list")

	cnt, err = list.Count(ctx)
	assertCountResult(t, 0, cnt, err, "count")
}

func assertNotErrorResult(t *testing.T, err error, prefix string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %s", prefix, err)
	}
}

func assertOkResult(t *testing.T, ok bool, err error, prefix string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %s", prefix, err)
	}
	if !ok {
		t.Fatalf("%s: expected be true", prefix)
	}
}

func assertNotOkResult(t *testing.T, ok bool, err error, prefix string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %s", prefix, err)
	}
	if ok {
		t.Fatalf("%s: expected be false", prefix)
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
