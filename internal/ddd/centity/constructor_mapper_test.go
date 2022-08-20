package centity

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConstructorMapperSuite struct {
	suite.Suite
}

func TestConstructorMapperSuite(t *testing.T) {
	suite.Run(t, new(ConstructorMapperSuite))
}

func (s *ConstructorMapperSuite) BeforeTest() {
	cm = newConstructorMap()
}

func (s *ConstructorMapperSuite) TestRegisterConstructor() {
	fakeAggregateRoot := &fakeAggregateRoot{}
	fakeConstructor := &fakeConstructor{}
	RegisterConstructor(fakeAggregateRoot, fakeConstructor.Constructor)
}

func (s *ConstructorMapperSuite) TestGetConstructor() {
	fakeAggregateRoot := &fakeAggregateRoot{}
	fakeConstructor := &fakeConstructor{}
	RegisterConstructor(fakeAggregateRoot, fakeConstructor.Constructor)

	got, err := GetConstructor(context.Background(), fakeAggregateRoot)
	s.NoError(err)
	got(nil)
	s.Equal(1, fakeConstructor.count)
}

func (s *ConstructorMapperSuite) TestGetNotExistedConstructorShouldReturnError() {
	cm = newConstructorMap()
	fakeAggregateRoot := &fakeAggregateRoot{}
	_, err := GetConstructor(context.Background(), fakeAggregateRoot)
	s.Error(err)
}

var _ AggregateRoot = &fakeAggregateRoot{}

type fakeAggregateRoot struct {
	AggregateRoot
}

func (a fakeAggregateRoot) Category() string {
	return "fakeAggregateRoot"
}

type fakeConstructor struct {
	count int
}

func (c *fakeConstructor) Constructor(events []DomainEvent) (AggregateRoot, error) {
	c.count++
	return &fakeAggregateRoot{}, nil
}
