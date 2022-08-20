package adpater

import (
	"context"
	"fmt"
	"sync"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/centity"
	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/cusecase"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/utils/inmem/klmap"
)

// t represents the pointer of the entity
type eventRepository[t centity.AggregateRoot] struct {
	mu       sync.RWMutex
	storage  *klmap.KLMap[centity.DomainEvent]
	eventBus cusecase.DomainEventBus
}

func NewEventRepositoryAdapter[t centity.AggregateRoot](eventBus cusecase.DomainEventBus) cusecase.Repository[t] {
	return &eventRepository[t]{
		storage:  klmap.NewKLMap[centity.DomainEvent](),
		eventBus: eventBus,
	}
}

func (r *eventRepository[t]) Save(ctx context.Context, aggregateRoot t) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage.Append(ctx, aggregateRoot.ID(), aggregateRoot.DomainEvents())
	if err := r.eventBus.PostAll(ctx, aggregateRoot); err != nil {
		return fmt.Errorf("event bus post all: %w", err)
	}
	aggregateRoot.ClearDomainEvents()
	return nil
}

func (r *eventRepository[t]) FindByID(ctx context.Context, id string) (t, error) {
	events, err := r.storage.Get(ctx, id)
	if err != nil {
		return r.nilObj(), fmt.Errorf("storage get: %w", err)
	}
	constructor, err := centity.GetConstructor(ctx, r.nilObj())
	if err != nil {
		return r.nilObj(), fmt.Errorf("get constructor")
	}
	aggregateRoot, err := constructor(events)
	if err != nil {
		return r.nilObj(), fmt.Errorf("constructor: %w", err)
	}
	return aggregateRoot.(t), nil
}

func (r *eventRepository[t]) nilObj() t {
	var nilObj t
	return nilObj
}
