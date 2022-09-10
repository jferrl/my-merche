package routing

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler func(c echo.Context) error

func WithRootHandler() Handler {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from bot!")
	}
}
