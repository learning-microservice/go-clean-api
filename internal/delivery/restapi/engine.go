package restapi

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"go-clean-api/config"
	"go-clean-api/internal/registry"
)

func New(cfg *config.Server, reg *registry.Registry) http.Handler {
	// setup echo engine
	engine := echo.New()
	engine.Logger = reg.Logger
	engine.Validator = reg.Validator

	// setup router
	setupRouter(engine, cfg, reg)

	return engine
}
