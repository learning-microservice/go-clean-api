package auth

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"go-clean-api/api/openapi"
	authuc "go-clean-api/internal/app/usecase/auth"
)

func Login(service authuc.LoginUsecase) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var req openapi.LoginRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := c.Validate(&req); err != nil {
			return err
		}

		ctx := c.Request().Context()

		// Call usecase function
		output, err := service.Execute(ctx, &authuc.LoginInput{
			Email:         string(req.Email),
			PlainPassword: req.Password,
		})

		// handle error response
		if err != nil {
			return err
		}

		// return success response
		return c.JSON(http.StatusOK, openapi.LoginResponse{
			Token: output.Token,
		})
	}
}
