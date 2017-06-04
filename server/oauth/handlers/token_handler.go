package handlers

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/request"
	"log"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/wireTypes"
	"strings"
	"github.com/danielsomerfield/authful/server/oauth"
	"time"
)

type CurrentTimeFn func() time.Time

type TokenGeneratorFn func() string

type TokenHandlerConfig struct {
	DefaultTokenExpiration int64
}

type TokenMetaData struct {
	token string
	expiration time.Time
	clientId string
}

type TokenStore interface {
	StoreToken(token string, tokenMetaData TokenMetaData)
}

func NewTokenHandler(
	config TokenHandlerConfig,
	clientLookupFn oauth.ClientLookupFn,
	tokenGenerator TokenGeneratorFn,
	tokenStore TokenStore,
	currentTimeFn CurrentTimeFn) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		TokenHandler(w, req, config, clientLookupFn, tokenGenerator, tokenStore, currentTimeFn)
	}
}

func TokenHandler(w http.ResponseWriter, req *http.Request, config TokenHandlerConfig,
	clientLookupFn oauth.ClientLookupFn, tokenGenerator TokenGeneratorFn, tokenStore TokenStore, currentTimeFn CurrentTimeFn) {

	if err := req.ParseForm(); err != nil {
		log.Printf("Failed with following error: %+v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	tokenRequest, err := request.ParseTokenRequest(*req)
	if err != nil {
		if err == request.ERR_INVALID_CLIENT {
			oauth.Unauthorized(err.Error(), w)
		} else {
			oauth.InvalidRequest(err.Error(), w)
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
		oauth.JsonError("invalid_client", "Invalid client.", "", 401, w)
		return
	}

	if tokenRequest.GrantType == "client_credentials" {

		requestedScopes := strings.Fields(tokenRequest.Scope)
		unknownScopes := elementsNotIn(requestedScopes, client.GetScopes())
		if len(unknownScopes) > 0 {
			log.Printf("Request contained unexpected scopes: %+v", unknownScopes)
			oauth.JsonError("invalid_scope", "a requested was unknown", "",
				http.StatusBadRequest, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		token := tokenGenerator()

		tokenStore.StoreToken(token, TokenMetaData{
			token: token,
			expiration: currentTimeFn().Add(time.Duration(config.DefaultTokenExpiration) * time.Second),
			clientId: tokenRequest.ClientId,
		})

		bytes, err := json.Marshal(wireTypes.TokenResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   config.DefaultTokenExpiration,
			Scope:       strings.Join(requestedScopes, " "),
		})
		oauth.WriteOrError(w, bytes, err)
	} else {
		oauth.JsonError("unsupported_grant_type", "the grant type was missing or unknown", "",
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
