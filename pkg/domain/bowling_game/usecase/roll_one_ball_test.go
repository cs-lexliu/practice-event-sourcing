package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/adpater"
	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/centity"
	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/cusecase"
	"github.com/cs-lexliu/practice-event-sourcing/pkg/domain/bowling_game/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type RollOneBallSuite struct {
	suite.Suite

	repository    BowlingGameRepository
	eventBus      cusecase.DomainEventBus
	eventListener *RollOneBallEventListener
}

func TestRollOneBallSuite(t *testing.T) {
	suite.Run(t, new(RollOneBallSuite))
}

func (s *RollOneBallSuite) SetupTest() {
	s.eventBus = adpater.NewDomainEventBus()
	s.eventListener = &RollOneBallEventListener{}
	s.eventBus.Register(context.Background(), "RollOneBallSuite", s.eventListener)

	//s.repository = adpater.NewEventRepositoryAdapter[*entity.BowlingGame](s.eventBus)
	s.repository = adpater.NewStateRepositoryAdapter[*entity.BowlingGame](s.eventBus)
}

func (s *RollOneBallSuite) TestRollOneBallShouldIncreaseScoresAndDecreasePins() {
	id := uuid.New().String()
	s.createBowlingGame(id)

	u := NewRollOneBallUseCase(s.repository)
	u.Execute(context.Background(), id, 1)

	s.Eventually(
		func() bool {
			return s.Equal(1, s.eventListener.rolledCount)
		},
		10*time.Millisecond,
		1*time.Millisecond,
	)

	got, err := s.repository.FindByID(context.Background(), id)
	s.NoError(err)
	s.Equal(1, got.Scores())
	s.Equal(9, got.Pins())
}

type RollOneBallEventListener struct {
	createdCount int
	rolledCount  int
}

func (l *RollOneBallEventListener) Consume(event centity.DomainEvent) {
	switch interface{}(event).(type) {
	case entity.BowlingGameCreatedEvent:
		l.createdCount++
	case entity.BowlingGameRolledEvent:
		l.rolledCount++
	}
}

func (s *RollOneBallSuite) createBowlingGame(id string) {
	u := NewCreateBowlingGameUseCase(s.repository)
	u.Execute(context.Background(), id)
}
