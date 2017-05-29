package request

import (
	"net/http"
)

type TokenRequest struct {
	GrantType string
	Scope     string
}

func ParseTokenRequest(httpRequest http.Request) (*TokenRequest, error) {
	tokenRequest := TokenRequest{}

	//TODO: add support for required fields for other grant types
	fields := map[string]*mapping{
		"grant_type": required(&tokenRequest.GrantType),
		"scope":      optional(&tokenRequest.Scope),
	}

	if err := ParseRequest(httpRequest, fields); err != nil {
		return nil, err
	}
	return &tokenRequest, nil
}
