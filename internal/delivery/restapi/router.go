package restapi

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"go-clean-api/config"
	"go-clean-api/internal/delivery/restapi/v1/auth"
	"go-clean-api/internal/registry"
)

// SetupRouter -.
// Swagger spec:
//
//	@title       Go Clean Demo API
//	@description Multi-domain clean architecture template with translation, user, and task management
//	@version     1.0
//	@host        localhost:8080
//	@BasePath    /v1
//	@securityDefinitions.apikey BearerAuth
//	@in header
//	@name Authorization
func setupRouter(engine *echo.Echo, _ *config.Server, reg *registry.Registry) {
	// setup global middlewares (順序に注意！！)
	engine.Use(
		middleware.RequestLogger(),
		middleware.Recover(),
	)

	// setup v1 APIs
	v1Group := engine.Group("v1")
	v1Group.POST("/auth/login", auth.Login(reg.UsecaseSet.AuthLogin))
	v1Group.POST("/auth/register", auth.Register(reg.UsecaseSet.AuthRegister))
}
