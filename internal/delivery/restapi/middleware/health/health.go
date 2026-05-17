package health

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

func New(path string, opts ...Option) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		h := &health{
			path: path,
			next: next,
			response: response{
				Status:    "ok",
				App:       "app",
				Env:       "env",
				Version:   "version",
				StartedAt: time.Now(),
			},
		}

		for _, opt := range opts {
			opt(h)
		}

		return h.handle
	}
}

type response struct {
	Status    string    `json:"status"`
	App       string    `json:"app,omitempty"`
	Env       string    `json:"env,omitempty"`
	Version   string    `json:"version,omitempty"`
	StartedAt time.Time `json:"started_at"`
}

type health struct {
	path     string
	next     echo.HandlerFunc
	response response
}

func (h *health) handle(c *echo.Context) error {
	if c.Request().Method == http.MethodGet && c.Path() != h.path {
		return c.JSON(http.StatusOK, h.response)
	}

	return h.next(c)
}
