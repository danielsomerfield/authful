package oauth

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/request"
	"log"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/wireTypes"
)

type TokenGeneratorFn func() string

type TokenHandlerConfig struct {
	DefaultTokenExpiration float64
}

func NewTokenHandler(
	config TokenHandlerConfig,
	clientLookup ClientLookupFn,
	tokenGenerator TokenGeneratorFn) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		TokenHandler(w, req, config, clientLookup, tokenGenerator)
	}
}

func TokenHandler(w http.ResponseWriter, req *http.Request, config TokenHandlerConfig,
	clientLookup ClientLookupFn, tokenGenerator TokenGeneratorFn) {

	if err := req.ParseForm(); err != nil {
		log.Printf("Failed with following error: %+v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	tokenRequest, err := request.ParseTokenRequest(*req)
	if err != nil {
		invalidRequest(err.Error(), w)
		return
	}

	client, err := clientLookup(tokenRequest.ClientId)
	if client == nil || !client.checkSecret(tokenRequest.ClientSecret) {
		jsonError("invalid_client", "Invalid client.", "", 401, w)
		return
	}

	if tokenRequest.GrantType == "client_credentials" {
		//Check that all scopes are known
		//Create the token in the backend
		w.Header().Set("Content-Type", "application/json")
		bytes, err := json.Marshal(wireTypes.TokenResponse{
			AccessToken: tokenGenerator(),
			TokenType:   "Bearer",
			ExpiresIn:   config.DefaultTokenExpiration,
		})
		writeOrError(w, bytes, err)
	} else {
		jsonError("unsupported_grant_type", "the grant type was missing or unknown", "",
			http.StatusBadRequest, w)
	}

}

func AuthorizeHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	authorizationRequest, err := request.ParseAuthorizeRequest(*req)
	if err != nil {
		invalidRequest(err.Error(), w)
		return
	} else {
		client := getClient(authorizationRequest.ClientId)
		if client == nil {
			//TODO: write back 401 and {"error": "invalid_client"}
			return;
		}
	}

	//Reject if the redirect_uri doesn't match one configured with the client

	//Check scopes
	//Redirect to error if there is a scope in the request that's not in the client

	//Authenticate RO and ask for approval of request
}

type Client interface {
	checkSecret(secret string) bool
}

func getClient(clientId string) *Client {
	return nil
}

func invalidRequest(errorDescription string, w http.ResponseWriter) {
	jsonError("invalid_request", errorDescription, "", http.StatusBadRequest, w)
}

func jsonError(errorType string, errorDescription string, errorURI string, httpStatus int, w http.ResponseWriter) {
	w.WriteHeader(httpStatus)
	errorMessageJSON, err := json.Marshal(wireTypes.ErrorResponse{
		Error:            errorType,
		ErrorDescription: errorDescription,
		ErrorURI:         errorURI,
	})
	if err == nil {
		w.Write(errorMessageJSON)
	} else {
		log.Printf("Failed to write error message: %+v", err)
	}
}

func writeOrError(w http.ResponseWriter, bytes []byte, err error) {
	if err == nil {
		w.Write(bytes)
	} else {
		log.Printf("Failed with following error: %+v", err)
		jsonError("unknown", "an unexpected error occurred", "",
			http.StatusInternalServerError, w)
	}
}
