package storage

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/ericogr/go-counter-online/counter"
)

type CounterDataPostgresql struct {
	extraParams string
}

var (
	POSTGRES_DRIVER_NAME           = "postgres"
	POSTGRES_TABLE_EXISTENCE_QUERY = `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema  = $1
			AND   table_name   = $2
		);`
	POSTGRES_TABLE_CREATION_STATEMENT = `
		CREATE TABLE counter (
			uuid varchar(36) PRIMARY KEY,
			name varchar(64) NOT NULL,
			count integer NOT NULL,
			date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`
	POSTGRES_FIND_COUNTER_QUERY = `
		SELECT uuid, name, count, date FROM counter WHERE uuid = $1`
	POSTGRES_INSERT_COUNTER_STATEMENT = `
		INSERT INTO counter (uuid, name, count, date) VALUES ($1, $2, $3, $4)`
	POSTGRES_UDATE_COUNTER_STATEMENT = `
		UPDATE counter SET count = $1 WHERE uuid = $2`
)

func (m *CounterDataPostgresql) DatastoreName() string {
	return POSTGRES_DRIVER_NAME
}

func (m *CounterDataPostgresql) Init(params string) (counter.CounterData, error) {
	log.Println("Initializing postgres...")

	m.extraParams = params

	// TODO: use migrations technique
	tableExist, err := m.checkTableExists()
	if err != nil {
		return nil, err
	}

	if !tableExist {
		log.Println("Table doesn't exist, creating...")
		err := m.createTable()
		if err != nil {
			return nil, err
		}
		log.Println("Table created successfully")
	} else {
		log.Println("Table already exists")
	}

	return m, nil
}

func (m *CounterDataPostgresql) Exists(uuid string) (counter.Counter, error) {
	sqlDB, err := sql.Open(POSTGRES_DRIVER_NAME, m.extraParams)
	if err != nil {
		return counter.Counter{}, err
	}
	defer sqlDB.Close()

	var rowCounter counter.Counter
	err = sqlDB.QueryRow(POSTGRES_FIND_COUNTER_QUERY, uuid).
		Scan(&rowCounter.UUID, &rowCounter.Name, &rowCounter.Count, &rowCounter.Date)
	if err != nil {
		return rowCounter, err
	}

	return rowCounter, err
}

func (m *CounterDataPostgresql) Create(userCounter counter.Counter) (counter.Counter, error) {
	sqlDB, err := sql.Open(POSTGRES_DRIVER_NAME, m.extraParams)
	if err != nil {
		return counter.Counter{}, err
	}
	defer sqlDB.Close()

	currentTimestamp := time.Now()
	_, err = sqlDB.Exec(POSTGRES_INSERT_COUNTER_STATEMENT,
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
	sqlDB, err := sql.Open(POSTGRES_DRIVER_NAME, m.extraParams)
	if err != nil {
		return counter.Counter{}, err
	}
	defer sqlDB.Close()

	newCounter := rowCounter.Count + 1
	_, err = sqlDB.Exec(POSTGRES_UDATE_COUNTER_STATEMENT,
		newCounter, userCounter.UUID,
	)
	if err != nil {
		return counter.Counter{}, err
	}

	userCounter.Count = newCounter

	return userCounter, nil

}

func (m *CounterDataPostgresql) checkTableExists() (bool, error) {
	sqlDB, err := sql.Open(POSTGRES_DRIVER_NAME, m.extraParams)
	if err != nil {
		return false, err
	}
	defer sqlDB.Close()

	var exist bool
	err = sqlDB.QueryRow(POSTGRES_TABLE_EXISTENCE_QUERY, "public", "counter").
		Scan(&exist)

	return exist, err
}

func (m *CounterDataPostgresql) createTable() error {
	sqlDB, err := sql.Open(POSTGRES_DRIVER_NAME, m.extraParams)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec(POSTGRES_TABLE_CREATION_STATEMENT)

	return err
}
