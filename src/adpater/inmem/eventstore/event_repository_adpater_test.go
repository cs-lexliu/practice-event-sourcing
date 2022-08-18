package eventstore

import (
	"testing"

	"github.com/cs-lexliu/practice-event-sourcing/src/bowling_game/entity"
)

func TestNewEventRepositoryAdapter(t *testing.T) {
	NewEventRepositoryAdapter[*entity.BowlingGame](nil)
}
