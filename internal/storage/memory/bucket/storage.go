package bucket

import (
	"context"
	"sync"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities/bucket"
)

type Storage struct {
	m  map[interface{}]bucket.Bucket
	mx sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		m:  make(map[interface{}]bucket.Bucket),
		mx: sync.RWMutex{},
	}
}

func (s *Storage) Add(ctx context.Context, bucket bucket.Bucket, key interface{}) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.m[key] = bucket
	return nil
}

func (s *Storage) Delete(ctx context.Context, key interface{}) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.m, key)
	return nil
}

func (s *Storage) Get(ctx context.Context, key interface{}) (bucket.Bucket, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	bucket, has := s.m[key]
	if !has {
		return nil, nil
	}
	return bucket, nil
}

func (s *Storage) Has(ctx context.Context, key interface{}) (bool, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	_, has := s.m[key]
	return has, nil
}

func (s *Storage) Count(context.Context) (int, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return len(s.m), nil
}
