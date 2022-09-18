package routing

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type Handler func(c echo.Context) error

type authorizer interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Client(ctx context.Context, t *oauth2.Token) *http.Client
}

type collector interface {
	Bootstrap(c *http.Client)
}

func WithRootHandler() Handler {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from bot!")
	}
}

func WithMercedesLoginHandler(auth authorizer) Handler {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL("login"))
	}
}

func WithMercedesLoginHandlerCallback(auth authorizer, coll collector) Handler {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		code := c.Request().URL.Query().Get("code")

		t, err := auth.Exchange(ctx, code)
		if err != nil {
			c.Logger().Errorf("Error exchanging code with access token: %v", err)
			return c.String(http.StatusBadRequest, "Error executing OAuth workflow")
		}

		coll.Bootstrap(auth.Client(ctx, t))

		return c.String(http.StatusOK, t.AccessToken)
	}
}
