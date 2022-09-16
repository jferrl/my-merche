package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
)

type Authorizer struct {
	mercedesAuthURL string
	clientID        string
	clientSecret    string
	scope           string
	basicAuthToken  string
	redirectURI     string
	client          *http.Client
}

type Opts struct {
	MercedesAuthURL string
	ClientID        string
	ClientSecret    string
	Scopes          []string
	RedirectURI     string
}

func New(o Opts) *Authorizer {
	bat := base64.StdEncoding.EncodeToString([]byte(o.ClientID + ":" + o.ClientSecret))

	return &Authorizer{
		basicAuthToken: bat,
		scope:          strings.Join(o.Scopes, " "),
		redirectURI:    o.RedirectURI,
		clientID:       o.ClientID,
		clientSecret:   o.ClientSecret,
		client:         cleanhttp.DefaultClient(),
	}
}

func (a *Authorizer) BuildMercedesLoginURL() string {
	return fmt.Sprintf(
		"%s/authorization.oauth2?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		a.mercedesAuthURL,
		a.clientID,
		a.redirectURI,
		a.scope,
		"login",
	)
}

func (a *Authorizer) ExchangeAuthCodeWithAccessToken(ctx context.Context, code string) (*OAuthAccessToken, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", a.redirectURI)
	data.Set("code", code)

	resp, err := a.oauthDo(ctx, data)
	if err != nil {
		return nil, err
	}

	var res OAuthAccessToken
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (a *Authorizer) RefreshTokens(ctx context.Context, rf string) (*OAuthAccessToken, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", rf)

	resp, err := a.oauthDo(ctx, data)
	if err != nil {
		return nil, err
	}

	var res OAuthAccessToken
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (a *Authorizer) oauthDo(ctx context.Context, data url.Values) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%v/token.oauth2", a.mercedesAuthURL),
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+a.basicAuthToken)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
