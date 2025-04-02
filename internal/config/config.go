package config

type Config struct {
	Monitors []Monitor `yaml:"monitors"`
}

type Monitor struct {
	Name     string           `yaml:"name"`
	Type     string           `yaml:"type"`
	Timeout  string           `yaml:"timeout"`
	Interval string           `yaml:"interval"`
	Http     *HTTPMonitor     `yaml:"http,omitempty"`
	Database *DatabaseMonitor `yaml:"database,omitempty"`
}

type HTTPMonitor struct {
	Url     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers,flow"`
}

type DatabaseMonitor struct {
	Driver           string  `yaml:"driver"`
	Query            string  `yaml:"query"`
	Host             *string `yaml:"host,omitempty"`
	Port             *string `yaml:"port,omitempty"`
	Username         *string `yaml:"username,omitempty"`
	Password         *string `yaml:"password,omitempty"`
	Database         *string `yaml:"database,omitempty"`
	ConnectionString *string `yaml:"connection_string,omitempty"`
}
