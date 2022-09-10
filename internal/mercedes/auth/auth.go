package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Authorizer struct {
	clientID, clientSecret, scope string
}

type Opts struct {
	ClientID, ClientSecret, Scope string
}

func New(o Opts) *Authorizer {
	return &Authorizer{
		clientID:     o.ClientID,
		clientSecret: o.ClientSecret,
		scope:        o.Scope,
	}
}

func (a *Authorizer) BuildMercedesLoginURL() string {
	return fmt.Sprintf(
		"https://id.mercedes-benz.com/as/authorization.oauth2?response_type=code&client_id=%v&redirect_uri=%v&scope=%v&state=%v",
		a.clientID,
		"https://my-merche.herokuapp.com/login/mercedes/callback",
		a.scope,
		"login",
	)
}

func (a *Authorizer) GetMercedesAccessToken(code string) string {
	authToken := base64.StdEncoding.EncodeToString([]byte(a.clientID + ":" + a.clientSecret))

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
