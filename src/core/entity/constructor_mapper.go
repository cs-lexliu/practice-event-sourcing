package entity

import (
	"context"
	"fmt"
	"reflect"
)

var mapper = map[string]Constructor{}

type Constructor func([]DomainEvent) interface{}

func RegisterConstructor(ctx context.Context, aggregateRoot AggregateRoot, constructor Constructor) {
	mapper[aggregateRoot.Category()] = constructor
}

func GetConstructor[t AggregateRoot](ctx context.Context, aggregateRoot *t) (Constructor, error) {
	obj := reflect.New(reflect.TypeOf(aggregateRoot).Elem().Elem())
	category := obj.MethodByName("Category").Call(nil)[0].String()
	constructor, ok := mapper[category]
	if !ok {
		return nil, fmt.Errorf("constructor not found")
	}
	return constructor, nil
}
