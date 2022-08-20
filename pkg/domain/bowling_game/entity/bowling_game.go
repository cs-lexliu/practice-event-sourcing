package entity

import (
	"fmt"

	"github.com/cs-lexliu/practice-event-sourcing/internal/ddd/centity"
)

var _ centity.AggregateRoot = &BowlingGame{}

type BowlingGame struct {
	*centity.AggregateRootCore
	id           string
	pins         int
	scores       int
	bonusChances []int
	ball         int
	frame        int
	status       BowlingGameStatus
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

func (b BowlingGame) ID() string {
	return b.id
}

func (b BowlingGame) Category() string {
	return "BowlingGame"
}

func (b BowlingGame) Pins() int {
	return b.pins
}

func (b BowlingGame) Scores() int {
	return b.scores
}

func (b BowlingGame) BonusChances() []int {
	return b.bonusChances
}

func (b BowlingGame) Ball() int {
	return b.ball
}

func (b BowlingGame) Frame() int {
	return b.frame
}

func (b BowlingGame) Status() BowlingGameStatus {
	return b.status
}

func (b *BowlingGame) When(event centity.DomainEvent) {
	switch event := interface{}(event).(type) {
	case BowlingGameCreatedEvent:
		b.id = event.id
		b.pins = 10
		b.frame = 1
	case BowlingGameRolledEvent:
		b.ball++
		for i, bonusChance := range b.bonusChances {
			if bonusChance > 0 {
				b.scores += event.hit
				b.bonusChances[i] = bonusChance - 1
			}
		}
		if !finalFrame(b.frame) {
			if strike(b.ball, b.pins, event.hit) {
				b.bonusChances = append(b.bonusChances, 2)
			}
			if spare(b.ball, b.pins, event.hit) {
				b.bonusChances = append(b.bonusChances, 1)
			}
		}
		b.scores += event.hit
		b.pins -= event.hit
		if !finalFrame(b.frame) {
			if b.pins == 0 || b.ball == 2 {
				b.pins = 10
				b.ball = 0
				b.frame++
			}
		} else {
			if gameFinished(b.ball, b.pins) {
				b.status = BowlingGameFinished
			}
			if b.pins == 0 {
				b.pins = 10
			}
		}
	}
}

func strike(ball, pins, hit int) bool {
	return ball == 1 && pins-hit == 0
}

func spare(ball, pins, hit int) bool {
	return ball == 2 && pins-hit == 0
}

func finalFrame(frame int) bool {
	return frame == 10
}

func gameFinished(ball, pins int) bool {
	return (ball == 2 && pins != 0) || ball == 3
}

func (b *BowlingGame) Insure() error {
	if b.pins > 10 {
		return fmt.Errorf("pins bigger than 10")
	}
	return nil
}

func Create(id string) (*BowlingGame, error) {
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
	if b.status == BowlingGameFinished {
		return fmt.Errorf("the game is finished")
	}
	if err := b.Apply(NewBowlingGameRolledEvent(b.id, hit)); err != nil {
		return fmt.Errorf("apply bowling game rolled event: %w", err)
	}
	return nil
}
