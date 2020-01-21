package bucket

import (
	"sync"
	"testing"
	"time"
)

type counter struct {
	value int
	mx    sync.RWMutex
}

func newCounter(val int) *counter {
	return &counter{value: val, mx: sync.RWMutex{}}
}

func (c *counter) inc() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.value++
}

func (c *counter) val() int {
	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.value
}

// Test that first packet is conform
func TestTokenBucket_ConformFirst(t *testing.T) {
	N := 10
	bucket := NewTokenBucketByLimitInMinute(uint(N))
	conform := bucket.IsConform(time.Now())
	if !conform {
		t.Fatalf("first packet must be conform always")
	}
}

// Test that now sent N packets are conform but N + 1 not anymore
func TestTokenBucket_ConformAllAtOnce(t *testing.T) {
	testConformWithTimeoutStep(t, 0, "without timeout")
}

// Test that sent N packets with timeout step 50ms are conform but N + 1 not anymore
func TestTokenBucket_ConformWithSmallTimeoutStep(t *testing.T) {
	testConformWithTimeoutStep(t, 50*time.Millisecond, "timeout step is 50ms")
}

// Test that sent N packets with timeout step 100ms are conform but N + 1 not anymore
func TestTokenBucket_ConformWithMiddleTimeoutStep(t *testing.T) {
	testConformWithTimeoutStep(t, 100*time.Millisecond, "timeout step is 100ms")
}

// Test that sent N packets with timeout step 5s are conform but N + 1 not anymore
func TestTokenBucket_ConformWithBigTimeoutStep(t *testing.T) {
	testConformWithTimeoutStep(t, 5*time.Second, "timeout step is 5s")
}

// Test that send N packets with some timeout step are conform but N + 1 not anymore
func testConformWithTimeoutStep(t *testing.T, timeout time.Duration, prefix string) {
	now := time.Unix(0, 0)
	N := 10
	bucket := NewTokenBucketByLimitInMinute(uint(N))
	for i := 0; i < N; i++ {
		conform := bucket.IsConform(now.Add(time.Duration(i) * timeout))
		if !conform {
			t.Fatalf("%s: packet %d must be conform", prefix, i+1)
		}
	}

	// (N + 1)th must non conform
	conform := bucket.IsConform(now.Add(time.Duration(N) * timeout))
	if conform {
		t.Fatalf("%s: %d packet must be not conform, cause bucket is overflowing", prefix, N+1)
	}
}

// Test that now sent N packets are conform but N + 1 not anymore
// Sending packet is concurrent
func TestTokenBucket_ConformAllAtOnceConcurrently(t *testing.T) {
	now := time.Now()

	N := 10
	bucket := NewTokenBucketByLimitInMinute(uint(N))

	// counter - number'
	conformCount := newCounter(0)
	nonConformCount := newCounter(0)

	wg := sync.WaitGroup{}

	for i := 0; i < N+1; i++ {
		wg.Add(1)
		go func(t time.Time) {
			conform := bucket.IsConform(t)
			if conform {
				conformCount.inc()
			} else {
				nonConformCount.inc()
			}
			wg.Done()
		}(now)
	}

	wg.Wait()

	if conformCount.val() != N {
		t.Errorf("unexpected count of conform packets %d instreadof %d", conformCount.val(), N)
	}

	if nonConformCount.val() != 1 {
		t.Errorf("unexpected count of nonconform packets %d instreadof %d", conformCount.val(), N)
	}

}

func TestTokenBucket_ConformAfterNeedDuration(t *testing.T) {
	now := time.Now()
	N := 10
	bucket := NewTokenBucketByLimitInMinute(uint(N))
	for i := 0; i < N; i++ {
		conform := bucket.IsConform(now)
		if !conform {
			t.Fatalf("packet %d must be conform", i+1)
		}
	}
	nextTime := now.Add(6 * time.Second)
	conform := bucket.IsConform(nextTime)
	if !conform {
		t.Fatalf("packet must be conform after timeout of %s", 6*time.Second)
	}
}
