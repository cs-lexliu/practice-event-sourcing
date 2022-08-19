package bowling_game

type Storage struct {
	data map[string][]BowlingGameEvent
}

func NewStorage() *Storage {
	return &Storage{
		data: map[string][]BowlingGameEvent{},
	}
}

func (s *Storage) Append(id string, events []BowlingGameEvent) error {
	s.data[id] = append(s.data[id], events...)
	return nil
}

func (s *Storage) Get(id string) ([]BowlingGameEvent, error) {
	events := s.data[id]
	return events, nil
}
