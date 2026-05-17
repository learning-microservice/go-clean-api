package httperror

import (
	goerrors "errors"
	"log/slog"

	"github.com/labstack/echo/v5"

	"go-clean-api/internal/domain/errors"
	"go-clean-api/pkg/validate"
)

func LogAttrs(status int, err error) []slog.Attr {
	if err == nil {
		return nil
	}

	// handle domain error
	if de := errors.AsError(err); de != nil {
		if status >= 500 {
			return []slog.Attr{
				slog.String("error_type", de.TypeName()),
				slog.String("error_message", de.Message()),
				slog.Any("error_stack", de.StackFrames()),
			}
		}
		return []slog.Attr{
			slog.String("error_type", de.TypeName()),
			slog.String("error_message", de.Message()),
		}
	}

	// handle validation error
	var verr *validate.Errors
	if ok := goerrors.As(err, &verr); ok {
		return []slog.Attr{
			slog.String("error_type", errors.TypeValidation.Name()),
			slog.String("error_message", verr.Error()),
		}
	}

	// handle echo http error
	var herr *echo.HTTPError
	if ok := goerrors.As(err, &herr); ok {
		message := herr.Error()
		if cause := herr.Unwrap(); cause != nil {
			message = cause.Error()
		}
		return []slog.Attr{
			slog.String("error_type", herr.Error()),
			slog.String("error_message", message),
		}
	}

	return []slog.Attr{
		slog.String("error_type", errors.TypeUnexpected.Name()),
		slog.String("error_message", err.Error()),
	}
}
