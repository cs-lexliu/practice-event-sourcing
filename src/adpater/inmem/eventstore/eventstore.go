package eventstore

import (
	"context"
	"fmt"
	"sync"
)

type storage[t any] struct {
	sync.RWMutex
	memory map[string][]t
}

func newStorage[t any]() *storage[t] {
	return &storage[t]{
		memory: map[string][]t{},
	}
}

func (s *storage[t]) Append(ctx context.Context, id string, objs []t) {
	s.Lock()
	defer s.Unlock()
	eventList := s.memory[id]
	s.memory[id] = append(eventList, objs...)
}

func (s *storage[t]) Get(ctx context.Context, id string) ([]t, error) {
	s.RLock()
	defer s.RUnlock()
	eventList, found := s.memory[id]
	if !found {
		return nil, fmt.Errorf("not found")
	}
	return eventList, nil
}
