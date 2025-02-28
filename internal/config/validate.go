package config

import (
	"wavezync/pulse-bridge/internal/config/translations"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func ValidateConfig(cfg *Config) error {
	validate = validator.New(validator.WithRequiredStructEnabled())

	var err error
	trans, err = translations.RegisterTranslations(validate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register translations")
		return err
	}

	if err := validate.Struct(cfg); err != nil {
		var errMsgs []string
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			errMsg := e.Translate(trans)
			errMsgs = append(errMsgs, errMsg)
			log.Error().
				Str("field", e.Field()).
				Str("tag", e.Tag()).
				Str("error", errMsg).
				Msg("Config validation error")
		}
		return err
	}

	for _, monitor := range cfg.Monitors {
		if err := validate.Struct(monitor); err != nil {
			var errMsgs []string
			validationErrors := err.(validator.ValidationErrors)
			for _, e := range validationErrors {
				errMsg := e.Translate(trans)
				errMsgs = append(errMsgs, errMsg)
				log.Error().
					Str("monitor", monitor.Name).
					Str("field", e.Field()).
					Str("tag", e.Tag()).
					Str("error", errMsg).
					Msg("Monitor validation error")
			}
			return err
		}
	}
	return nil
}
