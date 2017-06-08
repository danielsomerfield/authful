package token

import (
	"net/http"
	"log"
	"encoding/json"
	oauth_wire "github.com/danielsomerfield/authful/server/wire/oauth"
	oauth_service "github.com/danielsomerfield/authful/server/service/oauth"
	"strings"
	"time"
	"github.com/danielsomerfield/authful/server/handlers"
)

type CurrentTimeFn func() time.Time

type TokenGeneratorFn func() string

type TokenHandlerConfig struct {
	DefaultTokenExpiration int64
}

type StoreTokenFn func(token string, tokenMetaData oauth_service.TokenMetaData) error

func NewTokenHandler(
	config TokenHandlerConfig,
	clientLookupFn oauth_service.ClientLookupFn,
	tokenGenerator TokenGeneratorFn,
	storeTokenFn StoreTokenFn,
	currentTimeFn CurrentTimeFn) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		TokenHandler(w, req, config, clientLookupFn, tokenGenerator, storeTokenFn, currentTimeFn)
	}
}

func TokenHandler(w http.ResponseWriter, req *http.Request, config TokenHandlerConfig,
	clientLookupFn oauth_service.ClientLookupFn, tokenGenerator TokenGeneratorFn, storeTokenFn StoreTokenFn, currentTimeFn CurrentTimeFn) {

	if err := req.ParseForm(); err != nil {
		log.Printf("Failed with following error: %+v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	tokenRequest, err := oauth_wire.ParseTokenRequest(*req)
	if err != nil {
		if err == oauth_wire.ERR_INVALID_CLIENT {
			handlers.Unauthorized(err.Error(), w)
		} else {
			handlers.InvalidRequest(err.Error(), w)
		}
		return
	}

	client, err := clientLookupFn(tokenRequest.ClientId)
	if err != nil || client == nil || !client.CheckSecret(tokenRequest.ClientSecret) {
		if err != nil {
			log.Printf("Failure trying to look up client: %+v", err)
		} else if client == nil {
			log.Printf("Attempt to find invalid client by id \"%s\"", tokenRequest.ClientId)
		} else {
			log.Printf("Bad secret for client id \"%s\"", tokenRequest.ClientId)
		}
		handlers.JsonError("invalid_client", "Invalid client.", "", 401, w)
		return
	}

	if tokenRequest.GrantType == "client_credentials" {

		requestedScopes := strings.Fields(tokenRequest.Scope)
		unknownScopes := elementsNotIn(requestedScopes, client.GetScopes())
		if len(unknownScopes) > 0 {
			log.Printf("Request contained unexpected scopes: %+v", unknownScopes)
			handlers.JsonError("invalid_scope", "a requested was unknown", "",
				http.StatusBadRequest, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		token := tokenGenerator()

		storeTokenFn(token, oauth_service.TokenMetaData{
			Token:      token,
			Expiration: currentTimeFn().Add(time.Duration(config.DefaultTokenExpiration) * time.Second),
			ClientId:   tokenRequest.ClientId,
		})

		bytes, err := json.Marshal(oauth_wire.TokenResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   config.DefaultTokenExpiration,
			Scope:       strings.Join(requestedScopes, " "),
		})
		handlers.WriteOrError(w, bytes, err)
	} else {
		handlers.JsonError("unsupported_grant_type", "the grant type was missing or unknown", "",
			http.StatusBadRequest, w)
	}
}

func elementsNotIn(array []string, knownElements []string) []string {
	extraElements := []string{}

	for _, element := range array {
		if !contains(knownElements, element) {
			extraElements = append(extraElements, element)
		}
	}

	return extraElements
}

func contains(array []string, element string) bool {
	for _, e := range array {
		if element == e {
			return true;
		}
	}
	return false
}
