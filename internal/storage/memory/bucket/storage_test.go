package bucket

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
)

type emptyBucket struct {
	id int64
}

func newEmptyBucket() emptyBucket {
	return emptyBucket{
		id: rand.Int63(),
	}
}

func (e emptyBucket) IsConform(t time.Time) bool {
	return true
}

func TestStorage_Add(t *testing.T) {
	storage := NewStorage()

	var err error

	err = storage.Add(context.Background(), newEmptyBucket(), entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.1`")

	err = storage.Add(context.Background(), newEmptyBucket(), entities.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.2`")

	err = storage.Add(context.Background(), newEmptyBucket(), entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.1` (2nd time)")

	cnt, err := storage.Count(context.Background())
	assertCountResult(t, 2, cnt, err, "count after adds")
}

func TestStorage_Delete(t *testing.T) {
	storage := NewStorage()

	var err error

	err = storage.Add(context.Background(), newEmptyBucket(), entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.1`")

	err = storage.Add(context.Background(), newEmptyBucket(), entities.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.2`")

	// Delete not existing
	err = storage.Delete(context.Background(), entities.IP("127.0.0.3"))
	assertNotErrorResult(t, err, "delete bucket for `127.0.0.3`")

	err = storage.Delete(context.Background(), entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "delete bucket for `127.0.0.1`")

	cnt, err := storage.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count")
}

func TestStorage_Get(t *testing.T) {
	storage := NewStorage()

	var bucket entities.Bucket
	var err error

	bucket, err = storage.Get(context.Background(), entities.IP("127.0.0.1"))
	assertOkBucketGetResult(t, nil, bucket, err, "get bucket from empty storage")

	expectedBucket1 := newEmptyBucket()
	expectedBucket2 := newEmptyBucket()

	err = storage.Add(context.Background(), expectedBucket1, entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.1`")

	err = storage.Add(context.Background(), expectedBucket2, entities.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.1`")

	bucket1, err := storage.Get(context.Background(), entities.IP("127.0.0.1"))
	assertOkBucketGetResult(t, expectedBucket1, bucket1, err, "get bucket by id `127.0.0.1`")

	bucket2, err := storage.Get(context.Background(), entities.IP("127.0.0.2"))
	assertOkBucketGetResult(t, expectedBucket2, bucket2, err, "get bucket by id `127.0.0.2`")
}

func TestStorage_Has(t *testing.T) {
	storage := NewStorage()

	var ok bool
	var err error

	err = storage.Add(context.Background(), newEmptyBucket(), entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.1`")

	ok, err = storage.Has(context.Background(), entities.IP("127.0.0.1"))
	assertOkResult(t, ok, err, "has bucket for `127.0.0.1`")

	err = storage.Delete(context.Background(), entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "delete bucket for ip `127.0.0.1`")

	ok, err = storage.Has(context.Background(), entities.IP("127.0.0.1"))
	if err != nil {
		t.Fatalf("has bucket for `127.0.0.1`: unexpected error %s", err)
	}
	if ok {
		t.Fatalf("expected there is not bucket for ip `127.0.0.1` in storage")
	}
}

func TestStorage_Count(t *testing.T) {
	storage := NewStorage()

	var cnt int
	var err error

	cnt, err = storage.Count(context.Background())
	assertCountResult(t, 0, cnt, err, "count")

	err = storage.Add(context.Background(), newEmptyBucket(), entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.1`")

	cnt, err = storage.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count")

	err = storage.Add(context.Background(), newEmptyBucket(), entities.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "add bucket for ip `127.0.0.2`")

	cnt, err = storage.Count(context.Background())
	assertCountResult(t, 2, cnt, err, "count")

	err = storage.Delete(context.Background(), entities.IP("127.0.0.2"))
	assertNotErrorResult(t, err, "delete bucket for `127.0.0.1`")

	cnt, err = storage.Count(context.Background())
	assertCountResult(t, 1, cnt, err, "count")

	err = storage.Delete(context.Background(), entities.IP("127.0.0.1"))
	assertNotErrorResult(t, err, "delete bucket for `127.0.0.1`")

	cnt, err = storage.Count(context.Background())
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
		t.Fatalf("%s: expected be successful", prefix)
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

func assertOkBucketGetResult(t *testing.T, expected entities.Bucket, test entities.Bucket, err error, prefix string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %s", prefix, err)
	}
	if expected != test {
		t.Fatalf("%s: expected that test `%+v` be equals to `%+v`", prefix, test, expected)
	}
}
