package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/jferrl/my-merche/internal/bot"
	"github.com/jferrl/my-merche/internal/http/routing"
	"github.com/jferrl/my-merche/internal/mercedes"
	"golang.org/x/oauth2"

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
	mercedes_vehicle_id     = os.Getenv("MERCEDES_VEHICLE_ID")

	bot_admin = os.Getenv("BOT_ADMIN")
)

func main() {
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	bootstrap()
}

func bootstrap() {

	ctx := context.Background()

	oauthConf := &oauth2.Config{
		ClientID:     mercedes_client_id,
		ClientSecret: merdeces_client_secret,
		RedirectURL:  mercedes_login_callback,
		Scopes: []string{
			"offline_access",
			"mb:vehicle:mbdata:fuelstatus",
			"mb:vehicle:mbdata:payasyoudrive",
			"mb:vehicle:mbdata:vehiclelock",
			"mb:vehicle:mbdata:vehiclestatus",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("%s/authorization.oauth2", mercedes_auth_url),
			TokenURL:  fmt.Sprintf("%s/token.oauth2", mercedes_auth_url),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}

	collector := mercedes.NewCollector(mercedes.VehicleID(mercedes_vehicle_id))

	e := echo.New()

	botSvr := tbot.New(telegram_token,
		tbot.WithHTTPClient(cleanhttp.DefaultClient()),
		tbot.WithLogger(e.Logger),
	)
	bcli := botSvr.Client()

	botSvr.Use(bot.WithSecure(bot.Allowlist{
		bot_admin: true,
	}, bcli))

	botSvr.HandleMessage(bot.WithStartHandler(bcli))
	botSvr.HandleMessage(bot.WithLoginHandler(bcli))
	botSvr.HandleMessage(bot.WithVehicleStatusHandler(ctx, bcli, collector))

	go botSvr.Start()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())

	e.GET("/", echo.HandlerFunc(routing.WithRootHandler()))
	e.GET("/login/mercedes", echo.HandlerFunc(routing.WithMercedesLoginHandler(oauthConf)))
	e.GET("/login/mercedes/callback", echo.HandlerFunc(routing.WithMercedesLoginHandlerCallback(oauthConf, collector)))
	e.Logger.Fatal(e.Start(":" + port))
}
