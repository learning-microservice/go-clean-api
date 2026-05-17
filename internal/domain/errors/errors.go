// Package errors はドメイン全体で共有する業務エラー種別を定義する。
// 各 TypeXxx は pkg/errors.Type[FieldError](code, name) で HTTP 相当コードと JSON の type 用文字列を持つ。
//
// Wrap 時は外側の Type が優先される（httperror のステータス、TypeXxx.Is の判定、JSON の type/error）。
// 内側の code をクライアントに出したくないときに Wrap する（例: NotFound を InvalidCredentials で包む）。
package errors

import (
	"encoding/json"
	goerrors "errors"
	"fmt"

	pkgErrors "go-clean-api/pkg/errors"
)

var (
	TypeValidation         = pkgErrors.Type[FieldError](400, "validation error")
	TypeInvalidCredentials = pkgErrors.Type[FieldError](401, "invalid credentials")
	TypeUnauthenticated    = pkgErrors.Type[FieldError](401, "unauthenticated")
	TypeForbidden          = pkgErrors.Type[FieldError](403, "forbidden")
	TypeNotFound           = pkgErrors.Type[FieldError](404, "not found")
	TypeAlreadyExists      = pkgErrors.Type[FieldError](409, "already exists")
	TypeInitialization     = pkgErrors.Type[FieldError](500, "initialization error")
	TypeUnexpected         = pkgErrors.Type[FieldError](500, "unexpected")
	TypeUnavailable        = pkgErrors.Type[FieldError](503, "unavailable")
	TypeTimeout            = pkgErrors.Type[FieldError](504, "timeout")
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
