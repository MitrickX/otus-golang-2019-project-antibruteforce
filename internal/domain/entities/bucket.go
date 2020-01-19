package entities

import "time"

type Bucket interface {
	Conform(t time.Time) bool
}
