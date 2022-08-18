package entity

var mapper = map[string]Constuctor{}

type Constuctor func([]DomainEvent) interface{}

func RegisterConstructor(id string, constructor Constuctor) {
	mapper[id] = constructor
}

func GetConstuctor(id string) Constuctor {
	return mapper[id]
}
