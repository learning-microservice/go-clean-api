package httperror

import (
	goerrors "errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"go-clean-api/internal/domain/errors"
	pkgErrors "go-clean-api/pkg/errors"
	"go-clean-api/pkg/validate"
)

const (
	msgErrorInvalidParameter = "invalid parameter"
)

type errorResponse struct {
	Type    string `json:"type"`
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

func Encode(c *echo.Context, err error) error {
	if err == nil {
		return nil
	}
	status, resp := handleErrorResponse(err)
	return c.JSON(status, resp)
}

func handleErrorResponse(err error) (status int, resp *errorResponse) {
	// handle domain error
	if de := errors.AsError(err); de != nil {
		return httpStatus(de), &errorResponse{
			Type:    de.Type(),
			Error:   de.Message(),
			Details: de.Details(),
		}
	}

	// handle Bind/Validate validate.Errors
	var verr *validate.Errors
	if ok := goerrors.As(err, &verr); ok {
		return http.StatusBadRequest, &errorResponse{
			Type:    errors.TypeValidation.Name(),
			Error:   msgErrorInvalidParameter,
			Details: verr.Unwrap(),
		}
	}

	// handle echo http error
	var herr *echo.HTTPError
	if ok := goerrors.As(err, &herr); ok {
		cause := herr.Unwrap()
		if cause == nil {
			cause = herr
		}
		return herr.StatusCode(), &errorResponse{
			Type:  errors.TypeValidation.Name(),
			Error: msgErrorInvalidParameter,
			Details: []errors.FieldError{
				errors.NewFieldError("", cause.Error()),
			},
		}
	}

	// handle unexpected error
	return http.StatusInternalServerError, &errorResponse{
		Type:  errors.TypeUnexpected.Name(),
		Error: "unexpected error occurred",
	}
}

func httpStatus(de *pkgErrors.Error[errors.FieldError]) int {
	switch {
	case errors.TypeValidation.Is(de):
		return http.StatusBadRequest
	case errors.TypeInvalidCredentials.Is(de),
		errors.TypeUnauthenticated.Is(de):
		return http.StatusUnauthorized
	case errors.TypeForbidden.Is(de):
		return http.StatusForbidden
	case errors.TypeNotFound.Is(de):
		return http.StatusNotFound
	case errors.TypeAlreadyExists.Is(de):
		return http.StatusConflict
	case errors.TypeTimeout.Is(de):
		return http.StatusGatewayTimeout // または 408
	case errors.TypeUnavailable.Is(de):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
