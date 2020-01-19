package bucket

import (
	"fmt"
	"sync"
	"time"
)

type LeakyBucket struct {
	increment time.Duration // I
	limit     time.Duration // L
	content   time.Duration // X
	lct       int64         // LCT as UnixNano
	mx        sync.Mutex    // avoid race
}

func NewLeakyBucket(increment time.Duration, limit time.Duration) *LeakyBucket {
	if DEBUG {
		fmt.Printf("I=%s,L=%s\n", increment, limit)
	}
	return &LeakyBucket{increment: increment, limit: limit, lct: -1}
}

func (b *LeakyBucket) Conform(t time.Time) bool {
	b.mx.Lock()
	defer b.mx.Unlock()

	// for first packet we init LCT by time arrival of first package
	if b.lct == -1 {
		b.lct = t.UnixNano()
	}

	ta := t.UnixNano()

	if DEBUG {
		fmt.Printf("\nX=%s, LCT=%s, T=%s\n", b.content, time.Duration(b.lct), time.Duration(ta))
		fmt.Printf("T-LCT=%s\n", time.Duration(ta-b.lct))
	}

	auxiliary := b.content - time.Duration(ta-b.lct) // X' = X - (ta - LCT)
	if auxiliary <= 0 {
		auxiliary = 0
	}

	if DEBUG {
		fmt.Printf("X'=%s,X'+I=%s\n", auxiliary, auxiliary+b.increment)
		fmt.Printf("X' > L? %t\n", auxiliary > b.limit)
	}

	if auxiliary > b.limit { // X' > L?
		// nonconforming "packet", bucket is overflowing
		return false
	}

	b.content = auxiliary + b.increment // X = X' + I
	b.lct = ta                          // LCT = ta

	// confirming "packet"
	return true
}
