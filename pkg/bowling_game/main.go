package bowling_game

import "log"

func main() {
	s := NewStorage()

	b, _ := NewBowlingGame("b1")
	b.Roll(1)
	b.Roll(1)
	b.Roll(1)
	s.Append(b.id, b.events)
	log.Println(b)

	events, _ := s.Get(b.id)
	b2, _ := NewBowlingGameFromEvent(events)
	b2.Roll(1)
	s.Append(b2.id, b2.events)
	log.Println(b2)
}
