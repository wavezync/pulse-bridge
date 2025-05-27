package config

import (
	"errors"
	"fmt"
)

func ValidateConfig(config *Config) error {
	if config == nil {
		return errors.New("config cannot be nil")
	}

	if len(config.Monitors) == 0 {
		return errors.New("no monitors defined in the configuration")
	}

	for i, monitor := range config.Monitors {
		if err := validateMonitor(monitor, i); err != nil {
			return err
		}
	}

	return nil
}

func validateMonitor(monitor Monitor, index int) error {
	if monitor.Name == "" {
		return fmt.Errorf("monitor at index %d is missing a name", index)
	}

	if monitor.Type == "" {
		return fmt.Errorf("monitor '%s' is missing a type", monitor.Name)
	}

	if monitor.Timeout == "" {
		return fmt.Errorf("monitor '%s' is missing timeout value", monitor.Name)
	}

	if monitor.Interval == "" {
		return fmt.Errorf("monitor '%s' is missing interval value", monitor.Name)
	}

	switch monitor.Type {
	case "http":
		return validateHTTPMonitor(monitor, index)
	case "database":
		return validateDatabaseMonitor(monitor, index)
	default:
		return fmt.Errorf("monitor '%s' has unknown type: %s (monitor %d)", monitor.Name, monitor.Type, index)
	}
}

func validateHTTPMonitor(monitor Monitor, index int) error {
	if monitor.Http == nil {
		return fmt.Errorf("monitor '%s' is of type 'http' but HTTP configuration is missing at monitor %d", monitor.Name, index)
	}

	if monitor.Http.Url == "" {
		return fmt.Errorf("HTTP monitor '%s' is missing URL at monitor %d", monitor.Name, index)
	}

	if monitor.Http.Method == "" {
		return fmt.Errorf("HTTP monitor '%s' is missing method at monitor %d", monitor.Name, index)
	}

	return nil
}

func validateDatabaseMonitor(monitor Monitor, index int) error {
	if monitor.Database == nil {
		return fmt.Errorf("monitor '%s' is of type 'database' but database configuration is missing at monitor %d", monitor.Name, index)
	}

	validDrivers := []string{"postgres", "mysql", "mariadb", "mssql", "redis"}
	driverValid := false
	for _, driver := range validDrivers {
		if monitor.Database.Driver == driver {
			driverValid = true
			break
		}
	}

	if !driverValid {
		return fmt.Errorf("database monitor '%s' has invalid driver '%s', must be one of %v at monitor %d",
			monitor.Name, monitor.Database.Driver, validDrivers, index)
	}

	if monitor.Database.Driver != "redis" && monitor.Database.Query == "" {
		return fmt.Errorf("database monitor '%s' is missing query at monitor %d", monitor.Name, index)
	}

	hasConnectionString := monitor.Database.ConnectionString != nil && *monitor.Database.ConnectionString != ""
	hasHostPort := monitor.Database.Host != nil && *monitor.Database.Host != "" &&
		monitor.Database.Port != nil && *monitor.Database.Port != ""

	if !hasConnectionString && !hasHostPort {
		return fmt.Errorf("database monitor '%s' requires either connection_string or host/port configuration at monitor %d",
			monitor.Name, index)
	}

	if !hasConnectionString {
		if monitor.Database.Driver == "redis" && monitor.Database.Database == nil {
			return fmt.Errorf("redis monitor '%s' requires database number when using host/port configuration at monitor %d",
				monitor.Name, index)
		}

		if monitor.Database.Driver != "redis" {
			if monitor.Database.Username == nil || *monitor.Database.Username == "" {
				return fmt.Errorf("database monitor '%s' is missing username for %s driver at monitor %d",
					monitor.Name, monitor.Database.Driver, index)
			}

			if monitor.Database.Password == nil {
				return fmt.Errorf("database monitor '%s' is missing password for %s driver at monitor %d",
					monitor.Name, monitor.Database.Driver, index)
			}

			if monitor.Database.Database == nil || *monitor.Database.Database == "" {
				return fmt.Errorf("database monitor '%s' is missing database name for %s driver at monitor %d",
					monitor.Name, monitor.Database.Driver, index)
			}
		}
	}

	return nil
}
