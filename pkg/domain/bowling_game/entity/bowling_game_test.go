package entity

import (
	"testing"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/centity"
	"github.com/stretchr/testify/assert"
)

func TestNewBowlingGame(t *testing.T) {
	b, err := Create("1234")
	assert.NoError(t, err)
	assert.Equal(t, "1234", b.id)
}

func TestRollOneBallShouldIncreaseScoresAndDecreasePins(t *testing.T) {
	b, _ := Create("1234")

	err := b.Roll(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, b.scores)
	assert.Equal(t, 9, b.pins)
}

func TestRollOneBallWithBonusChangesShouldHaveBonusScores(t *testing.T) {
	b, _ := Create("")
	b.bonusChances = []int{1}

	err := b.Roll(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, b.scores)
	assert.Equal(t, []int{0}, b.bonusChances)
}

func TestRollTwoBallWithLeavingPinsShouldRestorePins(t *testing.T) {
	b, _ := Create("")

	b.Roll(1)
	b.Roll(1)
	assert.Equal(t, 10, b.pins)
}

func TestRollStrikePinsShouldRestorePins(t *testing.T) {
	b, _ := Create("")

	b.Roll(10)
	assert.Equal(t, 10, b.pins)
}

func TestRollStrikeShouldHaveTwoBonusChances(t *testing.T) {
	b, _ := Create("")

	b.Roll(10)
	assert.Equal(t, []int{2}, b.bonusChances)
}

func TestRollSparePinsShouldRestorePins(t *testing.T) {
	b, _ := Create("")

	b.Roll(1)
	b.Roll(9)
	assert.Equal(t, 10, b.pins)
}

func TestRollSpareShouldHaveOneBonusChances(t *testing.T) {
	b, _ := Create("")

	b.Roll(1)
	b.Roll(9)
	assert.Equal(t, []int{1}, b.bonusChances)
}

func TestFinalFrameShouldNotHaveBonusChance(t *testing.T) {
	b, _ := Create("")
	b.frame = 10

	b.Roll(10)
	var want []int
	assert.Equal(t, want, b.bonusChances)
}

func TestRollOneBallAfterGameFinishedShouldHaveError(t *testing.T) {
	b, _ := Create("")
	b.status = BowlingGameFinished

	err := b.Roll(1)
	assert.Error(t, err)
}

func TestRollThreeStrikeInFinalFrameShouldChangeToBowlingGameFinishedStatus(t *testing.T) {
	b, _ := Create("")
	b.frame = 10

	b.Roll(10)
	b.Roll(10)
	b.Roll(10)
	assert.Equal(t, BowlingGameFinished, b.status)
}

func TestRollTwoBallButNotCleanInFinalFrameShouldChangeToBowlingGameFinishedStatus(t *testing.T) {
	b, _ := Create("")
	b.frame = 10

	b.Roll(1)
	b.Roll(1)
	assert.Equal(t, BowlingGameFinished, b.status)
}

func TestPerfectGameShouldHaveThreeHundredScores(t *testing.T) {
	b, _ := Create("")
	for range make([]int, 12) {
		b.Roll(10)
	}
	assert.Equal(t, 300, b.scores)
}

func TestNewBowlingGameFromEvent(t *testing.T) {
	events := []centity.DomainEvent{
		NewBowlingGameCreatedEvent("123"),
		NewBowlingGameRolledEvent("123", 1),
	}

	b, err := newBowlingGameFromEvent(events)
	assert.NoError(t, err)
	assert.Nil(t, b.(*BowlingGame).DomainEvents())
}

func TestNewBowlingGameFromEventWithInvalidEventShouldReturnError(t *testing.T) {
	events := []centity.DomainEvent{
		NewBowlingGameCreatedEvent("123"),
		NewBowlingGameRolledEvent("123", -1),
	}

	_, err := newBowlingGameFromEvent(events)
	assert.Error(t, err)
}
