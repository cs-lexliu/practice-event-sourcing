package entity

import core "github.com/cs-lexliu/practice-event-sourcing/src/core/entity"

type BowlingGameEvents struct {
	core.DomainEvent
	name string
}

func (e BowlingGameEvents) String() string {
	return e.name
}

type BowlingGameCreated struct {
	BowlingGameEvents
	gameID string
}

func NewBowlingGameCreated(gameID string) BowlingGameCreated {
	return BowlingGameCreated{
		BowlingGameEvents: BowlingGameEvents{
			name: "BowlingGameCreated",
		},
		gameID: gameID,
	}
}

type BowlingGameRolledABall struct {
	BowlingGameEvents
	hit int
}

func NewBowlingGameRollABall(hit int) BowlingGameRolledABall {
	return BowlingGameRolledABall{
		BowlingGameEvents: BowlingGameEvents{
			name: "BowlingGameRolledABall",
		},
		hit: hit,
	}
}
