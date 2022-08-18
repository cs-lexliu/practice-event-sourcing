package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/cs-lexliu/practice-event-sourcing/src/adpater/inmem/pubsub"
	"github.com/cs-lexliu/practice-event-sourcing/src/adpater/inmem/statestore"
	"github.com/cs-lexliu/practice-event-sourcing/src/bowling_game/entity"
	"github.com/cs-lexliu/practice-event-sourcing/src/bowling_game/usecase/port/in"
	core "github.com/cs-lexliu/practice-event-sourcing/src/core/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateBowlingGameUseCase(t *testing.T) {
	ctx := context.Background()

	eventBus := pubsub.NewDomainEventBus()
	listener := &fakeCreateBowlingGameListener{}
	subscriber := pubsub.NewSubscriber("event-listener", listener.Notifier)
	eventBus.Register(ctx, subscriber)

	repository := statestore.NewStateRepositoryAdapter[*entity.BowlingGame](eventBus)
	createBowlingGameUseCase := NewCreateBowlingGameService(repository)
	input := in.CreateBowlingGameInput{
		GameID: uuid.New().String(),
	}

	output, err := createBowlingGameUseCase.execute(ctx, input)
	assert.NoError(t, err)
	_, err = repository.FindByID(ctx, output)
	assert.NoError(t, err)

	assert.Eventually(t, func() bool {
		return assert.Equal(t, 1, listener.count)
	}, 5*time.Millisecond, 1*time.Millisecond)
}

type fakeCreateBowlingGameListener struct {
	count int
}

func (l *fakeCreateBowlingGameListener) Notifier(event core.DomainEvent) {
	l.count++
}
