package pubsub

import (
	"context"
	"sync"
)

// Hub represents a message broker notifying the message to subscribers.
type Hub[T any] struct {
	sync.Mutex
	subs map[string]*Subscriber[T]
}

func (h *Hub[T]) Subscribe(ctx context.Context, name string, s *Subscriber[T]) error {
	h.Lock()
	h.subs[name] = s
	h.Unlock()

	go func() {
		select {
		case <-s.quit:
		case <-ctx.Done():
			h.Lock()
			delete(h.subs, name)
			h.Unlock()
		}
	}()

	go s.run(ctx)

	return nil
}

func (h *Hub[T]) Unsubscribe(ctx context.Context, name string) error {
	h.Lock()
	s := h.subs[name]
	delete(h.subs, name)
	h.Unlock()
	close(s.quit)
	return nil
}

func (h *Hub[T]) Subscribers() int {
	return len(h.subs)
}

func (h *Hub[T]) Publish(ctx context.Context, msg message) error {
	h.Lock()
	for _, s := range h.subs {
		s.publish(ctx, msg)
	}
	h.Unlock()

	return nil
}
