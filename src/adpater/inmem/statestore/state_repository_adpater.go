package statestore

import (
	"context"
	"fmt"
	"sync"

	"github.com/cs-lexliu/practice-event-sourcing/src/core/entity"
	"github.com/cs-lexliu/practice-event-sourcing/src/core/usecase"
)

type stateRepository[t entity.AggregateRoot] struct {
	sync.RWMutex
	storage  *storage[t]
	eventBus usecase.DomainEventBus
}

func NewStateRepositoryAdapter[t entity.AggregateRoot](eventBus usecase.DomainEventBus) usecase.Repository[t] {
	return &stateRepository[t]{
		storage:  newStorage[t](),
		eventBus: eventBus,
	}
}

func (r *stateRepository[t]) Save(ctx context.Context, aggregateRoot t) error {
	r.Lock()
	defer r.Unlock()
	r.storage.Set(ctx, aggregateRoot.ID(), aggregateRoot)
	r.eventBus.PostAll(ctx, aggregateRoot)
	aggregateRoot.ClearDomainEvents()
	return nil
}

func (r *stateRepository[t]) FindByID(ctx context.Context, id string) (t, error) {
	aggregateRoot, err := r.storage.Get(ctx, id)
	if err != nil {
		var empty t
		return empty, fmt.Errorf("statestore get: %w", err)
	}
	return aggregateRoot, nil
}
