package cusecase

import (
	"context"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/centity"
)

type DomainEventBus interface {
	Post(ctx context.Context, domainEvent centity.DomainEvent) error
	PostAll(ctx context.Context, aggregateRoot centity.AggregateRoot) error
	Register(ctx context.Context, name string, listener DomainEventListener) error
	Unregister(ctx context.Context, name string) error
}

type DomainEventListener interface {
	Consume(event centity.DomainEvent)
}
