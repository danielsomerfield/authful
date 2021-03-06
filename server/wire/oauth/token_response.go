package oauth

type TokenResponse struct {
	AccessToken  string        `json:"access_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    int64       `json:"expires_in"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	Scope        string        `json:"scope,omitempty"`
}
