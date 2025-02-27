package config

type Config struct {
	Monitors []Monitor `yaml:"monitors" json:"monitors"`
}

type Monitor struct {
	Name     string           `yaml:"name" json:"name"`
	Type     string           `yaml:"type" json:"type"`
	Timeout  string           `yaml:"timeout" json:"timeout"`
	Interval string           `yaml:"interval" json:"interval"`
	Http     *HTTPMonitor     `yaml:"http,omitempty" json:"http,omitempty"`
	Database *DatabaseMonitor `yaml:"database,omitempty" json:"database,omitempty"`
}

type HTTPMonitor struct {
	Url     string            `yaml:"url" json:"url"`
	Method  string            `yaml:"method" json:"method"`
	Headers map[string]string `yaml:"headers,flow" json:"headers"`
}

type DatabaseMonitor struct {
	Driver           string  `yaml:"driver" json:"driver"`
	Host             *string `yaml:"host,omitempty" json:"host,omitempty"`
	Port             *string `yaml:"port,omitempty" json:"port,omitempty"`
	Username         *string `yaml:"username,omitempty" json:"username,omitempty"`
	Password         *string `yaml:"password,omitempty" json:"password,omitempty"`
	Database         *string `yaml:"database,omitempty" json:"database,omitempty"`
	Query            string  `yaml:"query" json:"query"`
	ConnectionString *string `yaml:"connection_string,omitempty" json:"connection_string,omitempty"`
}
