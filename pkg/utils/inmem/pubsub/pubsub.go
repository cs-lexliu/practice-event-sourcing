package pubsub

func NewPubSub[T any]() *Hub[T] {
	return &Hub[T]{
		subs: map[string]*Subscriber[T]{},
	}
}

type message any
