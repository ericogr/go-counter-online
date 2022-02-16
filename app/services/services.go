package services

import (
	"time"

	"github.com/ericogr/go-counter-online/storage"
)

type Counter struct {
	UUID  string    `json:"uuid"`
	Name  string    `json:"name"`
	Count int       `json:"count"`
	Date  time.Time `json:"created_at"`
}

type CounterService interface {
	Exists(uuid string) (Counter, error)
	Create(counter Counter) (Counter, error)
	Increment(counter Counter) (Counter, error)
}

type DefaultCounterService struct {
	CounterData storage.CounterData
}

func (cs *DefaultCounterService) Exists(uuid string) (Counter, error) {
	data, err := cs.CounterData.Get(uuid)
	if err != nil {
		return Counter{}, err
	}

	return Counter{
		UUID:  data.UUID,
		Name:  data.Name,
		Count: data.Count,
		Date:  data.Date,
	}, nil
}

func (cs *DefaultCounterService) Create(counter Counter) (Counter, error) {
	storageCounter := storage.Counter{
		UUID:  counter.UUID,
		Name:  counter.Name,
		Date:  time.Now(),
		Count: 0,
	}
	data, err := cs.CounterData.Create(storageCounter)
	if err != nil {
		return Counter{}, err
	}

	return Counter{
		UUID:  data.UUID,
		Name:  data.Name,
		Date:  data.Date,
		Count: data.Count,
	}, nil
}

func (cs *DefaultCounterService) Increment(counter Counter) (Counter, error) {
	dataCounter, err := cs.CounterData.Get(counter.UUID)
	if err != nil {
		return Counter{}, err
	}

	dataCounter.Count++
	cs.CounterData.Update(dataCounter)

	return Counter{
		UUID:  dataCounter.UUID,
		Name:  dataCounter.Name,
		Date:  dataCounter.Date,
		Count: dataCounter.Count,
	}, nil
}
