package config

type Config struct {
	Monitors []Monitor `mapstructure:"monitors"`
}

type Monitor struct {
	Name     string           `mapstructure:"name"`
	Type     string           `mapstructure:"type"`
	Timeout  string           `mapstructure:"timeout"`
	Interval string           `mapstructure:"interval"`
	Http     *HTTPMonitor     `mapstructure:"http,omitempty"`
	Database *DatabaseMonitor `mapstructure:"database,omitempty"`
}

type HTTPMonitor struct {
	Url     string            `mapstructure:"url"`
	Method  string            `mapstructure:"method"`
	Headers map[string]string `mapstructure:"headers,flow"`
}

type DatabaseMonitor struct {
	Driver           string  `mapstructure:"driver"`
	Query            string  `mapstructure:"query"`
	Host             *string `mapstructure:"host,omitempty"`
	Port             *string `mapstructure:"port,omitempty"`
	Username         *string `mapstructure:"username,omitempty"`
	Password         *string `mapstructure:"password,omitempty"`
	Database         *string `mapstructure:"database,omitempty"`
	ConnectionString *string `mapstructure:"connection_string,omitempty"`
}
