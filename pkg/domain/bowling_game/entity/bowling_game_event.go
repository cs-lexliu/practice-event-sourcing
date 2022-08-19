package entity

import "github.com/cs-lexliu/practice-event-sourcing/internal/ddd/centity"

type BowlingGameEvent interface {
	centity.DomainEvent
	bowlingGameEvent()
}

type BowlingGameCreatedEvent struct {
	centity.DomainEventCore
	BowlingGameEvent
	id string
}

func NewBowlingGameCreatedEvent(id string) BowlingGameCreatedEvent {
	return BowlingGameCreatedEvent{
		DomainEventCore: centity.NewDomainEventCore("BowlingGameCreatedEvent"),
		id:              id,
	}
}

type BowlingGameRolledEvent struct {
	centity.DomainEventCore
	BowlingGameEvent
	id  string
	hit int
}

func NewBowlingGameRolledEvent(id string, hit int) BowlingGameRolledEvent {
	return BowlingGameRolledEvent{
		DomainEventCore: centity.NewDomainEventCore("BowlingGameRolledEvent"),
		id:              id,
		hit:             hit,
	}
}
