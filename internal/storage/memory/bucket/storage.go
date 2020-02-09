package bucket

import (
	"context"
	"sync"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
)

// Storage type, implements BucketStorage interface and keep buckets in runtime memory
type Storage struct {
	m  map[interface{}]entities.Bucket
	mx sync.RWMutex
}

// NewStorage construct new storage in runtime memory
func NewStorage() *Storage {
	return &Storage{
		m:  make(map[interface{}]entities.Bucket),
		mx: sync.RWMutex{},
	}
}

// Add bucket into storage by key
func (s *Storage) Add(ctx context.Context, bucket entities.Bucket, key interface{}) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.m[key] = bucket

	return nil
}

// Delete bucket from storage by key
func (s *Storage) Delete(ctx context.Context, key interface{}) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.m, key)

	return nil
}

// Get bucket from storage by key
func (s *Storage) Get(ctx context.Context, key interface{}) (entities.Bucket, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	bucket, has := s.m[key]
	if !has {
		return nil, nil
	}

	return bucket, nil
}

// Has storage bucket by key?
func (s *Storage) Has(ctx context.Context, key interface{}) (bool, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	_, has := s.m[key]

	return has, nil
}

// Count of total number of buckets in storage
func (s *Storage) Count(context.Context) (int, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return len(s.m), nil
}
