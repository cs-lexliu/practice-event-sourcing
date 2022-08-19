package bowling_game

type BowlingGame struct {
	id string
}

func NewBowlingGame(id string) *BowlingGame {
	b := &BowlingGame{}
	b.Apply(NewBowlingGameCreatedEvent(id))
	return b
}

func (b *BowlingGame) Apply(event BowlingGameEvent) {
	switch event := interface{}(event).(type) {
	case BowlingGameCreatedEvent:
		b.id = event.id
	}
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
