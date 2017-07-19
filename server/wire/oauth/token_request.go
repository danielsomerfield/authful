package oauth

import (
	"net/http"
	"errors"
	"log"
	"github.com/danielsomerfield/authful/server/wire"
	"github.com/danielsomerfield/authful/server/wire/authentication"
	"strings"
)

var GRANT_TYPE_CLIENT_CREDENTIALS = "client_credentials"

type TokenRequest struct {
	GrantType    string
	Scope        string
	ClientId     string
	ClientSecret string
}

var ERR_INVALID_CLIENT = errors.New("invalid_client")
var ERR_INVALID_REQUEST = errors.New("invalid_request")

func ParseTokenRequest(httpRequest http.Request) (*TokenRequest, error) {

	if err := httpRequest.ParseForm(); err != nil {
		log.Printf("Failed to parse form due to error: %+v", err)
		return nil, ERR_INVALID_REQUEST
	}

	tokenRequest := TokenRequest{}

	//TODO: add support for required fields for other grant types
	fields := map[string]*wire.Mapping{
		"grant_type":    wire.Required(&tokenRequest.GrantType),
		"scope":         wire.Optional(&tokenRequest.Scope),
		"client_id":     wire.Optional(&tokenRequest.ClientId),
		"client_secret": wire.Optional(&tokenRequest.ClientSecret),
	}

	if err := wire.ParseRequest(httpRequest, fields); err != nil {
		return nil, err
	}

	clientCredentials, err := authentication.ParseClientCredentialsBasicHeader(httpRequest)
	if clientCredentials != nil {
		if tokenRequest.ClientId != "" {
			log.Print("Invalid client: credentials in both header and body.")
			return nil, ERR_INVALID_CLIENT
		} else {
			tokenRequest.ClientId = clientCredentials.ClientId
			tokenRequest.ClientSecret = clientCredentials.ClientSecret
		}
	}

	if tokenRequest.GrantType == GRANT_TYPE_CLIENT_CREDENTIALS {
		if strings.TrimSpace(tokenRequest.ClientId) == "" || strings.TrimSpace(tokenRequest.ClientSecret) == "" {
			return nil, ERR_INVALID_REQUEST
		}
	}

	return &tokenRequest, err
}
