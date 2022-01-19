package counter

type CounterData interface {
	DatastoreName() string
	Init(params string) (CounterData, error)
	Exists(uuid string) (Counter, error)
	Create(Counter) (Counter, error)
	Increment(Counter) (Counter, error)
}
