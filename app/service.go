package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-kit/log/level"
	_ "github.com/lib/pq"
)

type CounterService interface {
	Init() error
	Increment(context context.Context, uuid string) (Counter, error)
	Create(context context.Context, uuid string, name string) (Counter, error)
	Terminate() error
}

type Counter struct {
	UUID  string    `json:"uuid"`
	Name  string    `json:"name"`
	Count int       `json:"count"`
	Date  time.Time `json:"created_at"`
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
		UPDATE counter SET count = $2, name = $3, date = $4 WHERE uuid = $1`
)

type PostgresCounterService struct {
	DatabaseParams string
	database       *sql.DB
}

func (cs *PostgresCounterService) Init() error {
	level.Info(logger).Log("msg", "Initializing postgres...")

	var err error
	cs.database, err = sql.Open(POSTGRES_DRIVER_NAME, cs.DatabaseParams)
	if err != nil {
		return err
	}

	tableExist, err := cs.checkTableExists(context.Background())
	if err != nil {
		return err
	}

	if !tableExist {
		level.Info(logger).Log("msg", "Table doesn't exist, creating...")
		err := cs.createTable(context.Background())
		if err != nil {
			return err
		}
		level.Info(logger).Log("msg", "Table created successfully")
	} else {
		level.Info(logger).Log("msg", "Table already exists")
	}

	return nil
}

func (cs *PostgresCounterService) Create(ctx context.Context, uuid string, name string) (Counter, error) {
	ctx, cancel := createContextWithTimeout(ctx, 10)
	defer cancel()

	_, err := cs.database.ExecContext(ctx, POSTGRES_INSERT_COUNTER_STATEMENT,
		uuid, name, 0, time.Now(),
	)
	if err != nil {
		return Counter{}, err
	}

	return Counter{
		UUID:  uuid,
		Name:  name,
		Count: 0,
		Date:  time.Now(),
	}, nil
}

func (cs *PostgresCounterService) Increment(ctx context.Context, uuid string) (Counter, error) {
	counter, err := cs.get(ctx, uuid)
	if err != nil {
		return Counter{}, err
	}

	counter.Count++

	return cs.update(ctx, counter)
}

func (cs *PostgresCounterService) Terminate() error {
	return cs.database.Close()
}

func (cs *PostgresCounterService) get(ctx context.Context, uuid string) (Counter, error) {
	ctx, cancel := createContextWithTimeout(ctx, 5)
	defer cancel()

	var rowCounter Counter
	err := cs.database.QueryRowContext(ctx, POSTGRES_FIND_COUNTER_QUERY, uuid).
		Scan(&rowCounter.UUID, &rowCounter.Name, &rowCounter.Count, &rowCounter.Date)
	if err != nil {
		return rowCounter, err
	}

	return rowCounter, err
}

func (cs *PostgresCounterService) update(ctx context.Context, userCounter Counter) (Counter, error) {
	ctx, cancel := createContextWithTimeout(ctx, 5)
	defer cancel()

	result, err := cs.database.ExecContext(ctx, POSTGRES_UDATE_COUNTER_STATEMENT,
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
		return Counter{}, fmt.Errorf("no uuid found: %s", userCounter.UUID)
	}

	return userCounter, nil
}

func (cs *PostgresCounterService) createTable(ctx context.Context) error {
	ctx, cancel := createContextWithTimeout(ctx, 5)
	defer cancel()

	_, err := cs.database.ExecContext(ctx, POSTGRES_TABLE_CREATION_STATEMENT)

	return err
}

func (cs *PostgresCounterService) checkTableExists(ctx context.Context) (bool, error) {
	ctx, cancel := createContextWithTimeout(ctx, 5)
	defer cancel()

	var exist bool
	err := cs.database.QueryRowContext(ctx, POSTGRES_TABLE_EXISTENCE_QUERY, "public", "counter").
		Scan(&exist)

	return exist, err
}

type CounterServiceMiddleware func(CounterService) CounterService

func createContextWithTimeout(ctx context.Context, timeoutSeconds int64) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
}
