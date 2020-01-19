package bucket

import "context"

type Storage interface {
	// Add bucket into storage by key
	Add(ctx context.Context, bucket Bucket, key interface{}) error
	// Delete bucket from storage by key
	Delete(ctx context.Context, key interface{}) error
	// Get bucket from storage by key
	Get(ctx context.Context, key interface{}) (Bucket, error)
	// Has storage bucket by key?
	Has(ctx context.Context, key interface{}) (bool, error)
	// Total count of buckets in storage
	Count(ctx context.Context) (int, error)
}
