package statestore

import (
	"context"
	"fmt"
	"sync"
)

type storage[t any] struct {
	sync.RWMutex
	entities map[string]t
}

func newStorage[t any]() *storage[t] {
	return &storage[t]{
		entities: map[string]t{},
	}
}

func (s *storage[t]) Set(ctx context.Context, k string, v t) {
	s.Lock()
	defer s.Unlock()
	s.entities[k] = v
}

func (s *storage[t]) Get(ctx context.Context, k string) (t, error) {
	s.RLock()
	defer s.RUnlock()
	v, found := s.entities[k]
	if !found {
		var empty t
		return empty, fmt.Errorf("not found")
	}
	return v, nil
}
