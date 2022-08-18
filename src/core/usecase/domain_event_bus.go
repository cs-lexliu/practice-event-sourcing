package usecase

import (
	"context"

	"github.com/cs-lexliu/practice-event-sourcing/src/core/entity"
)

type DomainEventBus interface {
	Post(ctx context.Context, domainEvent entity.DomainEvent) error
	PostAll(ctx context.Context, aggregateRoot entity.AggregateRoot) error
	Register(ctx context.Context, listener any) error
	Unregister(ctx context.Context, listener any) error
}
