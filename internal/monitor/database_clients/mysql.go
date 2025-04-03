package database_clients

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ExecMysqlQuery(useConnString bool, connString, host, port, username, password, database, query string, timeout time.Duration) error {
	var err error

	connectionStr := connString
	if !useConnString {
		connectionStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			username, password, host, port, database)
	}

	mysqlDB, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer mysqlDB.Close()

	// Optimize connection pool settings for health checks
	mysqlDB.SetConnMaxLifetime(timeout)
	mysqlDB.SetMaxOpenConns(1)
	mysqlDB.SetMaxIdleConns(0)

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err = mysqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = mysqlDB.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}

	return nil
}
