package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateConfig(cfg *Config) error {
	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("config validation error: %w", err)
	}

	for _, monitor := range cfg.Monitors {
		if err := validate.Struct(monitor); err != nil {
			return fmt.Errorf("monitor '%s' validation error: %w", monitor.Name, err)
		}
	}

	return nil
}
