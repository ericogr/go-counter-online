package storage

import (
	"fmt"

	"github.com/ericogr/go-counter-online/counter"
)

var counterMemoryDatabase = make(map[string]counter.Counter)

type CounterDataInMemory struct {
}

func (m *CounterDataInMemory) DatastoreName() string {
	return "memory"
}

func (m *CounterDataInMemory) Init(params string) (counter.CounterData, error) {
	return m, nil
}

func (m *CounterDataInMemory) Exists(uuid string) (counter.Counter, error) {
	return counterMemoryDatabase[uuid], nil
}

func (m *CounterDataInMemory) Create(counter counter.Counter) (counter.Counter, error) {
	if counter, ok := counterMemoryDatabase[counter.UUID]; ok {
		return counter, fmt.Errorf("counter %s already exists: %s", counter.UUID, counter.Name)
	}

	counterMemoryDatabase[counter.UUID] = counter
	return counter, nil
}

func (m *CounterDataInMemory) Increment(userCounter counter.Counter) (counter.Counter, error) {
	if counter, ok := counterMemoryDatabase[userCounter.UUID]; ok {
		counter.Count++
		counterMemoryDatabase[userCounter.UUID] = counter
		return counter, nil
	}

	return counter.Counter{}, fmt.Errorf("counter %s does not exist", userCounter.UUID)
}
