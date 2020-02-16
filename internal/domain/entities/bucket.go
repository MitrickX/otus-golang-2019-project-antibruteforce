package entities

import (
	"context"
	"time"
)

// Bucket interface, bucket is abstraction that implements rate limit conception
type Bucket interface {
	// IsConform checks is packet conform bucket
	IsConform(t time.Time) bool
	// IsActive checks if bucket active to this time
	IsActive(t time.Time) bool
}

// BucketStorage interface, interface for data struct where we keep buckets by keys
type BucketStorage interface {
	// Add bucket into storage by key
	Add(ctx context.Context, bucket Bucket, key interface{}) error
	// Delete bucket from storage by key
	Delete(ctx context.Context, key interface{}) error
	// Get bucket from storage by key
	Get(ctx context.Context, key interface{}) (Bucket, error)
	// Has storage bucket by key?
	Has(ctx context.Context, key interface{}) (bool, error)
	// Count of total number of buckets in storage
	Count(ctx context.Context) (int, error)
	// ClearNotActive clear not active to current time buckets from storage
	ClearNotActive(ctx context.Context, t time.Time) (int, error)
}
