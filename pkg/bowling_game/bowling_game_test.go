package bowling_game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBowlingGame(t *testing.T) {
	b, err := NewBowlingGame("1234")
	assert.NoError(t, err)
	assert.Equal(t, "1234", b.id)
}

func TestRollOneBallShouldIncreaseScoresAndDecreasePins(t *testing.T) {
	b, _ := NewBowlingGame("1234")

	err := b.Roll(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, b.scores)
	assert.Equal(t, 9, b.pins)
}
