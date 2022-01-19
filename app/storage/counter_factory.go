package storage

import (
	"fmt"

	"github.com/ericogr/go-counter-online/counter"
)

var counterDatastores []counter.CounterData = []counter.CounterData{
	&CounterDataInMemory{},
	&CounterDataPostgresql{},
}

func GetCounterInstance(datastore string, params string) (counter.CounterData, error) {
	for _, counterDatastore := range counterDatastores {
		if counterDatastore.DatastoreName() == datastore {
			return counterDatastore.Init(params)
		}
	}

	return nil, fmt.Errorf("datastore not found: %s", datastore)
}
