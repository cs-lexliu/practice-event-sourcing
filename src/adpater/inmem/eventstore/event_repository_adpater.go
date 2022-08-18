package eventstore

import (
	"context"
	"fmt"
	"sync"

	"github.com/cs-lexliu/practice-event-sourcing/src/core/entity"
	"github.com/cs-lexliu/practice-event-sourcing/src/core/usecase"
)

// t represents the pointer of the entity
type eventRepository[t entity.AggregateRoot] struct {
	sync.RWMutex
	storage  *storage[entity.DomainEvent]
	eventBus usecase.DomainEventBus
}

func NewEventRepositoryAdapter[t entity.AggregateRoot](eventBus usecase.DomainEventBus) usecase.Repository[t] {
	return &eventRepository[t]{
		storage:  newStorage[entity.DomainEvent](),
		eventBus: eventBus,
	}
}

func (r *eventRepository[t]) Save(ctx context.Context, aggregateRoot t) error {
	r.Lock()
	defer r.Unlock()
	r.storage.Append(ctx, aggregateRoot.ID(), aggregateRoot.DomainEvents())
	r.eventBus.PostAll(ctx, aggregateRoot)
	aggregateRoot.ClearDomainEvents()
	return nil
}

func (r *eventRepository[t]) FindByID(ctx context.Context, id string) (t, error) {
	events, err := r.storage.Get(ctx, id)
	if err != nil {
		return r.nilObj(), fmt.Errorf("storage get: %w", err)
	}
	constructor, err := entity.GetConstructor(ctx, new(t))
	if err != nil {
		return r.nilObj(), fmt.Errorf("get constructor")
	}
	return constructor(events).(t), nil
}

func (r *eventRepository[t]) nilObj() t {
	var nilObj t
	return nilObj
}
