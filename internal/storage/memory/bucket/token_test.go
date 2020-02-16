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
	bucket := NewTokenBucketByLimitInMinute(time.Now(), uint(N), time.Minute)

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
	//nolint:gomnd
	smallTimeout := 50 * time.Millisecond
	testConformWithTimeoutStep(t, smallTimeout, "timeout step is 50ms")
}

// Test that sent N packets with timeout step 100ms are conform but N + 1 not anymore
func TestTokenBucket_ConformWithMiddleTimeoutStep(t *testing.T) {
	//nolint:gomnd
	middleTimeout := 100 * time.Millisecond
	testConformWithTimeoutStep(t, middleTimeout, "timeout step is 100ms")
}

// Test that sent N packets with timeout step 5s are conform but N + 1 not anymore
func TestTokenBucket_ConformWithBigTimeoutStep(t *testing.T) {
	//nolint:gomnd
	bigTimeout := 5 * time.Second
	testConformWithTimeoutStep(t, bigTimeout, "timeout step is 5s")
}

// Test that send N packets with some timeout step are conform but N + 1 not anymore
func testConformWithTimeoutStep(t *testing.T, timeout time.Duration, prefix string) {
	now := time.Unix(0, 0)
	N := 10
	bucket := NewTokenBucketByLimitInMinute(time.Now(), uint(N), time.Minute)

	for i := 0; i < N; i++ {
		conform := bucket.IsConform(now.Add(time.Duration(i) * timeout))
		if !conform {
			t.Fatalf("%s: packet %d must be conform", prefix, i)
		}
	}

	// (N + 1)th must non conform
	conform := bucket.IsConform(now.Add(time.Duration(N) * timeout))
	if conform {
		//nolint:gomnd
		t.Fatalf("%s: %d packet must be not conform, cause bucket is overflowing", prefix, N+1)
	}
}

// Test that now sent N packets are conform but N + 1 not anymore
// Sending packet is concurrent
func TestTokenBucket_ConformAllAtOnceConcurrently(t *testing.T) {
	now := time.Now()

	N := 10
	bucket := NewTokenBucketByLimitInMinute(time.Now(), uint(N), time.Minute)

	// counter - number'
	conformCount := newCounter(0)
	nonConformCount := newCounter(0)

	wg := sync.WaitGroup{}

	for i := 0; i < N+1; i++ {
		//nolint:gomnd
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

	//nolint:gomnd
	if nonConformCount.val() != 1 {
		t.Errorf("unexpected count of nonconform packets %d instreadof %d", conformCount.val(), N)
	}
}

func TestTokenBucket_ConformAfterNeedDuration(t *testing.T) {
	now := time.Now()
	N := 10
	bucket := NewTokenBucketByLimitInMinute(time.Now(), uint(N), time.Minute)

	for i := 0; i < N; i++ {
		conform := bucket.IsConform(now)
		if !conform {
			t.Fatalf("packet %d must be conform", i)
		}
	}

	//nolint:gomnd
	timeout := 6 * time.Second

	nextTime := now.Add(timeout)

	conform := bucket.IsConform(nextTime)
	if !conform {
		t.Fatalf("packet must be conform after timeout of %s", timeout)
	}
}

func TestTokenBucket_Release(t *testing.T) {
	N := 10
	bucket := NewTokenBucketByLimitInMinute(time.Now(), uint(N), time.Minute)

	nowTime := time.Now()

	getCount := func() uint {
		bucket.mx.Lock()
		defer bucket.mx.Unlock()

		return bucket.count
	}

	// Arrival of packets is one token in minute (slow for this limit rate)
	// So after each arrival we has 9 tokens left, cause on next conform enough tokens will released
	// Bucket always has plenty of tokens

	currentTime := nowTime
	for i := 0; i < N; i++ {
		currentTime = currentTime.Add(time.Minute)

		conform := bucket.IsConform(currentTime)
		if !conform {
			t.Fatalf("packet %d must be conform", i)
		}

		count := getCount()

		//nolint:gomnd
		if count < uint(N-1) {
			t.Fatal("tokens in bucket must be plenty cause of slow arrival rates")
		}
	}
}

func TestTokenBucket_IsActiveWhenIsConformNotCalled(t *testing.T) {
	N := 10
	bucket := NewTokenBucketByLimitInMinute(time.Now(), uint(N), time.Minute)

	if bucket.IsActive(time.Now()) == false {
		t.Fatalf("bucket must be active right after construction")
	}

	nextTime := time.Now().Add(time.Minute).Add(time.Millisecond)
	if bucket.IsActive(nextTime) == true {
		t.Fatalf("bucket must be not active after wait timeout after constructor")
	}
}

func TestTokenBucket_IsActiveRightAfterIsConformCalled(t *testing.T) {
	N := 10

	nowTime := time.Now()
	bucket := NewTokenBucketByLimitInMinute(nowTime, uint(N), time.Minute)

	nextTime := time.Now().Add(time.Minute).Add(time.Millisecond)
	bucket.IsConform(nextTime)

	if bucket.IsActive(nextTime) == false {
		t.Fatalf("bucket must be active right after IsConform called")
	}
}

func TestTokenBucket_IsActiveAfterIsConformCalled(t *testing.T) {
	N := 10

	nowTime := time.Now()
	bucket := NewTokenBucketByLimitInMinute(nowTime, uint(N), time.Minute)

	nextTime := time.Now().Add(time.Minute)
	bucket.IsConform(nextTime)

	// bucket restored, but active timeout is not passed so bucket still active
	nextTime = nextTime.Add(59 * time.Second) //nolint:gomnd

	if bucket.IsActive(nextTime) == false {
		t.Fatalf("bucket must be active, because active timeout is not passed")
	}

	// active timeout is now passed so bucket must be not active
	nextTime = nextTime.Add(time.Second).Add(time.Millisecond)

	if bucket.IsActive(nextTime) == true {
		t.Fatalf("bucket must not be active, because active timeout is passed")
	}
}
