package main

import (
	"log"
	"os"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/jferrl/my-merche/internal/bot"
	"github.com/jferrl/my-merche/internal/http/routing"
	"github.com/jferrl/my-merche/internal/mercedes/auth"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yanzay/tbot/v2"
)

var (
	port = os.Getenv("PORT")

	telegram_token = os.Getenv("TELEGRAM_TOKEN")

	mercedes_client_id      = os.Getenv("MERCEDES_CLIENT_ID")
	merdeces_client_secret  = os.Getenv("MERCEDES_CLIENT_SECRET")
	mercedes_login_callback = os.Getenv("MERCEDES_LOGIN_CALLBACK")
	mercedes_auth_url       = os.Getenv("MERCEDES_AUTH_URL")

	bot_admin = os.Getenv("BOT_ADMIN")
)

func main() {
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	bootstrap()
}

func bootstrap() {
	authorizer := auth.New(
		auth.Opts{
			MercedesAuthURL: mercedes_auth_url,
			ClientID:        mercedes_client_id,
			ClientSecret:    merdeces_client_secret,
			Scopes: []string{
				"mb:vehicle:mbdata:fuelstatus",
				"mb:vehicle:mbdata:payasyoudrive",
				"mb:vehicle:mbdata:vehiclelock",
				"mb:vehicle:mbdata:vehiclestatus",
			},
			RedirectURI: mercedes_login_callback,
		},
	)

	e := echo.New()

	botSvr := tbot.New(telegram_token,
		tbot.WithHTTPClient(cleanhttp.DefaultClient()),
		tbot.WithLogger(e.Logger),
	)
	bcli := botSvr.Client()

	botSvr.Use(bot.WithSecure(bot.Allowlist{
		bot_admin: true,
	}, bcli))

	botSvr.HandleMessage(bot.WithLoginHandler(bcli))

	go botSvr.Start()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())

	e.GET("/", echo.HandlerFunc(routing.WithRootHandler()))
	e.GET("/login/mercedes", echo.HandlerFunc(routing.WithMercedesLoginHandler(authorizer)))
	e.GET("/login/mercedes/callback", echo.HandlerFunc(routing.WithMercedesLoginHandlerCallback(authorizer)))
	e.Logger.Fatal(e.Start(":" + port))
}
