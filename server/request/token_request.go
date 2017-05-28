package request

import "net/url"

type TokenRequest struct {
	GrantType string
	Scope     string
}

func ParseTokenRequest(values url.Values) (*TokenRequest, *ParseError) {
	tokenRequest := TokenRequest{}

	//TODO: add support for required fields for other grant types
	fields := map[string]*mapping{
		"grant_type": required(&tokenRequest.GrantType),
		"scope":      optional(&tokenRequest.Scope),
	}

	parseError := ParseRequest(values, fields)

	if parseError != nil {
		return nil, parseError
	} else {
		return &tokenRequest, nil
	}
}
