package httperror

import (
	goerrors "errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"go-clean-api/internal/domain/errors"
	"go-clean-api/pkg/validate"
)

const (
	msgErrorInvalidParameter = "invalid parameter"
	msgErrorUnexpected       = "unexpected error occurred"
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
		return de.TypeCode(), &errorResponse{
			Type:    de.TypeName(),
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
		Error: msgErrorUnexpected,
	}
}
