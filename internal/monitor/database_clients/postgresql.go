package database_clients

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func PgRun(useConnString bool, connString, host, port, username, password, database, query string) error {
	var err error

	connectionStr := connString
	if !useConnString {
		connectionStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			username, password, host, port, database)
	}

	pgDB, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer pgDB.Close()

	if err = pgDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = pgDB.Query(query)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}

	return nil
}
