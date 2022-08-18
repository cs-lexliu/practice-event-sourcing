package entity

type AggregateRoot interface {
	ID() string
	Category() string
	When(event DomainEvent)
	DomainEvents() []DomainEvent
	ClearDomainEvents()
	Replay(events []DomainEvent) AggregateRoot
}

type AggregateRootTemplate struct {
	iAggregateRoot AggregateRoot
	domainEvents   []DomainEvent
}

func NewAggregateRootTemple[T AggregateRoot](t T) *AggregateRootTemplate {
	return &AggregateRootTemplate{
		iAggregateRoot: t,
	}
}

func (a *AggregateRootTemplate) Apply(event DomainEvent) {
	a.iAggregateRoot.When(event)
	a.addDomainEvent(event)
}

func (a *AggregateRootTemplate) addDomainEvent(event DomainEvent) {
	a.domainEvents = append(a.domainEvents, event)
}

func (a *AggregateRootTemplate) DomainEvents() []DomainEvent {
	return a.domainEvents
}

func (a *AggregateRootTemplate) ClearDomainEvents() {
	a.domainEvents = nil
}
