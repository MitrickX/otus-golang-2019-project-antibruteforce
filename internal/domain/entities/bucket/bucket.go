package bucket

import (
	"time"
)

type Bucket interface {
	IsConform(t time.Time) bool
}
