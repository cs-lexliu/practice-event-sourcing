package bowling_game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBowlingGame(t *testing.T) {
	b := NewBowlingGame("1234")
	assert.Equal(t, "1234", b.id)
}
