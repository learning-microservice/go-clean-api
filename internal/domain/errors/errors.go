package errors

import (
	"encoding/json"
	goerrors "errors"
	"fmt"

	pkgErrors "go-clean-api/pkg/errors"
)

var (
	TypeValidation         = pkgErrors.Type[FieldError]("validation error")
	TypeInvalidCredentials = pkgErrors.Type[FieldError]("invalid credentials")
	TypeUnauthenticated    = pkgErrors.Type[FieldError]("unauthenticated")
	TypeForbidden          = pkgErrors.Type[FieldError]("forbidden")
	TypeNotFound           = pkgErrors.Type[FieldError]("not found")
	TypeAlreadyExists      = pkgErrors.Type[FieldError]("already exists")
	TypeInitialization     = pkgErrors.Type[FieldError]("initialization error")
	TypeUnexpected         = pkgErrors.Type[FieldError]("unexpected")
	TypeUnavailable        = pkgErrors.Type[FieldError]("unavailable")
	TypeTimeout            = pkgErrors.Type[FieldError]("timeout")
)

func AsError(err error) *pkgErrors.Error[FieldError] {
	return pkgErrors.AsError[FieldError](err)
}

func Is(err, target error) bool {
	return goerrors.Is(err, target)
}

type FieldError struct {
	field   string
	message string
}

func NewFieldError(field, message string) FieldError {
	return FieldError{
		field:   field,
		message: message,
	}
}

func (e FieldError) Field() string {
	return e.field
}

func (e FieldError) Message() string {
	return e.message
}

func (e FieldError) Error() string {
	return fmt.Sprintf("field %s: %s", e.field, e.message)
}

func (e FieldError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Field   string `json:"field,omitempty"`
		Message string `json:"message"`
	}{
		Field:   e.field,
		Message: e.message,
	})
}
