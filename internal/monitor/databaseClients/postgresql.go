package databaseClients

import (
	"context"
	"database/sql"
	"fmt"
	"wavezync/pulse-bridge/internal/types"

	_ "github.com/lib/pq"
)

func ExecPgQuery(useConnString bool, config DatabaseClientConfig) *types.MonitorError {
	var err error

	connectionStr := config.ConnString
	if !useConnString {
		connectionStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.Username, config.Password, config.Host, config.Port, config.Dbname)
	}

	pgDB, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return types.NewConfigError(fmt.Errorf("failed to open database connection: %w", err))
	}
	defer pgDB.Close()

	pgDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	pgDB.SetMaxOpenConns(config.MaxOpenConns)
	pgDB.SetMaxIdleConns(config.MaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if err = pgDB.PingContext(ctx); err != nil {
		return types.NewClientError(fmt.Errorf("failed to ping database: %w", err))
	}

	if config.Query != "" {
		_, err = pgDB.QueryContext(ctx, config.Query)
		if err != nil {
			return types.NewClientError(fmt.Errorf("query execution failed: %w", err))
		}
	}

	return nil
}
