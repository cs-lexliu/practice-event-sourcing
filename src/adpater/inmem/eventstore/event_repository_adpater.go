package eventstore

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/cs-lexliu/practice-event-sourcing/src/core/entity"
	"github.com/cs-lexliu/practice-event-sourcing/src/core/usecase"
)

// t represents the type of entity
type eventRepository[t entity.AggregateRoot] struct {
	sync.RWMutex
	storage  *storage[entity.DomainEvent]
	eventBus usecase.DomainEventBus
}

func NewEventRepositoryAdapter[t entity.AggregateRoot](eventBus usecase.DomainEventBus) usecase.Repository[t] {
	log.Println(reflect.ValueOf(new(t)).Elem().Type().Elem())
	a := reflect.New(reflect.TypeOf(new(t)).Elem().Elem())
	log.Println(a.MethodByName("Category").Call(nil)[0])

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
		var empty t
		return empty, fmt.Errorf("statestore get: %w", err)
	}
	obj := reflect.New(reflect.TypeOf(new(t)).Elem().Elem())
	category := obj.MethodByName("Category").Call(nil)[0]
	constructor := entity.GetConstuctor(category.String())
	return constructor(events).(t), nil
}
