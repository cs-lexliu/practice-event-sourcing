package klmap

import (
	"context"
	"fmt"
	"sync"
)

// KLMap represents a thread safe store map a key to list of data.
type KLMap[t any] struct {
	mu    sync.RWMutex
	store map[string][]t
}

func NewKLMap[t any]() *KLMap[t] {
	return &KLMap[t]{
		store: map[string][]t{},
	}
}

func (s *KLMap[t]) Append(ctx context.Context, id string, objs []t) {
	s.mu.Lock()
	defer s.mu.Unlock()
	eventList := s.store[id]
	s.store[id] = append(eventList, objs...)
}

func (s *KLMap[t]) Get(ctx context.Context, id string) ([]t, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	eventList, found := s.store[id]
	if !found {
		return nil, fmt.Errorf("not found")
	}
	return eventList, nil
}
