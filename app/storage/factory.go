package storage

import (
	"fmt"
)

var counterDatastores []CounterData = []CounterData{
	&CounterDataInMemory{},
	&CounterDataPostgresql{},
}

func GetStoreInstance(datastore string, params string) (CounterData, error) {
	for _, counterDatastore := range counterDatastores {
		if counterDatastore.DatastoreName() == datastore {
			return counterDatastore.Init(params)
		}
	}

	return nil, fmt.Errorf("datastore not found: %s", datastore)
}
