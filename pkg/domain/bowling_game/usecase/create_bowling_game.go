package usecase

import (
	"context"
	"fmt"

	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/entity"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/usecase/port/in"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/usecase/port/out"
)

type CreateBowlingGame struct {
	repository out.BowlingGameRepository
}

func NewCreateBowlingGameUseCase(repository out.BowlingGameRepository) *CreateBowlingGame {
	return &CreateBowlingGame{
		repository: repository,
	}
}

func (u *CreateBowlingGame) Execute(ctx context.Context, input in.CreateBowlingGameInput) error {
	b, err := entity.Create(input.BowlingGameID)
	if err != nil {
		return fmt.Errorf("new bowling game: %w", err)
	}
	if err := u.repository.Save(ctx, b); err != nil {
		return fmt.Errorf("repository save: %w", err)
	}
	return nil
}
