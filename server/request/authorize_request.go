package request

import (
	"net/http"
)

type AuthorizeRequest struct {
	ResponseType string
	ClientId     string
	RedirectURI  string
	Scope        string
	State        string
}

func ParseAuthorizeRequest(httpRequest http.Request) (*AuthorizeRequest, error) {
	authorizeRequest := AuthorizeRequest{}

	fields := map[string]*mapping{
		"response_type": required(&authorizeRequest.ResponseType),
		"client_id":     required(&authorizeRequest.ClientId),
		"redirect_uri":  optional(&authorizeRequest.RedirectURI),
		"scope":         optional(&authorizeRequest.Scope),
		"state":         optional(&authorizeRequest.State),
	}

	if err := ParseRequest(httpRequest, fields); err != nil {
		return nil, err
	}
	return &authorizeRequest, nil
}
