package usecase

import (
	"context"
	"fmt"
)

type RollOneBall struct {
	repository BowlingGameRepository
}

func NewRollOneBallUseCase(repository BowlingGameRepository) *RollOneBall {
	return &RollOneBall{
		repository: repository,
	}
}

func (u *RollOneBall) Execute(ctx context.Context, id string, hit int) error {
	b, err := u.repository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repository find by id: %w", err)
	}
	if err := b.Roll(hit); err != nil {
		return fmt.Errorf("roll: %w", err)
	}
	if err := u.repository.Save(ctx, b); err != nil {
		return fmt.Errorf("repository save: %w", err)
	}
	return nil
}
