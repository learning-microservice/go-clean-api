package errors

import (
	goerrors "errors"

	pkgErrors "go-clean-api/pkg/errors"
)

var (
	TypeValidation         = pkgErrors.Type("validation error")
	TypeInvalidCredentials = pkgErrors.Type("invalid credentials")
	TypeUnauthenticated    = pkgErrors.Type("unauthenticated")
	TypeForbidden          = pkgErrors.Type("forbidden")
	TypeNotFound           = pkgErrors.Type("not found")
	TypeAlreadyExists      = pkgErrors.Type("already exists")
	TypeInitialization     = pkgErrors.Type("initialization error")
	TypeUnexpected         = pkgErrors.Type("unexpected")
	TypeUnavailable        = pkgErrors.Type("unavailable")
	TypeTimeout            = pkgErrors.Type("timeout")
)

func AsError(err error) *pkgErrors.Error {
	return pkgErrors.AsError(err)
}

func Is(err, target error) bool {
	return goerrors.Is(err, target)
}
