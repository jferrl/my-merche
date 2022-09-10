package routing

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler func(c echo.Context) error

type authorizer interface {
	BuildMercedesLoginURL() string
}

func WithRootHandler() Handler {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from bot!")
	}
}

func WithMercedesLoginHandler(auth authorizer) Handler {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, auth.BuildMercedesLoginURL())
	}
}

func WithMercedesLoginHandlerCallback(auth authorizer) Handler {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Logged in")
	}
}
