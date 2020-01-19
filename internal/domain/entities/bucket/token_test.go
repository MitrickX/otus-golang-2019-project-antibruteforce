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

// Test that first packet is conformed
func TestTokenBucket_ConformFirst(t *testing.T) {
	N := 10
	bucket := NewTokenBucketByLimitInMinute(uint(N))
	conformed := bucket.Conform(time.Now())
	if !conformed {
		t.Fatalf("first packet must be conformed always")
	}
}

// Test that now sent N packets are conformed but N + 1 not anymore
func TestTokenBucket_ConformAllAtOnce(t *testing.T) {
	testConformWithTimeoutStep(t, 0, "without timeout")
}

// Test that sent N packets with timeout step 50ms are conformed but N + 1 not anymore
func TestTokenBucket_ConformWithSmallTimeoutStep(t *testing.T) {
	testConformWithTimeoutStep(t, 50*time.Millisecond, "timeout step is 50ms")
}

// Test that sent N packets with timeout step 100ms are conformed but N + 1 not anymore
func TestTokenBucket_ConformWithMiddleTimeoutStep(t *testing.T) {
	testConformWithTimeoutStep(t, 100*time.Millisecond, "timeout step is 100ms")
}

// Test that sent N packets with timeout step 5s are conformed but N + 1 not anymore
func TestTokenBucket_ConformWithBigTimeoutStep(t *testing.T) {
	testConformWithTimeoutStep(t, 5*time.Second, "timeout step is 5s")
}

// Test that send N packets with some timeout step are conformed but N + 1 not anymore
func testConformWithTimeoutStep(t *testing.T, timeout time.Duration, prefix string) {
	now := time.Unix(0, 0)
	N := 10
	bucket := NewTokenBucketByLimitInMinute(uint(N))
	for i := 0; i < N; i++ {
		conformed := bucket.Conform(now.Add(time.Duration(i) * timeout))
		if !conformed {
			t.Fatalf("%s: packet %d must be conformed", prefix, i+1)
		}
	}

	// (N + 1)th must non conformed
	conformed := bucket.Conform(now.Add(time.Duration(N) * timeout))
	if conformed {
		t.Fatalf("%s: %d packet must be not conformed, cause bucket is overflowing", prefix, N+1)
	}
}

// Test that now sent N packets are conformed but N + 1 not anymore
// Sending packet is concurrent
func TestTokenBucket_ConformAllAtOnceConcurrently(t *testing.T) {
	now := time.Now()

	N := 10
	bucket := NewTokenBucketByLimitInMinute(uint(N))

	// counter - number'
	conformedCount := newCounter(0)
	nonConformedCount := newCounter(0)

	wg := sync.WaitGroup{}

	for i := 0; i < N+1; i++ {
		wg.Add(1)
		go func(t time.Time) {
			conformed := bucket.Conform(t)
			if conformed {
				conformedCount.inc()
			} else {
				nonConformedCount.inc()
			}
			wg.Done()
		}(now)
	}

	wg.Wait()

	if conformedCount.val() != N {
		t.Errorf("unexpected count of conformed packets %d instreadof %d", conformedCount.val(), N)
	}

	if nonConformedCount.val() != 1 {
		t.Errorf("unexpected count of nonconformed packets %d instreadof %d", conformedCount.val(), N)
	}

}

func TestTokenBucket_ConformAfterNeedDuration(t *testing.T) {
	now := time.Now()
	N := 10
	bucket := NewTokenBucketByLimitInMinute(uint(N))
	for i := 0; i < N; i++ {
		conformed := bucket.Conform(now)
		if !conformed {
			t.Fatalf("packet %d must be conformed", i+1)
		}
	}
	nextTime := now.Add(6 * time.Second)
	conformed := bucket.Conform(nextTime)
	if !conformed {
		t.Fatalf("packet must be conformed after timeout of %s", 6*time.Second)
	}
}
