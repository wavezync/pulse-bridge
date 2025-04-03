package monitor

import (
	"fmt"
	"slices"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/monitor/databaseClients"
)

func DatabaseMonitor(monitor *config.Monitor) error {
	if monitor.Database == nil {
		return fmt.Errorf("database configuration is missing")
	}

	isSupported := slices.Contains([]string{"mysql", "mariadb", "postgres", "sqlite", "mssql"}, monitor.Database.Driver)
	if !isSupported {
		return fmt.Errorf("unsupported database driver: %s", monitor.Database.Driver)
	}

	useConnString, clientConfig, err := databaseClients.PrepareDBClientConfig(monitor)
	if err != nil {
		return err
	}

	switch monitor.Database.Driver {
	case "postgres":
		return databaseClients.ExecPgQuery(useConnString, clientConfig)
	case "mysql", "mariadb":
		return databaseClients.ExecMysqlQuery(useConnString, clientConfig)
	case "mssql":
		return databaseClients.ExecMssqlQuery(useConnString, clientConfig)
	default:
		return fmt.Errorf("driver %s is supported but not implemented yet", monitor.Database.Driver)
	}
}
