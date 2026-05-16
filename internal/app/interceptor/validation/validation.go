package validation

import (
	"context"

	"go-clean-api/internal/app"
	"go-clean-api/internal/domain/errors"
	"go-clean-api/pkg/validate"
)

func Wrap[I, O any](validator *validate.Validator) app.Wrap[I, O] {
	return func(next app.Interactor[I, O]) app.Interactor[I, O] {
		i := &validation[I, O]{
			validator: validator,
			next:      next,
		}
		return i
	}
}

type validation[I, O any] struct {
	validator *validate.Validator
	next      app.Interactor[I, O]
}

func (v *validation[I, O]) Execute(ctx context.Context, input I) (O, error) {
	err := v.validator.ValidateCtx(ctx, input)
	if err != nil {
		var zero O
		return zero, errors.TypeValidation.Wrap(err, "validation error")
	}
	return v.next.Execute(ctx, input)
}
