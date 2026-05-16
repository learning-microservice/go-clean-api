package validate

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	goValidator "github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/iancoleman/strcase"
)

func New(opts ...Option) (*Validator, error) {
	// default locale
	en := en.New()

	// setup validator
	validator := Validator{
		validate:      goValidator.New(),
		translator:    ut.New(en, en),
		defaultLocale: en.Locale(),
	}

	// register default translation (english)
	if err := validator.registerTranslation(en, translations.RegisterDefaultTranslations, true); err != nil {
		return nil, err
	}

	// setup options
	for _, opt := range opts {
		if err := opt(&validator); err != nil {
			return nil, err
		}
	}

	// setup locale handler
	if validator.localeHandler == nil {
		validator.localeHandler = func(_ context.Context) string {
			return validator.defaultLocale
		}
	}

	return &validator, nil
}

type Validator struct {
	validate      *goValidator.Validate
	translator    *ut.UniversalTranslator
	fieldNames    map[string]string
	defaultLocale string
	localeHandler func(context.Context) string
}

func (v *Validator) Validate(input any) error {
	return v.ValidateCtx(context.Background(), input)
}

func (v *Validator) ValidateCtx(ctx context.Context, input any) error {
	if cause := v.validate.StructCtx(ctx, input); cause != nil {
		validationErrors := cause.(goValidator.ValidationErrors)

		errs := make(FieldErrors, len(validationErrors))
		if errors.As(cause, &validationErrors) {
			locale := v.localeHandler(ctx)
			translator, found := v.translator.GetTranslator(locale)
			if !found {
				// TODO: warning log
				return cause
			}
			for i, e := range validationErrors {
				errs[i] = FieldError{
					Field:   strcase.ToSnake(e.StructField()), // TODO: snake or camel case (using strcase library)
					Message: e.Translate(translator),
				}
			}
		}
		return &errs
	}
	return nil
}

func (v *Validator) registerTranslation(translator locales.Translator, translationFunc RegisterTranslationFunc, defaultFlag bool) error {
	locale := translator.Locale()
	trans, found := v.translator.GetTranslator(locale)
	if !found {
		if err := v.translator.AddTranslator(translator, true); err != nil {
			return fmt.Errorf("failed to add %s translator: %w", locale, err)
		}
		trans, found = v.translator.GetTranslator(translator.Locale())
		if !found {
			return fmt.Errorf("failed to find %s translator", locale)
		}
	}
	if err := translationFunc(v.validate, trans); err != nil {
		return fmt.Errorf("failed to register %s translations: %w", locale, err)
	}
	if defaultFlag {
		v.defaultLocale = locale
	}
	return nil
}
