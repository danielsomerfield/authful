package request

import (
	"net/url"
)

type AuthorizeRequest struct {
	ResponseType string
	ClientId     string
	RedirectURI  string
	Scope        string
	State        string
}

func ParseAuthorizeRequest(values url.Values) (*AuthorizeRequest, *ParseError) {
	authorizeRequest := AuthorizeRequest{}

	fields := map[string]*mapping{
		"response_type": required(&authorizeRequest.ResponseType),
		"client_id":     required(&authorizeRequest.ClientId),
		"redirect_uri":  optional(&authorizeRequest.RedirectURI),
		"scope":         optional(&authorizeRequest.Scope),
		"state":         optional(&authorizeRequest.State),
	}

	parseError := ParseRequest(values, fields)

	if parseError != nil {
		return nil, parseError
	} else {
		return &authorizeRequest, nil
	}
}
