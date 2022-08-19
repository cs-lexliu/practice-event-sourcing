package cusecase

import (
	"context"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/centity"
)

type Repository[t centity.AggregateRoot] interface {
	Save(ctx context.Context, entity t) error
	FindByID(ctx context.Context, id string) (t, error)
}
