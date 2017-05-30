package request

import (
	"net/http"
)

var GRANT_TYPE_CLIENT_CREDENTIALS = "client_credentials"

type TokenRequest struct {
	GrantType    string
	Scope        string
	ClientId     string
	ClientSecret string
}

func ParseTokenRequest(httpRequest http.Request) (*TokenRequest, error) {
	tokenRequest := TokenRequest{}

	//TODO: add support for required fields for other grant types
	fields := map[string]*mapping{
		"grant_type":    required(&tokenRequest.GrantType),
		"scope":         optional(&tokenRequest.Scope),
		"client_id":     optional(&tokenRequest.ClientId),
		"client_secret": optional(&tokenRequest.ClientSecret),
	}

	if err := ParseRequest(httpRequest, fields); err != nil {
		return nil, err
	}
	return &tokenRequest, nil
}
