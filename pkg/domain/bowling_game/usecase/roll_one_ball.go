package usecase

import (
	"context"
	"fmt"

	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/usecase/port/in"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/usecase/port/out"
)

type RollOneBall struct {
	repository out.BowlingGameRepository
}

func NewRollOneBallUseCase(repository out.BowlingGameRepository) *RollOneBall {
	return &RollOneBall{
		repository: repository,
	}
}

func (u *RollOneBall) Execute(ctx context.Context, input in.RollOneBallInput) error {
	b, err := u.repository.FindByID(ctx, input.BowlingGameID)
	if err != nil {
		return fmt.Errorf("repository find by id: %w", err)
	}
	if err := b.Roll(input.Hit); err != nil {
		return fmt.Errorf("roll: %w", err)
	}
	if err := u.repository.Save(ctx, b); err != nil {
		return fmt.Errorf("repository save: %w", err)
	}
	return nil
}
