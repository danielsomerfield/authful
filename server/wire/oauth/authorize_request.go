package oauth

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/wire"
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

	fields := map[string]*wire.Mapping{
		"response_type": wire.Required(&authorizeRequest.ResponseType),
		"client_id":     wire.Required(&authorizeRequest.ClientId),
		"redirect_uri":  wire.Optional(&authorizeRequest.RedirectURI),
		"scope":         wire.Optional(&authorizeRequest.Scope),
		"state":         wire.Optional(&authorizeRequest.State),
	}

	if err := wire.ParseRequest(httpRequest, fields); err != nil {
		return nil, err
	}
	return &authorizeRequest, nil
}
