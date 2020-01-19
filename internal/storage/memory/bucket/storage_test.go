package bucket

import (
	"math/rand"
	"testing"
	"time"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities/bucket"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities/ip"
)

type emptyBucket struct {
	id int64
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newEmptyBucket() emptyBucket {
	return emptyBucket{
		id: rand.Int63(),
	}
}

func (e emptyBucket) Conform(t time.Time) bool {
	return true
}

func TestStorage_Add(t *testing.T) {
	storage := NewStorage()

	var ok bool
	var err error

	ok, err = storage.Add(newEmptyBucket(), ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add bucket for ip `127.0.0.1`")

	ok, err = storage.Add(newEmptyBucket(), ip.IP("127.0.0.2"))
	assertOkResult(t, ok, err, "add bucket for ip `127.0.0.2`")

	ok, err = storage.Add(newEmptyBucket(), ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add bucket for ip `127.0.0.1` (2nd time)")

	cnt, err := storage.Count()
	assertCountResult(t, 2, cnt, err, "count")

}

func TestStorage_Delete(t *testing.T) {
	storage := NewStorage()

	var ok bool
	var err error

	ok, err = storage.Add(newEmptyBucket(), ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add bucket for ip `127.0.0.1`")

	ok, err = storage.Add(newEmptyBucket(), ip.IP("127.0.0.2"))
	assertOkResult(t, ok, err, "add bucket for ip `127.0.0.2`")

	// Delete not existing
	ok, err = storage.Delete(ip.IP("127.0.0.3"))
	assertOkResult(t, ok, err, "delete bucket for `127.0.0.3`")

	ok, err = storage.Delete(ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "delete bucket for `127.0.0.1`")

	cnt, err := storage.Count()
	assertCountResult(t, 1, cnt, err, "count")
}

func TestStorage_Get(t *testing.T) {
	storage := NewStorage()

	var bucket bucket.Bucket
	var ok bool
	var err error

	bucket, err = storage.Get(ip.IP("127.0.0.1"))
	assertOkBucketGetResult(t, nil, bucket, err, "get bucket from empty storage")

	expectedBucket1 := newEmptyBucket()
	expectedBucket2 := newEmptyBucket()

	ok, err = storage.Add(expectedBucket1, ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add bucket for ip `127.0.0.1`")

	ok, err = storage.Add(expectedBucket2, ip.IP("127.0.0.2"))
	assertOkResult(t, ok, err, "add bucket for ip `127.0.0.1`")

	bucket1, err := storage.Get(ip.IP("127.0.0.1"))
	assertOkBucketGetResult(t, expectedBucket1, bucket1, err, "get bucket by id `127.0.0.1`")

	bucket2, err := storage.Get(ip.IP("127.0.0.2"))
	assertOkBucketGetResult(t, expectedBucket2, bucket2, err, "get bucket by id `127.0.0.2`")

}

func TestStorage_Has(t *testing.T) {
	storage := NewStorage()

	var ok bool
	var err error

	ok, err = storage.Add(newEmptyBucket(), ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add bucket for ip `127.0.0.1`")

	ok, err = storage.Has(ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "has bucket for `127.0.0.1`")

	ok, err = storage.Delete(ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "delete bucket for ip `127.0.0.1`")

	ok, err = storage.Has(ip.IP("127.0.0.1"))
	if err != nil {
		t.Fatalf("has bucket for `127.0.0.1`: unexpected error %s", err)
	}
	if ok {
		t.Fatalf("expected there is not bucket for ip `127.0.0.1` in storage")
	}
}

func TestStorage_Count(t *testing.T) {
	storage := NewStorage()

	var ok bool
	var cnt int
	var err error

	cnt, err = storage.Count()
	assertCountResult(t, 0, cnt, err, "count")

	ok, err = storage.Add(newEmptyBucket(), ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "add bucket for ip `127.0.0.1`")

	cnt, err = storage.Count()
	assertCountResult(t, 1, cnt, err, "count")

	ok, err = storage.Add(newEmptyBucket(), ip.IP("127.0.0.2"))
	assertOkResult(t, ok, err, "add bucket for ip `127.0.0.2`")

	cnt, err = storage.Count()
	assertCountResult(t, 2, cnt, err, "count")

	ok, err = storage.Delete(ip.IP("127.0.0.2"))
	assertOkResult(t, ok, err, "delete bucket for `127.0.0.1`")

	cnt, err = storage.Count()
	assertCountResult(t, 1, cnt, err, "count")

	ok, err = storage.Delete(ip.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "delete bucket for `127.0.0.1`")

	cnt, err = storage.Count()
	assertCountResult(t, 0, cnt, err, "count")

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

func assertOkBucketGetResult(t *testing.T, expected bucket.Bucket, test bucket.Bucket, err error, prefix string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %s", prefix, err)
	}
	if expected != test {
		t.Fatalf("%s: expected that test `%+v` be equals to `%+v`", prefix, test, expected)
	}
}
