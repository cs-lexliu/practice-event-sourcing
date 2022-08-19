package usecase

import (
	"context"
	"fmt"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/cusecase"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/entity"
)

type BowlingGameRepository cusecase.Repository[*entity.BowlingGame]

type CreateBowlingGame struct {
	repository BowlingGameRepository
}

func NewCreateBowlingGameUseCase(repository BowlingGameRepository) *CreateBowlingGame {
	return &CreateBowlingGame{
		repository: repository,
	}
}

func (u *CreateBowlingGame) Execute(ctx context.Context, id string) error {
	b, err := entity.NewBowlingGame(id)
	if err != nil {
		return fmt.Errorf("new bowling game: %w", err)
	}
	if err := u.repository.Save(ctx, b); err != nil {
		return fmt.Errorf("repository save: %w", err)
	}
	return nil
}
