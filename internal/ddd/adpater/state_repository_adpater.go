package adpater

import (
	"context"
	"fmt"
	"sync"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/centity"
	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/cusecase"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/utils/inmem/kvmap"
)

type stateRepository[t centity.AggregateRoot] struct {
	mu       sync.RWMutex
	store    *kvmap.KVMap[t]
	eventBus cusecase.DomainEventBus
}

func NewStateRepositoryAdapter[t centity.AggregateRoot](eventBus cusecase.DomainEventBus) cusecase.Repository[t] {
	return &stateRepository[t]{
		store:    kvmap.NewKVMap[t](),
		eventBus: eventBus,
	}
}

func (r *stateRepository[t]) Save(ctx context.Context, aggregateRoot t) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store.Set(ctx, aggregateRoot.ID(), aggregateRoot)
	r.eventBus.PostAll(ctx, aggregateRoot)
	aggregateRoot.ClearDomainEvents()
	return nil
}

func (r *stateRepository[t]) FindByID(ctx context.Context, id string) (t, error) {
	aggregateRoot, err := r.store.Get(ctx, id)
	if err != nil {
		var empty t
		return empty, fmt.Errorf("statestore get: %w", err)
	}
	return aggregateRoot, nil
}
