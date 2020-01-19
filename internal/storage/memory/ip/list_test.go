package ip

import (
	"context"
	"testing"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities/ip"
)

func TestList_Add(t *testing.T) {
	list := NewList()

	var err error

	err = list.Add(context.Background(), ip.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1`")

	err = list.Add(context.Background(), ip.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "add ip `127.0.0.2`")

	err = list.Add(context.Background(), ip.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1` (2nd time)")

	cnt, err := list.Count(context.Background())
	assertCountResult(t, 2, cnt, err, "count")
}

func TestList_Delete(t *testing.T) {
	list := NewList()

	var err error

	err = list.Add(context.Background(), ip.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1`")

	err = list.Add(context.Background(), ip.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "add ip `127.0.0.2`")

	// Delete not existing
	err = list.Delete(context.Background(), ip.IP("127.0.0.3"))
	assertNotErrorResult(t, err, "delete `127.0.0.3`")

	err = list.Delete(context.Background(), ip.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "delete `127.0.0.1`")

	cnt, err := list.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count")
}

func TestList_Has(t *testing.T) {
	list := NewList()

	var ok bool
	var err error

	err = list.Add(context.Background(), ip.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1`")

	ok, err = list.Has(context.Background(), ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "has `127.0.0.1`")

	err = list.Delete(context.Background(), ip.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "delete ip `127.0.0.1`")

	ok, err = list.Has(context.Background(), ip.IP("127.0.0.1"))
	if err != nil {
		t.Fatalf("has `127.0.0.1`: unexpected error %s", err)
	}
	if ok {
		t.Fatalf("expected there is not ip `127.0.0.1` in list")
	}
}

func TestList_Count(t *testing.T) {
	list := NewList()

	var cnt int
	var err error

	cnt, err = list.Count(context.Background())
	assertCountResult(t, 0, cnt, err, "count")

	err = list.Add(context.Background(), ip.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add ip `127.0.0.1`")

	cnt, err = list.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count")

	err = list.Add(context.Background(), ip.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "add ip `127.0.0.2`")

	cnt, err = list.Count(context.Background())
	assertCountResult(t, 2, cnt, err, "count")

	err = list.Delete(context.Background(), ip.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "delete `127.0.0.1`")

	cnt, err = list.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count")

	err = list.Delete(context.Background(), ip.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "delete `127.0.0.1`")

	cnt, err = list.Count(context.Background())
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
		t.Fatalf("%s: expected be successfull", prefix)
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
