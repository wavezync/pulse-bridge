package databaseClients

import (
	"fmt"
	"time"
	"wavezync/pulse-bridge/internal/config"
)

type DatabaseClientConfig struct {
	ConnString      string
	Host            string
	Port            string
	Username        string
	Password        string
	Dbname          string
	Query           string
	Timeout         time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func PrepareDBClientConfig(monitor *config.Monitor) (bool, DatabaseClientConfig, error) {

	timeout, err := time.ParseDuration(monitor.Timeout)
	if err != nil {
		return false, DatabaseClientConfig{}, fmt.Errorf("invalid timeout format: %w", err)
	}

	params := DatabaseClientConfig{
		Query:           monitor.Database.Query,
		Timeout:         timeout,
		MaxOpenConns:    1, // Default
		MaxIdleConns:    0, // Default
		ConnMaxLifetime: timeout,
	}

	hasConnString := monitor.Database.ConnectionString != nil && *monitor.Database.ConnectionString != ""

	hasIndividualParams := monitor.Database.Host != nil && *monitor.Database.Host != "" &&
		monitor.Database.Port != nil && *monitor.Database.Port != "" &&
		monitor.Database.Username != nil && *monitor.Database.Username != "" &&
		monitor.Database.Password != nil && *monitor.Database.Password != "" &&
		monitor.Database.Database != nil && *monitor.Database.Database != ""

	if hasConnString {
		params.ConnString = *monitor.Database.ConnectionString
		return true, params, nil
	} else if hasIndividualParams {
		params.Host = *monitor.Database.Host
		params.Port = *monitor.Database.Port
		params.Username = *monitor.Database.Username
		params.Password = *monitor.Database.Password
		params.Dbname = *monitor.Database.Database
		return false, params, nil
	}

	return false, DatabaseClientConfig{}, fmt.Errorf("either connection_string or host and port must be specified")
}
