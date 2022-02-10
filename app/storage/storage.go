package storage

import "time"

type Counter struct {
	UUID  string    `json:"uuid"`
	Name  string    `json:"name"`
	Count int       `json:"count"`
	Date  time.Time `json:"created_at"`
}

type CounterData interface {
	DatastoreName() string
	Init(params string) (CounterData, error)
	Terminate() error
	Get(uuid string) (Counter, error)
	Create(Counter) (Counter, error)
	Update(Counter) (Counter, error)
}
