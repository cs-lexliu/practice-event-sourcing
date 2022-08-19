package pubsub

import (
	"context"
	"sync"
)

// Subscriber represents the message consumer listening the message from hub.
type Subscriber[T any] struct {
	sync.Mutex
	name       string
	notifierFn func(T)
	handler    chan message
	quit       chan struct{}
}

func NewSubscriber[T any](name string, notifierFn func(T)) *Subscriber[T] {
	return &Subscriber[T]{
		name:       name,
		notifierFn: notifierFn,
		handler:    make(chan message, 100),
		quit:       make(chan struct{}),
	}
}

func (s *Subscriber[T]) run(ctx context.Context) {
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

func (s *Subscriber[T]) publish(ctx context.Context, msg message) {
	select {
	case <-ctx.Done():
		return
	case s.handler <- msg:
	default:
	}
}
