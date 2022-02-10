package storage

import (
	"fmt"
)

var counterMemoryDatabase = make(map[string]Counter)

type CounterDataInMemory struct {
}

func (m *CounterDataInMemory) DatastoreName() string {
	return "memory"
}

func (m *CounterDataInMemory) Init(params string) (CounterData, error) {
	return m, nil
}

func (m *CounterDataInMemory) Terminate() error {
	return nil
}

func (m *CounterDataInMemory) Get(uuid string) (Counter, error) {
	if value, ok := counterMemoryDatabase[uuid]; ok {
		return value, nil
	}

	return Counter{}, fmt.Errorf("counter %s does not exist", uuid)
}

func (m *CounterDataInMemory) Create(counter Counter) (Counter, error) {
	if counter, ok := counterMemoryDatabase[counter.UUID]; ok {
		return counter, fmt.Errorf("counter %s already exists: %s", counter.UUID, counter.Name)
	}

	counterMemoryDatabase[counter.UUID] = counter
	return counter, nil
}

func (m *CounterDataInMemory) Update(counter Counter) (Counter, error) {
	if dataCounter, ok := counterMemoryDatabase[counter.UUID]; ok {
		dataCounter.Name = counter.Name
		dataCounter.Date = counter.Date
		dataCounter.Count = counter.Count
		dataCounter.Date = counter.Date
		counterMemoryDatabase[counter.UUID] = dataCounter
		return counter, nil
	}

	return Counter{}, fmt.Errorf("counter %s does not exist", counter.UUID)
}
