package bucket

import (
	"testing"
	"time"
)

func TestLeakyBucket_Conform(t *testing.T) {
	increment := time.Duration(4)
	limit := time.Duration(6)

	bucket := NewLeakyBucket(increment, limit)

	type packet struct {
		arrival   int64
		conformed bool
	}

	expected := []packet{
		{
			arrival:   2,
			conformed: true,
		},
		{
			arrival:   3,
			conformed: true,
		},
		{
			arrival:   6,
			conformed: true,
		},
		{
			arrival:   9,
			conformed: true,
		},
		{
			arrival:   11,
			conformed: false,
		},
		{
			arrival:   16,
			conformed: true,
		},
		{
			arrival:   23,
			conformed: true,
		},
		{
			arrival:   24,
			conformed: true,
		},
		{
			arrival:   25,
			conformed: true,
		},
		{
			arrival:   26,
			conformed: false,
		},
		{
			arrival:   30,
			conformed: true,
		},
	}

	for i := 0; i < len(expected); i++ {
		tm := time.Unix(0, expected[i].arrival)
		conformed := bucket.Conform(tm)
		if conformed != expected[i].conformed {
			t.Errorf("unexpected that packet %d with time arrival %d has conformed status(%t)", i, expected[i].arrival, conformed)
		}
	}

}
