package adpater

import (
	"context"
	"fmt"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/centity"
	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/cusecase"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/utils/inmem/pubsub"
)

type domainEventBusAdapter struct {
	bus *pubsub.Hub[centity.DomainEvent]
}

func NewDomainEventBus() cusecase.DomainEventBus {
	return &domainEventBusAdapter{
		bus: pubsub.NewPubSub[centity.DomainEvent](),
	}
}

func (d domainEventBusAdapter) Post(ctx context.Context, domainEvent centity.DomainEvent) error {
	if err := d.bus.Publish(ctx, domainEvent); err != nil {
		return fmt.Errorf("bus publish event: %w", err)
	}
	return nil
}

func (d domainEventBusAdapter) PostAll(ctx context.Context, aggregateRoot centity.AggregateRoot) error {
	for _, e := range aggregateRoot.DomainEvents() {
		if err := d.bus.Publish(ctx, e); err != nil {
			return fmt.Errorf("bus publish event: %w", err)
		}
	}
	return nil
}

func (d domainEventBusAdapter) Register(ctx context.Context, name string, listener cusecase.DomainEventListener) error {
	sub := pubsub.NewSubscriber(name, listener.Consume)
	if err := d.bus.Subscribe(ctx, name, sub); err != nil {
		return fmt.Errorf("bus subscribe listener: %w", err)
	}
	return nil
}

func (d domainEventBusAdapter) Unregister(ctx context.Context, name string) error {
	if err := d.bus.Unsubscribe(ctx, name); err != nil {
		return fmt.Errorf("bus unsubscribe listener: %w", err)
	}
	return nil
}
