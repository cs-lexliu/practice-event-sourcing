package pubsub

import (
	"context"
	"fmt"

	"github.com/cs-lexliu/practice-event-sourcing/src/core/entity"
	"github.com/cs-lexliu/practice-event-sourcing/src/core/usecase"
)

type domainEventBusAdapter struct {
	bus *hub[entity.DomainEvent]
}

func NewDomainEventBus() usecase.DomainEventBus {
	return &domainEventBusAdapter{
		bus: NewPubSub[entity.DomainEvent](),
	}
}

func (d domainEventBusAdapter) Post(ctx context.Context, domainEvent entity.DomainEvent) error {
	if err := d.bus.Publish(ctx, domainEvent); err != nil {
		return fmt.Errorf("bus publish event: %w", err)
	}
	return nil
}

func (d domainEventBusAdapter) PostAll(ctx context.Context, aggregateRoot entity.AggregateRoot) error {
	for _, e := range aggregateRoot.DomainEvents() {
		if err := d.bus.Publish(ctx, e); err != nil {
			return fmt.Errorf("bus publish event: %w", err)
		}
	}
	return nil
}

func (d domainEventBusAdapter) Register(ctx context.Context, listener any) error {
	if err := d.bus.Subscribe(ctx, listener.(*subscriber[entity.DomainEvent])); err != nil {
		return fmt.Errorf("bus subscribe listener: %w", err)
	}
	return nil
}

func (d domainEventBusAdapter) Unregister(ctx context.Context, listener any) error {
	if err := d.bus.Unsubscribe(ctx, listener.(*subscriber[entity.DomainEvent])); err != nil {
		return fmt.Errorf("bus unsubscribe listener: %w", err)
	}
	return nil
}
