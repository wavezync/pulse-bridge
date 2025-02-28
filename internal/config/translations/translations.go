package translations

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rs/zerolog/log"
)

var (
	uni   *ut.UniversalTranslator
	trans ut.Translator
)

func RegisterTranslations(validate *validator.Validate) (ut.Translator, error) {
	enLocale := en.New()
	uni = ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Error().Err(err).Msg("Failed to register default translations")
		return nil, err
	}

	translations := []struct {
		tag         string
		translation string
		override    bool
	}{
		{
			tag:         "required",
			translation: "{0} is required",
			override:    true,
		},
		{
			tag:         "required_if",
			translation: "{0} is required when {1} is {2}",
			override:    true,
		},
		{
			tag:         "oneof",
			translation: "{0} must be one of: {1}",
			override:    true,
		},
		{
			tag:         "excluded_unless",
			translation: "{0} should not be set unless {1} is {2}",
			override:    true,
		},
		{
			tag:         "excluded_with",
			translation: "{0} cannot be used together with connection_string",
			override:    true,
		},
		{
			tag:         "required_without",
			translation: "{0} is required when connection_string is not provided",
			override:    true,
		},
		{
			tag:         "required_without_all",
			translation: "{0} is required when individual connection parameters are not provided",
			override:    true,
		},
	}

	for _, t := range translations {
		if err := validate.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc); err != nil {
			log.Error().
				Err(err).
				Str("tag", t.tag).
				Msg("Failed to register custom translation")
			return nil, err
		}
	}
	return trans, nil
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		return fe.Error()
	}
	return t
}
