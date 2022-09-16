package auth

type Session struct {
	token *OAuthAccessToken
}

func NewSession() *Session {
	return &Session{}
}
