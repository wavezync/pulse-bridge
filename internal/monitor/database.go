package monitor

import (
	"fmt"
	"slices"
	"wavezync/pulse-bridge/internal/config"
	"wavezync/pulse-bridge/internal/monitor/databaseClients"
	"wavezync/pulse-bridge/internal/types"
)

func DatabaseMonitor(monitor *config.Monitor) *types.MonitorError {
	if monitor.Database == nil {
		return types.NewConfigError(fmt.Errorf("database configuration is missing"))
	}

	isSupported := slices.Contains([]string{"mysql", "mariadb", "postgres", "mssql", "redis"}, monitor.Database.Driver)
	if !isSupported {
		return types.NewConfigError(fmt.Errorf("unsupported database driver: %s", monitor.Database.Driver))
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
	case "redis":
		return databaseClients.ExecRedisQuery(useConnString, clientConfig)
	default:
		return types.NewConfigError(fmt.Errorf("driver %s is not supported yet", monitor.Database.Driver))
	}
}
