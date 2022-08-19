package bowling_game

import "fmt"

type BowlingGame struct {
	id     string
	pins   int // pins not bigger than 10
	scores int // scores not negative
}

func NewBowlingGame(id string) (*BowlingGame, error) {
	b := &BowlingGame{}
	if err := b.Apply(NewBowlingGameCreatedEvent(id)); err != nil {
		return nil, fmt.Errorf("apply: %w", err)
	}
	return b, nil
}

func NewBowlingGameFromEvent(events []BowlingGameEvent) (*BowlingGame, error) {
	b := &BowlingGame{}
	for _, event := range events {
		b.Apply(event)
	}
	return b, nil
}

func (b *BowlingGame) Roll(hit int) error {
	if hit < 0 {
		return fmt.Errorf("hit value is negative")
	}
	if err := b.Apply(NewBowlingGameRolledEvent(b.id, hit)); err != nil {
		return fmt.Errorf("apply: %w", err)
	}
	return nil
}

func (b *BowlingGame) Apply(event BowlingGameEvent) error {
	switch event := interface{}(event).(type) {
	case BowlingGameCreatedEvent:
		b.id = event.id
		b.pins = 10
	case BowlingGameRolledEvent:
		b.scores += event.hit
		b.pins -= event.hit
	}
	if err := b.Check(); err != nil {
		return fmt.Errorf("check: %w", err)
	}
	return nil
}

func (b *BowlingGame) Check() error {
	if b.pins > 10 {
		return fmt.Errorf("pins bigger than 10")
	}
	if b.scores < 0 {
		return fmt.Errorf("score is negative")
	}
	return nil
}

type BowlingGameEvent interface {
	bowlingGameEvent()
}

type BowlingGameCreatedEvent struct {
	BowlingGameEvent
	id string
}

func NewBowlingGameCreatedEvent(id string) BowlingGameCreatedEvent {
	return BowlingGameCreatedEvent{
		id: id,
	}
}

type BowlingGameRolledEvent struct {
	BowlingGameEvent
	id  string
	hit int
}

func NewBowlingGameRolledEvent(id string, hit int) BowlingGameRolledEvent {
	return BowlingGameRolledEvent{
		id:  id,
		hit: hit,
	}
}
