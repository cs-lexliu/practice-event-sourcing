package pubsub

import (
	"context"
	"sync"
)

type hub[T any] struct {
	sync.Mutex
	subs map[*subscriber[T]]struct{}
}

func NewPubSub[T any]() *hub[T] {
	return &hub[T]{
		subs: map[*subscriber[T]]struct{}{},
	}
}

func (h *hub[T]) Subscribe(ctx context.Context, s *subscriber[T]) error {
	h.Lock()
	h.subs[s] = struct{}{}
	h.Unlock()

	go func() {
		select {
		case <-s.quit:
		case <-ctx.Done():
			h.Lock()
			delete(h.subs, s)
			h.Unlock()
		}
	}()

	go s.run(ctx)

	return nil
}

func (h *hub[T]) Unsubscribe(ctx context.Context, s *subscriber[T]) error {
	h.Lock()
	delete(h.subs, s)
	h.Unlock()
	close(s.quit)
	return nil
}

func (h *hub[T]) Subscribers() int {
	return len(h.subs)
}

func (h *hub[T]) Publish(ctx context.Context, msg message) error {
	h.Lock()
	for s := range h.subs {
		s.publish(ctx, msg)
	}
	h.Unlock()

	return nil
}

type message any

type subscriber[T any] struct {
	sync.Mutex
	name       string
	notifierFn func(T)
	handler    chan message
	quit       chan struct{}
}

func NewSubscriber[T any](name string, notifierFn func(T)) *subscriber[T] {
	return &subscriber[T]{
		name:       name,
		notifierFn: notifierFn,
		handler:    make(chan message, 100),
		quit:       make(chan struct{}),
	}
}

func (s *subscriber[T]) run(ctx context.Context) {
	for {
		select {
		case msg := <-s.handler:
			s.notifierFn(msg.(T))
		case <-s.quit:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *subscriber[T]) publish(ctx context.Context, msg message) {
	select {
	case <-ctx.Done():
		return
	case s.handler <- msg:
	default:
	}
}
