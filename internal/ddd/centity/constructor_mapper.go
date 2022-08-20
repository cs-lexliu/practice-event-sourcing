package centity

import (
	"context"
	"fmt"
	"reflect"

	"github.com/cs-lexliu/practice-event-sourcing/pkg/utils/inmem/kvmap"
)

var cm *constructorMap

func init() {
	cm = newConstructorMap()
}

type Constructor func([]DomainEvent) (AggregateRoot, error)

type constructorMap struct {
	store *kvmap.KVMap[Constructor]
}

func newConstructorMap() *constructorMap {
	return &constructorMap{
		store: kvmap.NewKVMap[Constructor](),
	}
}

func RegisterConstructor(aggregateRoot AggregateRoot, constructor Constructor) {
	cm.store.Set(context.Background(), aggregateRoot.Category(), constructor)
}

func GetConstructor(ctx context.Context, aggregateRoot AggregateRoot) (Constructor, error) {
	constructor, err := cm.store.Get(ctx, getAggregateRootCategory(aggregateRoot))
	if err != nil {
		return nil, fmt.Errorf("constructor get: %w", err)
	}

	return constructor, nil
}

func getAggregateRootCategory(aggregateRoot AggregateRoot) string {
	obj := reflect.New(reflect.TypeOf(aggregateRoot).Elem())
	return obj.MethodByName("Category").Call(nil)[0].String()
}
