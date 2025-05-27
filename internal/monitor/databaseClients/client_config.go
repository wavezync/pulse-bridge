package databaseClients

import (
	"fmt"
	"time"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/types"
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

func PrepareDBClientConfig(monitor *config.Monitor) (bool, DatabaseClientConfig, *types.MonitorError) {
	timeout, err := time.ParseDuration(monitor.Timeout)
	if err != nil {
		return false, DatabaseClientConfig{}, types.NewConfigError(fmt.Errorf("invalid timeout format: %w", err))
	}

	params := DatabaseClientConfig{
		Query:           monitor.Database.Query,
		Timeout:         timeout,
		MaxOpenConns:    1,
		MaxIdleConns:    0,
		ConnMaxLifetime: timeout,
	}

	hasConnString := monitor.Database.ConnectionString != nil && *monitor.Database.ConnectionString != "" && monitor.Database.Driver != ""

	var hasIndividualParams bool

	hasHostPort := monitor.Database.Host != nil && *monitor.Database.Host != "" &&
		monitor.Database.Port != nil && *monitor.Database.Port != ""

	if !hasHostPort {
		hasIndividualParams = false
	} else if monitor.Database.Driver == "redis" {
		hasIndividualParams = monitor.Database.Database != nil && *monitor.Database.Database != ""
	} else {
		hasIndividualParams = monitor.Database.Username != nil && *monitor.Database.Username != "" &&
			monitor.Database.Password != nil &&
			monitor.Database.Database != nil && *monitor.Database.Database != ""
	}

	if hasConnString {
		params.ConnString = *monitor.Database.ConnectionString
		return true, params, nil
	} else if hasIndividualParams {
		params.Host = *monitor.Database.Host
		params.Port = *monitor.Database.Port

		// Safe assignment for username - might be nil for Redis
		if monitor.Database.Username != nil {
			params.Username = *monitor.Database.Username
		}

		// Safe assignment for password - might be nil
		if monitor.Database.Password != nil {
			params.Password = *monitor.Database.Password
		}

		// Database name should be safe since we checked in hasIndividualParams
		params.Dbname = *monitor.Database.Database
		return false, params, nil
	}

	return false, DatabaseClientConfig{}, types.NewConfigError(fmt.Errorf("either connection_string or host and port must be specified"))
}
