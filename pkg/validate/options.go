package validate

import (
	"context"
	"reflect"

	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	goValidator "github.com/go-playground/validator/v10"
)

type Option func(*Validator) error

type RegisterTranslationFunc = func(v *goValidator.Validate, trans ut.Translator) (err error)

func WithTranslator(translator locales.Translator, translationFunc RegisterTranslationFunc, defaultFlag bool) Option {
	return func(v *Validator) error {
		return v.registerTranslation(translator, translationFunc, defaultFlag)
	}
}

func WithLocaleHandler(handler func(context.Context) string) Option {
	return func(v *Validator) error {
		v.localeHandler = handler
		return nil
	}
}

func WithFieldNameMap(locale string, fieldNames map[string]string) Option {
	return func(v *Validator) error {
		v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			if name, ok := fieldNames[field.Name]; ok {
				return name
			}
			return field.Name
		})
		return nil
	}
}
