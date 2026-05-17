package restapi

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"

	"go-clean-api/config"
	"go-clean-api/internal/delivery/restapi/httperror"
	"go-clean-api/internal/registry"
)

func New(cfg *config.Server, reg *registry.Registry) http.Handler {
	// setup echo engine
	engine := echo.New()
	engine.Logger = reg.Logger
	engine.Validator = reg.Validator
	engine.HTTPErrorHandler = func(c *echo.Context, err error) {
		if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
			if resp.Committed {
				return // response has been already sent to the client by handler or some middleware
			}
		}
		if err := httperror.Encode(c, err); err != nil {
			reg.Logger.Error("failed to encode error response", "error", err)
		}
	}

	// setup router
	setupRouter(engine, cfg, reg)

	if strings.EqualFold(cfg.Log.Level, "debug") {
		routes := engine.Router().Routes()
		reg.Logger.Info("echo routing table", "routes", routes)
	}

	return engine
}
