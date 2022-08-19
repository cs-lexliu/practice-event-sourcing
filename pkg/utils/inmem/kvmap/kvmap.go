package kvmap

import (
	"context"
	"fmt"
	"sync"
)

// KVMap represents a thread safe key value map.
type KVMap[t any] struct {
	mu    sync.RWMutex
	store map[string]t
}

func NewKVMap[t any]() *KVMap[t] {
	return &KVMap[t]{
		store: map[string]t{},
	}
}

func (s *KVMap[t]) Set(ctx context.Context, k string, v t) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[k] = v
}

func (s *KVMap[t]) Get(ctx context.Context, k string) (t, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, found := s.store[k]
	if !found {
		var empty t
		return empty, fmt.Errorf("not found")
	}
	return v, nil
}
