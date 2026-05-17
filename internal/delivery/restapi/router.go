package restapi

import (
	"github.com/labstack/echo/v5"

	"go-clean-api/config"
	"go-clean-api/internal/delivery/restapi/middleware/accesslog"
	"go-clean-api/internal/delivery/restapi/middleware/health"
	"go-clean-api/internal/delivery/restapi/middleware/recovery"
	"go-clean-api/internal/delivery/restapi/v1/auth"
	"go-clean-api/internal/registry"
)

// SetupRouter -.
func setupRouter(engine *echo.Echo, cfg *config.Server, reg *registry.Registry) {
	// setup global middlewares (順序に注意！！)
	engine.Use(
		health.New("/health",
			health.WithApp(cfg.APP.Name),
			health.WithVersion(cfg.APP.Version),
			health.WithEnv(cfg.APP.Env),
		),
		accesslog.New(reg.Logger),
		recovery.New(),
	)

	// setup v1 APIs
	v1Group := engine.Group("v1")
	v1Group.POST("/auth/login", auth.Login(reg.UsecaseSet.AuthLogin))
	v1Group.POST("/auth/register", auth.Register(reg.UsecaseSet.AuthRegister))
}
