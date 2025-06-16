package config

import (
	"testing"
)

func strPtr(s string) *string { return &s }

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name:    "empty monitors",
			config:  &Config{Monitors: []Monitor{}},
			wantErr: true,
		},
		{
			name: "missing monitor name",
			config: &Config{Monitors: []Monitor{
				{Name: "", Type: "http", Timeout: "5s", Interval: "1m", Http: &HTTPMonitor{Url: "http://a", Method: "GET"}},
			}},
			wantErr: true,
		},
		{
			name: "valid http monitor",
			config: &Config{Monitors: []Monitor{
				{Name: "test", Type: "http", Timeout: "5s", Interval: "1m", Http: &HTTPMonitor{Url: "http://a", Method: "GET"}},
			}},
			wantErr: false,
		},
		{
			name: "http monitor missing url",
			config: &Config{Monitors: []Monitor{
				{Name: "test", Type: "http", Timeout: "5s", Interval: "1m", Http: &HTTPMonitor{Url: "", Method: "GET"}},
			}},
			wantErr: true,
		},
		{
			name: "valid postgres monitor with connection string",
			config: &Config{Monitors: []Monitor{
				{
					Name: "db1", Type: "database", Timeout: "5s", Interval: "1m",
					Database: &DatabaseMonitor{
						Driver:           "postgres",
						Query:            "SELECT 1",
						ConnectionString: strPtr("postgres://user:pass@host/db"),
					},
				},
			}},
			wantErr: false,
		},
		{
			name: "valid mysql monitor with host/port",
			config: &Config{Monitors: []Monitor{
				{
					Name: "db2", Type: "database", Timeout: "5s", Interval: "1m",
					Database: &DatabaseMonitor{
						Driver:   "mysql",
						Query:    "SELECT 1",
						Host:     strPtr("localhost"),
						Port:     strPtr("3306"),
						Username: strPtr("user"),
						Password: strPtr("pass"),
						Database: strPtr("db"),
					},
				},
			}},
			wantErr: false,
		},
		{
			name: "redis monitor missing database number",
			config: &Config{Monitors: []Monitor{
				{
					Name: "redis1", Type: "database", Timeout: "5s", Interval: "1m",
					Database: &DatabaseMonitor{
						Driver: "redis",
						Host:   strPtr("localhost"),
						Port:   strPtr("6379"),
					},
				},
			}},
			wantErr: true,
		},
		{
			name: "database monitor invalid driver",
			config: &Config{Monitors: []Monitor{
				{
					Name: "db3", Type: "database", Timeout: "5s", Interval: "1m",
					Database: &DatabaseMonitor{
						Driver: "oracle",
						Query:  "SELECT 1",
					},
				},
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
