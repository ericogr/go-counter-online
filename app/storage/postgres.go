package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type CounterDataPostgresql struct {
	extraParams string
	database    *sql.DB
}

var (
	POSTGRES_DRIVER_NAME           = "postgres"
	POSTGRES_MIN_RECONN            = 10 * time.Second
	POSTGRES_MAX_RECONN            = time.Minute
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
		UPDATE counter SET count = $2, name = $3, date = $4 WHERE uuid = $1`
)

func (m *CounterDataPostgresql) DatastoreName() string {
	return POSTGRES_DRIVER_NAME
}

func (m *CounterDataPostgresql) Init(params string) (CounterData, error) {
	log.Println("Initializing postgres...")

	m.extraParams = params
	var err error
	m.database, err = sql.Open(POSTGRES_DRIVER_NAME, m.extraParams)
	if err != nil {
		return m, err
	}

	// TODO: use migrations...
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

func (m *CounterDataPostgresql) Terminate() error {
	return m.database.Close()
}

func (m *CounterDataPostgresql) Get(uuid string) (Counter, error) {
	var rowCounter Counter
	err := m.database.QueryRow(POSTGRES_FIND_COUNTER_QUERY, uuid).
		Scan(&rowCounter.UUID, &rowCounter.Name, &rowCounter.Count, &rowCounter.Date)
	if err != nil {
		return rowCounter, err
	}

	return rowCounter, err
}

func (m *CounterDataPostgresql) Create(counter Counter) (Counter, error) {
	_, err := m.database.Exec(POSTGRES_INSERT_COUNTER_STATEMENT,
		counter.UUID, counter.Name, counter.Count, counter.Date,
	)
	if err != nil {
		return Counter{}, err
	}

	return counter, nil
}

func (m *CounterDataPostgresql) Update(userCounter Counter) (Counter, error) {
	result, err := m.database.Exec(POSTGRES_UDATE_COUNTER_STATEMENT,
		userCounter.UUID,
		userCounter.Count,
		userCounter.Name,
		userCounter.Date,
	)
	if err != nil {
		return Counter{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return Counter{}, err
	}

	if rows == 0 {
		return Counter{}, fmt.Errorf("no counter found: %s", userCounter.UUID)
	}

	return userCounter, nil
}

func (m *CounterDataPostgresql) checkTableExists() (bool, error) {
	var exist bool
	err := m.database.QueryRow(POSTGRES_TABLE_EXISTENCE_QUERY, "public", "counter").
		Scan(&exist)

	return exist, err
}

func (m *CounterDataPostgresql) createTable() error {
	_, err := m.database.Exec(POSTGRES_TABLE_CREATION_STATEMENT)

	return err
}
