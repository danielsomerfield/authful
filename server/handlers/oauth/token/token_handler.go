package token

import (
	"net/http"
	"log"
	"encoding/json"
	oauth_wire "github.com/danielsomerfield/authful/server/wire/oauth"
	oauth_service "github.com/danielsomerfield/authful/server/service/oauth"
	"strings"
	"time"
	oauth_handlers "github.com/danielsomerfield/authful/server/handlers/oauth"
	"github.com/danielsomerfield/authful/common/util"
)

type CurrentTimeFn func() time.Time

type TokenGeneratorFn func() string

type TokenHandlerConfig struct {
	DefaultTokenExpiration int64
}

func NewTokenHandler(
	config TokenHandlerConfig,
	clientLookupFn oauth_service.ClientLookupFn,
	tokenGenerator TokenGeneratorFn,
	storeTokenFn oauth_service.StoreTokenMetaDataFn,
	currentTimeFn CurrentTimeFn) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		TokenHandler(w, req, config, clientLookupFn, tokenGenerator, storeTokenFn, currentTimeFn)
	}
}

func TokenHandler(w http.ResponseWriter,
	req *http.Request,
	config TokenHandlerConfig,
	clientLookupFn oauth_service.ClientLookupFn,
	tokenGenerator TokenGeneratorFn,
	storeTokenFn oauth_service.StoreTokenMetaDataFn,
	currentTimeFn CurrentTimeFn) {

	tokenRequest, err := oauth_wire.ParseTokenRequest(*req)
	if err != nil {
		if err == oauth_wire.ERR_INVALID_CLIENT {
			oauth_handlers.Unauthorized(err.Error(), w)
		} else {
			oauth_handlers.InvalidRequest(err.Error(), w)
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
		oauth_handlers.JsonOAuthError("invalid_client", "Invalid client.", "", 401, w)
		return
	}

	if tokenRequest.GrantType == "client_credentials" {

		requestedScopes := strings.Fields(tokenRequest.Scope)
		unknownScopes := util.ElementsNotIn(requestedScopes, client.GetScopes())
		if len(unknownScopes) > 0 {
			log.Printf("Request contained unexpected scopes: %+v", unknownScopes)
			oauth_handlers.JsonOAuthError("invalid_scope", "a requested was unknown", "",
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
		oauth_handlers.WriteOrJsonOAuthError(w, bytes, err)
	} else {
		oauth_handlers.JsonOAuthError("unsupported_grant_type", "the grant type was missing or unknown", "",
			http.StatusBadRequest, w)
	}
}