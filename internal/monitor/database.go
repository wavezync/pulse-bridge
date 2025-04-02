package monitor

import (
	"fmt"
	"slices"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/monitor/database_clients"
)

func DatabaseMonitor(monitor *config.Monitor) error {
	if monitor.Database == nil {
		return fmt.Errorf("database configuration is missing")
	}

	isSupported := slices.Contains([]string{"mysql", "postgres", "sqlite"}, monitor.Database.Driver)
	if !isSupported {
		return fmt.Errorf("unsupported database driver: %s", monitor.Database.Driver)
	}

	useConnString, connParams, err := prepareDbConnectionParams(monitor.Database)
	if err != nil {
		return err
	}

	switch monitor.Database.Driver {
	case "postgres":
		return database_clients.PgRun(
			useConnString,
			connParams.connString,
			connParams.host,
			connParams.port,
			connParams.username,
			connParams.password,
			connParams.dbname,
			monitor.Database.Query,
		)
	default:
		return fmt.Errorf("driver %s is supported but not implemented yet", monitor.Database.Driver)
	}
}

type DatabaseConnectionParams struct {
	connString string
	host       string
	port       string
	username   string
	password   string
	dbname     string
}

func prepareDbConnectionParams(dbConfig *config.DatabaseMonitor) (bool, DatabaseConnectionParams, error) {
	params := DatabaseConnectionParams{}

	hasConnString := dbConfig.ConnectionString != nil && *dbConfig.ConnectionString != ""

	hasIndividualParams := dbConfig.Host != nil && *dbConfig.Host != "" &&
		dbConfig.Port != nil && *dbConfig.Port != "" &&
		dbConfig.Username != nil && *dbConfig.Username != "" &&
		dbConfig.Password != nil && *dbConfig.Password != "" &&
		dbConfig.Database != nil && *dbConfig.Database != ""

	if hasConnString {
		params.connString = *dbConfig.ConnectionString
		return true, params, nil
	} else if hasIndividualParams {
		params.host = *dbConfig.Host
		params.port = *dbConfig.Port
		params.username = *dbConfig.Username
		params.password = *dbConfig.Password
		params.dbname = *dbConfig.Database
		return false, params, nil
	}

	return false, params, fmt.Errorf("either connection_string or host and port must be specified")
}
