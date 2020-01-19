package bucket

import (
	"time"
)

type Bucket interface {
	Conform(t time.Time) bool
}
