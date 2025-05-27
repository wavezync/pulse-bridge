package databaseClients

import (
	"context"
	"database/sql"
	"fmt"
	"wavezync/pulse-bridge/internal/types"

	_ "github.com/go-sql-driver/mysql"
)

func ExecMysqlQuery(useConnString bool, config DatabaseClientConfig) *types.MonitorError {
	var err error

	connectionStr := config.ConnString
	if !useConnString {
		connectionStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.Username, config.Password, config.Host, config.Port, config.Dbname)
	}

	mysqlDB, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return types.NewConfigError(fmt.Errorf("failed to open database connection: %w", err))
	}
	defer mysqlDB.Close()

	mysqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	mysqlDB.SetMaxOpenConns(config.MaxOpenConns)
	mysqlDB.SetMaxIdleConns(config.MaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if err = mysqlDB.PingContext(ctx); err != nil {
		return types.NewClientError(fmt.Errorf("failed to ping database: %w", err))
	}

	if config.Query != "" {
		_, err = mysqlDB.QueryContext(ctx, config.Query)
		if err != nil {
			return types.NewClientError(fmt.Errorf("query execution failed: %w", err))
		}
	}

	return nil
}
