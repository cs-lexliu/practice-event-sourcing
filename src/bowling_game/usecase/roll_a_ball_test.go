package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/cs-lexliu/practice-event-sourcing/src/adpater/inmem/eventstore"
	"github.com/cs-lexliu/practice-event-sourcing/src/adpater/inmem/pubsub"
	"github.com/cs-lexliu/practice-event-sourcing/src/bowling_game/entity"
	"github.com/cs-lexliu/practice-event-sourcing/src/bowling_game/usecase/port/in"
	"github.com/cs-lexliu/practice-event-sourcing/src/bowling_game/usecase/port/out"
	coreentity "github.com/cs-lexliu/practice-event-sourcing/src/core/entity"
	coreusecase "github.com/cs-lexliu/practice-event-sourcing/src/core/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type RollABallUseCaseTestSuite struct {
	suite.Suite

	eventBus   coreusecase.DomainEventBus
	listener   *fakeRollABallListener
	repository out.BowlingGameRepository
}

func TestRollABallUseCaseTestSuite(t *testing.T) {
	suite.Run(t, &RollABallUseCaseTestSuite{})
}

func (s *RollABallUseCaseTestSuite) SetupTest() {
	ctx := context.Background()

	coreentity.RegisterConstructor(entity.BowlingGame{}.Category(), entity.BowlingGameConstuctor)

	s.listener = &fakeRollABallListener{}

	eventBus := pubsub.NewDomainEventBus()
	subscriber := pubsub.NewSubscriber("event-listener", s.listener.Notifier)
	eventBus.Register(ctx, subscriber)
	s.eventBus = eventBus

	//s.repository = statestore.NewStateRepositoryAdapter[*entity.BowlingGame](eventBus)
	s.repository = eventstore.NewEventRepositoryAdapter[*entity.BowlingGame](eventBus)
}

func (s *RollABallUseCaseTestSuite) TestRollABallUseCase() {
	gameID := s.createBowlingGame()

	ctx := context.Background()
	rollABallUseCase := NewRollABallServiceService(s.repository)
	input := in.RollABallInput{
		GameID: gameID,
		Hit:    1,
	}

	gotGameID, err := rollABallUseCase.execute(ctx, input)
	s.NoError(err)
	s.Equal(gameID, gotGameID)
	_, err = s.repository.FindByID(ctx, gotGameID)
	s.NoError(err)

	s.Eventually(func() bool {
		return s.Equal(1, s.listener.count)
	}, 5*time.Millisecond, 1*time.Millisecond)
}

func (s *RollABallUseCaseTestSuite) TestRollTwentyBallHasTwentyScore() {
	gameID := s.createBowlingGame()

	ctx := context.Background()
	rollABallUseCase := NewRollABallServiceService(s.repository)
	input := in.RollABallInput{
		GameID: gameID,
		Hit:    1,
	}

	for i := 0; i < 20; i++ {
		_, _ = rollABallUseCase.execute(ctx, input)
	}
	output, err := s.repository.FindByID(ctx, gameID)
	s.NoError(err)
	s.Equal(20, output.Score())
}

func (s *RollABallUseCaseTestSuite) TestRollTwentyOneBallHasTwentyScore() {
	gameID := s.createBowlingGame()
	s.rollManyBall(gameID, 1, 20)

	ctx := context.Background()
	output, err := s.repository.FindByID(ctx, gameID)
	s.NoError(err)
	s.Equal(20, output.Score())
}

func (s *RollABallUseCaseTestSuite) TestUnableToRollABallAfterGameFinished_TwentyBallWithoutSpareAndStrike() {
	gameID := s.createBowlingGame()
	s.rollManyBall(gameID, 1, 20)

	ctx := context.Background()
	rollABallUseCase := NewRollABallServiceService(s.repository)
	_, err := rollABallUseCase.execute(ctx, in.RollABallInput{
		GameID: gameID,
		Hit:    1,
	})
	s.Error(err)
}

func (s *RollABallUseCaseTestSuite) TestRollABallShouldNotHasNegativeHitValue() {
	gameID := s.createBowlingGame()

	ctx := context.Background()
	rollABallUseCase := NewRollABallServiceService(s.repository)
	input := in.RollABallInput{
		GameID: gameID,
		Hit:    -1,
	}
	_, err := rollABallUseCase.execute(ctx, input)
	s.Error(err)
}

func (s *RollABallUseCaseTestSuite) TestRollABallWithSpareHasNoLeavingPinsAndNotBonusRemainHit() {
	gameID := s.createBowlingGame()
	s.hasSpare(gameID)

	ctx := context.Background()
	output, err := s.repository.FindByID(ctx, gameID)
	s.NoError(err)
	s.Equal(10, output.LeavingPins())
	s.Equal(2, output.BonusRemainHit())
}

