package oauth

import (
	"net/http"
	"encoding/base64"
	"regexp"
	"errors"
	"log"
	"github.com/danielsomerfield/authful/server/wire"
)

var GRANT_TYPE_CLIENT_CREDENTIALS = "client_credentials"

type TokenRequest struct {
	GrantType    string
	Scope        string
	ClientId     string
	ClientSecret string
}

var ERR_INVALID_CLIENT = errors.New("invalid_client")

func ParseTokenRequest(httpRequest http.Request) (*TokenRequest, error) {

	if err := httpRequest.ParseForm(); err != nil {
		return nil, err
	}

	tokenRequest := TokenRequest{}

	//TODO: add support for required fields for other grant types
	//TODO: e.g. : need a way to make sure client_id and client_secret are there if grant_type == "client_credentials"
	fields := map[string]*wire.Mapping{
		"grant_type":    wire.Required(&tokenRequest.GrantType),
		"scope":         wire.Optional(&tokenRequest.Scope),
		"client_id":     wire.Optional(&tokenRequest.ClientId),
		"client_secret": wire.Optional(&tokenRequest.ClientSecret),
	}

	if err := wire.ParseRequest(httpRequest, fields); err != nil {
		return nil, err
	}

	authHeader := httpRequest.Header.Get("Authorization")
	if authHeader != "" {
		if tokenRequest.ClientId != "" {
			log.Print("Invalid client: credentials in both header and body.")
			return nil, ERR_INVALID_CLIENT
		}
		encodedToken := regexp.MustCompile("Basic ([a-zA-Z0-9]*)").FindStringSubmatch(string(authHeader))
		if len(encodedToken) > 1 {
			bearerBytes, err := base64.RawStdEncoding.DecodeString(encodedToken[1])
			if err != nil {
				return nil, err
			}
			bearerTokenString := string(bearerBytes)
			creds := regexp.MustCompile("(.*):(.*)").FindStringSubmatch(bearerTokenString)
			if len(creds) > 1 {
				tokenRequest.ClientId = creds[1]
				if len(creds) > 2 {
					tokenRequest.ClientSecret = creds[2]
				}
			}
		}

	}
	return &tokenRequest, nil
}
