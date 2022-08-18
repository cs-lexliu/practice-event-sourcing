package entity

type AggregateRoot interface {
	ID() string
	Category() string
	When(event DomainEvent)
	DomainEvents() []DomainEvent
	ClearDomainEvents()
	Replay(events []DomainEvent) AggregateRoot
}

type AggregateRootCore struct {
	iAggregateRoot AggregateRoot
	domainEvents   []DomainEvent
}

func NewAggregateRootTemple[T AggregateRoot](t T) *AggregateRootCore {
	return &AggregateRootCore{
		iAggregateRoot: t,
	}
}

func (a *AggregateRootCore) Apply(event DomainEvent) {
	a.iAggregateRoot.When(event)
	a.addDomainEvent(event)
}

func (a *AggregateRootCore) addDomainEvent(event DomainEvent) {
	a.domainEvents = append(a.domainEvents, event)
}

func (a *AggregateRootCore) DomainEvents() []DomainEvent {
	return a.domainEvents
}

func (a *AggregateRootCore) ClearDomainEvents() {
	a.domainEvents = nil
}
