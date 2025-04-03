package databaseClients

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

func ExecMssqlQuery(useConnString bool, config DatabaseClientConfig) error {
	var err error

	connectionStr := config.ConnString
	if !useConnString {
		connectionStr = fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s;encrypt=disable",
			config.Host, config.Port, config.Username, config.Password, config.Dbname)
	}

	mssqlDB, err := sql.Open("sqlserver", connectionStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer mssqlDB.Close()

	// Apply connection pool settings
	mssqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	mssqlDB.SetMaxOpenConns(config.MaxOpenConns)
	mssqlDB.SetMaxIdleConns(config.MaxIdleConns)

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if err = mssqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = mssqlDB.QueryContext(ctx, config.Query)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}

	return nil
}
