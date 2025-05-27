package databaseClients

import (
	"context"
	"database/sql"
	"fmt"
	"wavezync/pulse-bridge/internal/types"

	_ "github.com/denisenkom/go-mssqldb"
)

func ExecMssqlQuery(useConnString bool, config DatabaseClientConfig) *types.MonitorError {
	var err error

	connectionStr := config.ConnString
	if !useConnString {
		connectionStr = fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s;encrypt=disable",
			config.Host, config.Port, config.Username, config.Password, config.Dbname)
	}

	mssqlDB, err := sql.Open("sqlserver", connectionStr)
	if err != nil {
		return types.NewConfigError(fmt.Errorf("failed to open database connection: %w", err))
	}
	defer mssqlDB.Close()

	mssqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	mssqlDB.SetMaxOpenConns(config.MaxOpenConns)
	mssqlDB.SetMaxIdleConns(config.MaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if err = mssqlDB.PingContext(ctx); err != nil {
		return types.NewClientError(fmt.Errorf("failed to ping database: %w", err))
	}

	if config.Query != "" {
		_, err = mssqlDB.QueryContext(ctx, config.Query)
		if err != nil {
			return types.NewClientError(fmt.Errorf("query execution failed: %w", err))
		}
	}

	return nil
}
