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
	"github.com/yanzay/tbot/v2"
)

var (
	port = os.Getenv("PORT")

	ttoken = os.Getenv("TELEGRAM_TOKEN")

	clientID     = os.Getenv("MERCEDES_CLIENT_ID")
	clientSecret = os.Getenv("MERCEDES_CLIENT_SECRET")

	admin = os.Getenv("BOT_ADMIN")
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
			MercedesAuthURL: "https://id.mercedes-benz.com/as/",
			ClientID:        clientID,
			ClientSecret:    clientSecret,
			Scopes: []string{
				"mb:vehicle:mbdata:fuelstatus",
				"mb:vehicle:mbdata:payasyoudrive",
				"mb:vehicle:mbdata:vehiclelock",
				"mb:vehicle:mbdata:vehiclestatus",
			},
			RedirectURI: "https://my-merche.herokuapp.com/login/mercedes/callback",
		},
	)

	e := echo.New()

	botSvr := tbot.New(ttoken,
		tbot.WithHTTPClient(cleanhttp.DefaultClient()),
		tbot.WithLogger(e.Logger),
	)
	bcli := botSvr.Client()

	botSvr.Use(bot.WithStat())
	botSvr.Use(bot.WithSecure(bot.Allowlist{
		admin: true,
	}, bcli))

	botSvr.HandleMessage(bot.WithLoginHandler(bcli))

	go botSvr.Start()

	e.GET("/", echo.HandlerFunc(routing.WithRootHandler()))
	e.GET("/login/mercedes", echo.HandlerFunc(routing.WithMercedesLoginHandler(authorizer)))
	e.GET("/login/mercedes/callback", echo.HandlerFunc(routing.WithMercedesLoginHandlerCallback(authorizer)))
	e.Logger.Fatal(e.Start(":" + port))
}
