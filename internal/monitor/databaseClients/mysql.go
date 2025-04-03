package databaseClients

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ExecMysqlQuery(useConnString bool, config DatabaseClientConfig) error {
	var err error

	connectionStr := config.ConnString
	if !useConnString {
		connectionStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.Username, config.Password, config.Host, config.Port, config.Dbname)
	}

	mysqlDB, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer mysqlDB.Close()

	// Apply connection pool settings
	mysqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	mysqlDB.SetMaxOpenConns(config.MaxOpenConns)
	mysqlDB.SetMaxIdleConns(config.MaxIdleConns)

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if err = mysqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = mysqlDB.QueryContext(ctx, config.Query)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}

	return nil
}
