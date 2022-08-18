package usecase

import (
	"context"
	"fmt"

	"github.com/cs-lexliu/practice-event-sourcing/src/bowling_game/usecase/port/in"
	"github.com/cs-lexliu/practice-event-sourcing/src/bowling_game/usecase/port/out"
)

type RollABallService struct {
	repository out.BowlingGameRepository
}

func NewRollABallServiceService(repository out.BowlingGameRepository) *RollABallService {
	return &RollABallService{
		repository: repository,
	}
}

func (s *RollABallService) execute(ctx context.Context, input in.RollABallInput) (string, error) {
	bowlingGame, err := s.repository.FindByID(ctx, input.GameID)
	if err != nil {
		return "", fmt.Errorf("repository find bowling game by id: %w", err)
	}
	if err := bowlingGame.RollABall(input.Hit); err != nil {
		return "", fmt.Errorf("bowling game roll a ball: %w", err)
	}
	if err := s.repository.Save(ctx, bowlingGame); err != nil {
		return "", fmt.Errorf("repository save bowling game: %w", err)
	}
	return bowlingGame.ID(), nil
}
