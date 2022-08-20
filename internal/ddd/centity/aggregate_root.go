package centity

import "fmt"

type AggregateRoot interface {
	ID() string
	Category() string
	When(event DomainEvent)
	Insure() error
	DomainEvents() []DomainEvent
	ClearDomainEvents()
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

func (a *AggregateRootCore) Apply(event DomainEvent) error {
	if err := a.iAggregateRoot.Insure(); err != nil {
		return fmt.Errorf("pre-insure: %w", err)
	}
	a.iAggregateRoot.When(event)
	if err := a.iAggregateRoot.Insure(); err != nil {
		return fmt.Errorf("post-insure: %w", err)
	}
	a.domainEvents = append(a.domainEvents, event)
	return nil
}

func (a *AggregateRootCore) DomainEvents() []DomainEvent {
	return a.domainEvents
}

func (a *AggregateRootCore) ClearDomainEvents() {
	a.domainEvents = nil
}
