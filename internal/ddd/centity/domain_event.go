package centity

type DomainEvent interface {
	domainEvent()
}

type DomainEventCore struct {
	DomainEvent
	name string
}

func NewDomainEventCore(name string) DomainEventCore {
	return DomainEventCore{}
}
