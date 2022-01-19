package storage

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/ericogr/go-counter-online/counter"
)

type CounterDataPostgresql struct {
	extraParams string
}

func (m *CounterDataPostgresql) DatastoreName() string {
	return "postgresql"
}

func (m *CounterDataPostgresql) Init(params string) (counter.CounterData, error) {
	m.extraParams = params
	return m, nil
}

func (m *CounterDataPostgresql) Exists(uuid string) (counter.Counter, error) {
	sqlDB, err := sql.Open("postgres", m.extraParams)
	if err != nil {
		return counter.Counter{}, err
	}
	defer sqlDB.Close()

	var rowCounter counter.Counter
	err = sqlDB.QueryRow("SELECT uuid, name, count, date FROM counter WHERE uuid = $1", uuid).
		Scan(&rowCounter.UUID, &rowCounter.Name, &rowCounter.Count, &rowCounter.Date)
	if err != nil {
		return rowCounter, err
	}

	return rowCounter, err
}

func (m *CounterDataPostgresql) Create(userCounter counter.Counter) (counter.Counter, error) {
	sqlDB, err := sql.Open("postgres", m.extraParams)
	if err != nil {
		return counter.Counter{}, err
	}
	defer sqlDB.Close()

	currentTimestamp := time.Now()
	_, err = sqlDB.Exec("INSERT INTO counter (uuid, name, count, date) VALUES ($1, $2, $3, $4)",
		userCounter.UUID, userCounter.Name, userCounter.Count, currentTimestamp,
	)
	if err != nil {
		return counter.Counter{}, err
	}

	userCounter.Date = currentTimestamp

	return userCounter, nil
}

func (m *CounterDataPostgresql) Increment(userCounter counter.Counter) (counter.Counter, error) {
	rowCounter, err := m.Exists(userCounter.UUID)
	if err != nil {
		return rowCounter, err
	}
	sqlDB, err := sql.Open("postgres", m.extraParams)
	if err != nil {
		return counter.Counter{}, err
	}
	defer sqlDB.Close()

	newCounter := rowCounter.Count + 1
	_, err = sqlDB.Exec("UPDATE counter SET count = $1 WHERE uuid = $2",
		newCounter, userCounter.UUID,
	)
	if err != nil {
		return counter.Counter{}, err
	}

	userCounter.Count = newCounter

	return userCounter, nil

}
