package database_clients

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func ExecMssqlQuery(useConnString bool, connString, host, port, username, password, database, query string, timeout time.Duration) error {
	var err error

	connectionStr := connString
	if !useConnString {
		connectionStr = fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s;encrypt=disable",
			host, port, username, password, database)
	}

	mssqlDB, err := sql.Open("sqlserver", connectionStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer mssqlDB.Close()

	// Optimize connection pool settings for health checks
	mssqlDB.SetConnMaxLifetime(timeout)
	mssqlDB.SetMaxOpenConns(1)
	mssqlDB.SetMaxIdleConns(0)

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err = mssqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = mssqlDB.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}

	return nil
}
