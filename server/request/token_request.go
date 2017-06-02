package request

import (
	"net/http"
	"encoding/base64"
	"regexp"
	"errors"
	"log"
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

	authHeader := httpRequest.Header.Get("Authorization")
	if authHeader != "" {
		if tokenRequest.ClientId != "" {
			log.Print("Invalid client: credentials in both header and body.")
			return nil, ERR_INVALID_CLIENT
		}
		encodedToken := regexp.MustCompile("Bearer ([a-zA-Z0-9]*)").FindStringSubmatch(string(authHeader))
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
