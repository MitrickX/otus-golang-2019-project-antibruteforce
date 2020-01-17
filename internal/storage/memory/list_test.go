package memory

import (
	"testing"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
)

func TestIPList_Add(t *testing.T) {
	list := NewIPList()

	var ok bool
	var err error

	ok, err = list.Add(entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add ip `127.0.0.1`")

	ok, err = list.Add(entities.IP("127.0.0.2"))
	assertOkResult(t, ok, err, "add ip `127.0.0.2`")

	ok, err = list.Add(entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add ip `127.0.0.1` (2nd time)")

	cnt, err := list.Count()
	assertCountResult(t, 2, cnt, err, "count")
}

func TestIPList_Delete(t *testing.T) {
	list := NewIPList()

	var ok bool
	var err error

	ok, err = list.Add(entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add ip `127.0.0.1`")

	ok, err = list.Add(entities.IP("127.0.0.2"))
	assertOkResult(t, ok, err, "add ip `127.0.0.2`")

	// Delete not existing
	ok, err = list.Delete(entities.IP("127.0.0.3"))
	assertOkResult(t, ok, err, "delete `127.0.0.3`")

	ok, err = list.Delete(entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "delete `127.0.0.1`")

	cnt, err := list.Count()
	assertCountResult(t, 1, cnt, err, "count")
}

func TestIPList_Has(t *testing.T) {
	list := NewIPList()

	var ok bool
	var err error

	ok, err = list.Add(entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add ip `127.0.0.1`")

	ok, err = list.Has(entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "has `127.0.0.1`")

	ok, err = list.Delete(entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "delete ip `127.0.0.1`")

	ok, err = list.Has(entities.IP("127.0.0.1"))
	if err != nil {
		t.Fatalf("has `127.0.0.1`: unexpected error %s", err)
	}
	if ok {
		t.Fatalf("expected there is not ip `127.0.0.1` in list")
	}
}

func TestIPList_Count(t *testing.T) {
	list := NewIPList()

	var ok bool
	var cnt int
	var err error

	cnt, err = list.Count()
	assertCountResult(t, 0, cnt, err, "count")

	ok, err = list.Add(entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add ip `127.0.0.1`")

	cnt, err = list.Count()
	assertCountResult(t, 1, cnt, err, "count")

	ok, err = list.Add(entities.IP("127.0.0.2"))
	assertOkResult(t, ok, err, "add ip `127.0.0.2`")

	cnt, err = list.Count()
	assertCountResult(t, 2, cnt, err, "count")

	ok, err = list.Delete(entities.IP("127.0.0.2"))
	assertOkResult(t, ok, err, "delete `127.0.0.1`")

	cnt, err = list.Count()
	assertCountResult(t, 1, cnt, err, "count")

	ok, err = list.Delete(entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "delete `127.0.0.1`")

	cnt, err = list.Count()
	assertCountResult(t, 0, cnt, err, "count")

}

func assertOkResult(t *testing.T, ok bool, err error, id string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %s", id, err)
	}
	if !ok {
		t.Fatalf("%s: expected be successfull", id)
	}
}

func assertCountResult(t *testing.T, expected int, count int, err error, id string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %s", id, err)
	}
	if count != expected {
		t.Fatalf("%s: unexpected count %d instreadof %d", id, count, expected)
	}
}
