package validation

import (
	"context"

	"go-clean-api/internal/app"
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
		// TODO: pkg/validate Errors はそのまま返却
		var zero O
		return zero, err
	}
	return v.next.Execute(ctx, input)
}

/*
func convertFieldErrors(errs []error) []errors.FieldError {
	details := make([]errors.FieldError, 0, len(errs))
	for _, e := range errs {
		var fe errors.FieldError
		if !goerrors.As(e, &fe) {
			fe = errors.NewFieldError("", e.Error())
		}
		details = append(details, fe)
	}
	return details
}
*/
