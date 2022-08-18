package usecase

import (
	"context"
)

type Repository[t any] interface {
	Save(ctx context.Context, entity t) error
	FindByID(ctx context.Context, id string) (t, error)
}
