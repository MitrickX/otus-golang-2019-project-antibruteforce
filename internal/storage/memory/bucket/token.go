package bucket

import (
	"sync"
	"time"
)

type TokenBucket struct {
	count    uint          // count of tokens in bucket
	duration time.Duration // duration in which release one token (increment count)
	limit    uint          // max possible count of tokens in bucket
	lct      int64         // last conform unix nano timestamp
	mx       sync.Mutex    // avoid race
}

func NewTokenBucket(limit uint, duration time.Duration) *TokenBucket {
	return &TokenBucket{
		count:    limit,
		duration: duration,
		limit:    limit,
	}
}

func NewTokenBucketByLimitInMinute(limit uint) *TokenBucket {
	d := time.Minute / time.Duration(limit)
	return NewTokenBucket(limit, d)
}

func (b *TokenBucket) IsConform(t time.Time) bool {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.releaseTokens(t)

	if b.count > 0 {
		b.count--            // conform packet, consume one token
		b.lct = t.UnixNano() // last conform time

		return true
	}

	return false
}

// inner helper, no need lock mutex cause already locked
func (b *TokenBucket) releaseTokens(t time.Time) {
	tms := t.UnixNano()

	elapsed := tms - b.lct
	if elapsed <= 0 {
		elapsed = 0
	}

	releaseCount := uint(elapsed / int64(b.duration))

	missingCount := b.limit - b.count
	if releaseCount > missingCount {
		releaseCount = missingCount // to prevent overflow limit of bucket
	}

	b.count += releaseCount
}
