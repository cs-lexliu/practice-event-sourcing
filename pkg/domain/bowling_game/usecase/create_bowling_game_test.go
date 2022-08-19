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

type CreateBowlingGameSuite struct {
	suite.Suite

	repository    BowlingGameRepository
	eventBus      cusecase.DomainEventBus
	eventListener *CreateBowlingGameEventListener
}

func TestCreateBowlingGameSuite(t *testing.T) {
	suite.Run(t, new(CreateBowlingGameSuite))
}

func (s *CreateBowlingGameSuite) SetupTest() {
	s.eventBus = adpater.NewDomainEventBus()
	s.eventListener = &CreateBowlingGameEventListener{}
	s.eventBus.Register(context.Background(), "CreateBowlingGameSuite", s.eventListener)

	s.repository = adpater.NewEventRepositoryAdapter[*entity.BowlingGame](s.eventBus)
}

func (s *CreateBowlingGameSuite) TestCreateBowlingGameGenerateBowlingGameCreatedEvent() {
	u := NewCreateBowlingGameUseCase(s.repository)
	u.Execute(context.Background(), uuid.New().String())
	s.Eventually(
		func() bool {
			return s.Equal(1, s.eventListener.count)
		},
		10*time.Millisecond,
		1*time.Millisecond,
	)
}

type CreateBowlingGameEventListener struct {
	count int
}

func (l *CreateBowlingGameEventListener) Consume(event centity.DomainEvent) {
	switch interface{}(event).(type) {
	case entity.BowlingGameCreatedEvent:
		l.count++
	}
}
