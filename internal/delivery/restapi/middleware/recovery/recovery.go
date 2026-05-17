package recovery

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func New() echo.MiddlewareFunc {
	// echo.Recover のまま利用。panic 時は PanicStackError が返り、
	// accesslog（RequestLogger, HandleError: true）が REQUEST_PANIC を出力する。
	return middleware.Recover()
}
