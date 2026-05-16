package auth

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"go-clean-api/api/openapi"
	authuc "go-clean-api/internal/app/usecase/auth"
	"go-clean-api/internal/domain/errors"
)

func Register(service authuc.RegisterUsecase) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var req openapi.RegisterRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
				Type:  errors.TypeValidation.Name(),
				Error: err.Error(),
			})
		}

		if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
				Type:  errors.TypeValidation.Name(),
				Error: err.Error(),
			})
		}

		ctx := c.Request().Context()

		// Call usecase function
		output, err := service.Execute(ctx, &authuc.RegisterInput{
			Name:          req.Name,
			Email:         string(req.Email),
			PlainPassword: req.Password,
		})

		// handle error response
		if err != nil {
			status := http.StatusInternalServerError
			errType := errors.TypeUnexpected
			if errors.TypeAlreadyExists.Is(err) {
				status = http.StatusConflict
				errType = errors.TypeAlreadyExists
			}
			return c.JSON(status, openapi.ErrorResponse{
				Type:  errType.Name(),
				Error: err.Error(),
			})
		}

		// return success response
		return c.JSON(http.StatusCreated, openapi.RegisterResponse{
			Id: output.ID.Int(),
		})
	}
}
