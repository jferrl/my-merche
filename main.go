package main

import (
	"log"
	"os"

	"github.com/jferrl/my-merche/internal/bot"
	"github.com/jferrl/my-merche/internal/http/routing"
	"github.com/jferrl/my-merche/internal/mercedes/auth"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/labstack/echo/v4"
	"github.com/yanzay/tbot/v2"
)

var (
	port = os.Getenv("PORT")

	ttoken = os.Getenv("TELEGRAM_TOKEN")

	clientID     = os.Getenv("MERCEDES_CLIENT_ID")
	clientSecret = os.Getenv("MERCEDES_CLIENT_SECRET")
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
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scope:        "mb:vehicle:mbdata:fuelstatus",
			RedirectURI:  "https://my-merche.herokuapp.com/login/mercedes/callback",
		},
	)

	e := echo.New()

	b := tbot.New(ttoken)
	bc := b.Client()

	b.HandleMessage(bot.WithLoginHandler(bc))

	e.Logger.Fatal(b.Start())

	e.GET("/", echo.HandlerFunc(routing.WithRootHandler()))
	e.GET("/login/mercedes", echo.HandlerFunc(routing.WithMercedesLoginHandler(authorizer)))
	e.GET("/login/mercedes/callback", echo.HandlerFunc(routing.WithMercedesLoginHandlerCallback(authorizer)))
	e.Logger.Fatal(e.Start(":" + port))
}
