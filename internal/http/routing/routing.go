package routing

import (
	"context"
	"net/http"

	"github.com/jferrl/my-merche/internal/mercedes/auth"
	"github.com/labstack/echo/v4"
)

type Handler func(c echo.Context) error

type authorizer interface {
	BuildMercedesLoginURL() string
	ExchangeAuthCodeWithAccessToken(ctx context.Context, code string) (*auth.OAuthAccessToken, error)
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
		code := c.Request().URL.Query().Get("code")

		c.Logger().Infof("Auth code: %s", code)

		_, err := auth.ExchangeAuthCodeWithAccessToken(c.Request().Context(), code)
		if err != nil {
			c.Logger().Errorf("Error exchanging code with access token: %v", err)
			return c.String(http.StatusBadRequest, "Error executing OAuth workflow")
		}

		return c.String(http.StatusOK, "Authorized")
	}
}
