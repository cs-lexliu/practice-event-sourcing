package entity

import (
	"fmt"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/centity"
)

var _ centity.AggregateRoot = &BowlingGame{}

type BowlingGame struct {
	*centity.AggregateRootCore
	id     string
	pins   int
	scores int
}

func newBowlingGame() *BowlingGame {
	b := &BowlingGame{}
	b.AggregateRootCore = centity.NewAggregateRootTemple(b)
	return b
}

func init() {
	centity.RegisterConstructor(&BowlingGame{}, newBowlingGameFromEvent)
}

func newBowlingGameFromEvent(events []centity.DomainEvent) (interface{}, error) {
	b := newBowlingGame()
	for _, event := range events {
		if err := b.Apply(event); err != nil {
			return nil, fmt.Errorf("apply: %w", err)
		}
	}
	b.ClearDomainEvents()
	return b, nil
}

func (b BowlingGame) Id() string {
	return b.id
}

func (b BowlingGame) Pins() int {
	return b.pins
}

func (b BowlingGame) Scores() int {
	return b.scores
}

func (b BowlingGame) ID() string {
	return b.id
}

func (b BowlingGame) Category() string {
	return "BowlingGame"
}

func (b *BowlingGame) When(event centity.DomainEvent) {
	switch event := interface{}(event).(type) {
	case BowlingGameCreatedEvent:
		b.id = event.id
		b.pins = 10
	case BowlingGameRolledEvent:
		b.scores += event.hit
		b.pins -= event.hit
	}
}

func (b *BowlingGame) Insure() error {
	if b.pins > 10 {
		return fmt.Errorf("pins bigger than 10")
	}
	if b.scores < 0 {
		return fmt.Errorf("score is negative")
	}
	return nil
}

func NewBowlingGame(id string) (*BowlingGame, error) {
	b := newBowlingGame()
	if err := b.Apply(NewBowlingGameCreatedEvent(id)); err != nil {
		return nil, fmt.Errorf("apply bowling game created event: %w", err)
	}
	return b, nil
}

func (b *BowlingGame) Roll(hit int) error {
	if hit < 0 {
		return fmt.Errorf("hit value is negative")
	}
	if err := b.Apply(NewBowlingGameRolledEvent(b.id, hit)); err != nil {
		return fmt.Errorf("apply bowling game rolled event: %w", err)
	}
	return nil
}
