package bucket

import (
	"testing"
	"time"
)

type packet struct {
	arrival int64
	conform bool
}

func TestLeakyBucket_Conform(t *testing.T) {
	//nolint:gomnd
	increment := time.Duration(4)

	//nolint:gomnd
	limit := time.Duration(6)

	bucket := NewLeakyBucket(increment, limit)

	expected := getExpectedResult()

	for i := 0; i < len(expected); i++ {
		tm := time.Unix(0, expected[i].arrival)

		conform := bucket.IsConform(tm)
		if conform != expected[i].conform {
			t.Errorf("unexpected that packet %d with time arrival %d has conform status(%t)", i, expected[i].arrival, conform)
		}
	}
}

func getExpectedResult() []packet {
	return []packet{
		{
			//nolint:gomnd
			arrival: 2,
			conform: true,
		},
		{
			//nolint:gomnd
			arrival: 3,
			conform: true,
		},
		{
			//nolint:gomnd
			arrival: 6,
			conform: true,
		},
		{
			//nolint:gomnd
			arrival: 9,
			conform: true,
		},
		{
			//nolint:gomnd
			arrival: 11,
			conform: false,
		},
		{
			//nolint:gomnd
			arrival: 16,
			conform: true,
		},
		{
			//nolint:gomnd
			arrival: 23,
			conform: true,
		},
		{
			//nolint:gomnd
			arrival: 24,
			conform: true,
		},
		{
			//nolint:gomnd
			arrival: 25,
			conform: true,
		},
		{
			//nolint:gomnd
			arrival: 26,
			conform: false,
		},
		{
			//nolint:gomnd
			arrival: 30,
			conform: true,
		},
	}
}
