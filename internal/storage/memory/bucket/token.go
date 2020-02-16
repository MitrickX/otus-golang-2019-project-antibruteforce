package bucket

import (
	"sync"
	"time"
)

// TokenBucket structure, implements Bucket interface and "Buckets of token" data structure
type TokenBucket struct {
	// count of tokens in bucket
	count uint

	// duration in which release one token (increment count)
	duration time.Duration

	// max possible count of tokens in bucket
	limit uint

	// last conform unix nano timestamp
	lastConformTime int64

	// last time when constructor or isConform had been called
	lastActiveTime int64

	// max duration from lastActiveTime to current time that could signal that bucket already is not active
	activeDuration time.Duration

	// avoid race
	mx sync.Mutex
}

// NewTokenBucket construct new TokenBucket structure by limit and duration
func NewTokenBucket(t time.Time, limit uint, duration time.Duration, activeDuration time.Duration) *TokenBucket {
	return &TokenBucket{
		count:          limit,
		duration:       duration,
		limit:          limit,
		lastActiveTime: t.UnixNano(),
		activeDuration: activeDuration,
	}
}

// NewTokenBucketByLimitInMinute construct new TokenBucket structure by rate in minute
func NewTokenBucketByLimitInMinute(t time.Time, limit uint, activeDuration time.Duration) *TokenBucket {
	d := time.Minute / time.Duration(limit)
	return NewTokenBucket(t, limit, d, activeDuration)
}

// IsConform checks is packet conform bucket
func (b *TokenBucket) IsConform(t time.Time) bool {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.lastActiveTime = t.UnixNano()

	b.releaseTokens(t)

	if b.count > 0 {
		b.count--                            // conform packet, consume one token
		b.lastConformTime = b.lastActiveTime // last conform time

		return true
	}

	return false
}

// IsActive checks if bucket active to this time
func (b *TokenBucket) IsActive(t time.Time) bool {
	b.mx.Lock()
	defer b.mx.Unlock()

	// we must release tokens because since last IsConform could be pass enough time to full bucket
	b.releaseTokens(t)

	// if bucket is not full bucket is active yet
	// bucket is not active if it is full for a while and no tokens are requested for a while
	if b.count < b.limit {
		return true
	}

	tms := t.UnixNano()

	elapsed := tms - b.lastActiveTime
	if elapsed <= 0 {
		elapsed = 0
	}

	return elapsed <= int64(b.activeDuration)
}

// inner helper, no need lock mutex cause already locked
func (b *TokenBucket) releaseTokens(t time.Time) {
	tms := t.UnixNano()

	elapsed := tms - b.lastConformTime
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