func (s *RollABallUseCaseTestSuite) TestRollABallAfterSpareHasBonusScoreAndLeavingPinsIsRestored() {
	gameID := s.createBowlingGame()
	s.hasSpare(gameID)

	ctx := context.Background()
	rollABallUseCase := NewRollABallServiceService(s.repository)
	_, _ = rollABallUseCase.execute(ctx, in.RollABallInput{
		GameID: gameID,
		Hit:    1,
	})
	output, err := s.repository.FindByID(ctx, gameID)
	s.NoError(err)
	s.Equal(12, output.Score())
	s.Equal(9, output.LeavingPins())
}

func (s *RollABallUseCaseTestSuite) TestRollABallWithStrikeHasNoLeavingPinsAndOneBonusRemainHit() {
	gameID := s.createBowlingGame()
	s.hasStrike(gameID)

	ctx := context.Background()
	output, err := s.repository.FindByID(ctx, gameID)
	s.NoError(err)
	s.Equal(10, output.LeavingPins())
	s.Equal(2, output.BonusRemainHit())
}

func (s *RollABallUseCaseTestSuite) TestRollABallAfterStrikeHasTwoBonusScoreAndLeavingPinsIsRestored() {
	gameID := s.createBowlingGame()
	s.hasStrike(gameID)

	ctx := context.Background()
	rollABallUseCase := NewRollABallServiceService(s.repository)
	_, _ = rollABallUseCase.execute(ctx, in.RollABallInput{
		GameID: gameID,
		Hit:    1,
	})
	_, _ = rollABallUseCase.execute(ctx, in.RollABallInput{
		GameID: gameID,
		Hit:    1,
	})
	output, err := s.repository.FindByID(ctx, gameID)
	s.NoError(err)
	s.Equal(14, output.Score())
	s.Equal(10, output.LeavingPins())
}

func (s *RollABallUseCaseTestSuite) TestRollABallWithAPerfectGame() {
	gameID := s.createBowlingGame()
	s.hasManyStrike(gameID, 12)

	ctx := context.Background()
	output, err := s.repository.FindByID(ctx, gameID)
	s.NoError(err)
	s.Equal(300, output.Score())
}

func (s *RollABallUseCaseTestSuite) TestRollABallWithSpareInFinalFrameHasOneBonusFrame() {
	gameID := s.createBowlingGame()
	s.hasManySpare(gameID, 10)

	ctx := context.Background()
	rollABallUseCase := NewRollABallServiceService(s.repository)
	_, err := rollABallUseCase.execute(ctx, in.RollABallInput{
		GameID: gameID,
		Hit:    1,
	})
	s.NoError(err)
}

func (s *RollABallUseCaseTestSuite) createBowlingGame() string {
	createBowlingGameUseCase := NewCreateBowlingGameService(s.repository)
	gameID := uuid.New().String()
	input := in.CreateBowlingGameInput{
		GameID: gameID,
	}

	_, _ = createBowlingGameUseCase.execute(context.Background(), input)
	return gameID
}

func (s *RollABallUseCaseTestSuite) rollManyBall(gameID string, hit int, times int) {
	ctx := context.Background()
	rollABallUseCase := NewRollABallServiceService(s.repository)
	for i := 0; i < times; i++ {
		_, _ = rollABallUseCase.execute(ctx, in.RollABallInput{
			GameID: gameID,
			Hit:    hit,
		})
	}
}

func (s *RollABallUseCaseTestSuite) hasManySpare(gameID string, times int) {
	for i := 0; i < times; i++ {
		s.hasSpare(gameID)
	}
}

func (s *RollABallUseCaseTestSuite) hasSpare(gameID string) {
	ctx := context.Background()
	rollABallUseCase := NewRollABallServiceService(s.repository)
	_, _ = rollABallUseCase.execute(ctx, in.RollABallInput{
		GameID: gameID,
		Hit:    1,
	})
	_, _ = rollABallUseCase.execute(ctx, in.RollABallInput{
		GameID: gameID,
		Hit:    9,
	})
}

func (s *RollABallUseCaseTestSuite) hasManyStrike(gameID string, times int) {
	for i := 0; i < times; i++ {
		s.hasStrike(gameID)
	}
}

func (s *RollABallUseCaseTestSuite) hasStrike(gameID string) {
	ctx := context.Background()
	rollABallUseCase := NewRollABallServiceService(s.repository)
	_, _ = rollABallUseCase.execute(ctx, in.RollABallInput{
		GameID: gameID,
		Hit:    10,
	})
}

type fakeRollABallListener struct {
	count int
}

func (l *fakeRollABallListener) Notifier(domainEvent coreentity.DomainEvent) {
	switch (domainEvent).(type) {
	case entity.BowlingGameRolledABall:
		l.count++
	}
}
