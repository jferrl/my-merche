package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/heroku/x/hmetrics/onload"
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

	http.HandleFunc("/login/mercedes/", mercedesLoginHandler)
	http.HandleFunc("/login/mercedes/callback", mercedesCallbackHandler)

	bot := tbot.New(ttoken)
	c := bot.Client()
	bot.HandleMessage("/info", func(m *tbot.Message) {
		c.SendMessage(m.Chat.ID, "hello!")
	})
	bot.HandleMessage("/authorize", func(m *tbot.Message) {
		c.SendMessage(m.Chat.ID, "Open https://my-merche.herokuapp.com/login/mercedes")
	})
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Panic(
		http.ListenAndServe("127.0.0.1:"+port, nil),
	)
}

func mercedesLoginHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://id.mercedes-benz.com/as/authorization.oauth2?response_type=code&client_id=%v&redirect_uri=%v&scope=%v&state=%v",
		clientID,
		"https://my-merche.herokuapp.com/login/mercedes/callback",
		"mb:vehicle:mbdata:fuelstatus",
		"login",
	)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func mercedesCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	getMercedesAccessToken(code)
}

func getMercedesAccessToken(code string) string {
	authToken := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", "https://my-merche.herokuapp.com/login/mercedes/callback")
	data.Set("code", code)

	req, err := http.NewRequest(
		http.MethodPost,
		"https://id.mercedes-benz.com/as/token.oauth2",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Authorization", "Basic "+authToken)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	type mercedesAccessTokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	var res mercedesAccessTokenResponse
	json.Unmarshal(body, &res)

	return res.AccessToken
}
