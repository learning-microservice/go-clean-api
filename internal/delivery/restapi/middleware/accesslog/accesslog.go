package accesslog

import (
	goerrors "errors"
	"log/slog"
	"slices"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"go-clean-api/internal/delivery/restapi/httperror"
	"go-clean-api/internal/domain/errors"
)

var (
	msgSuccessRequest = "REQUEST"
	msgErrorRequest   = "REQUEST_ERROR"
	msgWarningRequest = "REQUEST_WARNING"
	msgPanicRequest   = "REQUEST_PANIC"
)

// New - .
func New(logger *slog.Logger) echo.MiddlewareFunc {
	isSystemError := func(status int) bool {
		return status >= 500
	}

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:       true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogRequestID:     true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogContentLength: true,
		LogResponseSize:  true,
		HandleError:      true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			ctx := c.Request().Context()

			logAttrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("host", v.Host),
				slog.String("bytes_in", v.ContentLength),
				slog.Int64("bytes_out", v.ResponseSize),
				slog.String("user_agent", v.UserAgent),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("request_id", v.RequestID),
			}

			// 正常終了: handler で error が返却されない
			if v.Error == nil {
				logger.LogAttrs(ctx, slog.LevelInfo, msgSuccessRequest, logAttrs...)
				return nil
			}

			// Panic Error: handler で panic が発生
			var panicErr *middleware.PanicStackError
			if ok := goerrors.As(v.Error, &panicErr); ok {
				stack := string(panicErr.Stack)
				stack = strings.ReplaceAll(stack, "\n\t", ": ")

				var message string
				if cause := panicErr.Unwrap(); cause != nil {
					message = cause.Error()
				} else {
					message = panicErr.Error()
				}
				logAttrs = slices.Concat(logAttrs, []slog.Attr{
					slog.String("error_type", errors.TypeUnexpected.Name()),
					slog.String("error_message", message),
					slog.Any("error_stack", strings.Split(stack, ",")),
				})

				logger.LogAttrs(ctx, slog.LevelError, msgPanicRequest, logAttrs...)
				return nil
			}

			// handle error attributes
			attrs := httperror.LogAttrs(v.Status, v.Error)
			logAttrs = slices.Concat(logAttrs, attrs)

			// handle http status
			message := msgWarningRequest
			level := slog.LevelWarn

			if isSystemError(v.Status) {
				message = msgErrorRequest
				level = slog.LevelError
			}

			logger.LogAttrs(ctx, level, message, logAttrs...)
			return nil
		},
	})
}
