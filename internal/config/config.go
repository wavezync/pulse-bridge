package config

type Config struct {
	Monitors []Monitor `yaml:"monitors"`
}

type Monitor struct {
	Name     string           `yaml:"name" validate:"required"`
	Type     string           `yaml:"type" validate:"required,oneof=http database"`
	Timeout  string           `yaml:"timeout" validate:"required"`
	Interval string           `yaml:"interval" validate:"required"`
	Http     *HTTPMonitor     `yaml:"http,omitempty" validate:"required_if=Type http,excluded_unless=Type http"`
	Database *DatabaseMonitor `yaml:"database,omitempty" validate:"required_if=Type database,excluded_unless=Type database"`
}

type HTTPMonitor struct {
	Url     string            `yaml:"url" validate:"required"`
	Method  string            `yaml:"method" validate:"required"`
	Headers map[string]string `yaml:"headers,flow"`
}

type DatabaseMonitor struct {
	Driver           string  `yaml:"driver" validate:"required"`
	Query            string  `yaml:"query" validate:"required"`
	Host             *string `yaml:"host,omitempty" validate:"required_without=ConnectionString,excluded_with=ConnectionString"`
	Port             *string `yaml:"port,omitempty" validate:"required_without=ConnectionString,excluded_with=ConnectionString"`
	Username         *string `yaml:"username,omitempty" validate:"required_without=ConnectionString,excluded_with=ConnectionString"`
	Password         *string `yaml:"password,omitempty" validate:"required_without=ConnectionString,excluded_with=ConnectionString"`
	Database         *string `yaml:"database,omitempty" validate:"required_without=ConnectionString,excluded_with=ConnectionString"`
	ConnectionString *string `yaml:"connection_string,omitempty" validate:"required_without_all=Host Port Username Password Database"`
}
