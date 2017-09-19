package client

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/wire"
	"io/ioutil"
	"log"
	"github.com/danielsomerfield/authful/server/handlers"
	"strings"
	"fmt"
)

func NewProtectedHandler(registerClientFn oauth.RegisterClientFn, lookup oauth.ClientLookupFn) http.HandlerFunc {
	return handlers.Protect(NewRegisterClientHandler(registerClientFn), "administrate", lookup)
}

func NewRegisterClientHandler(
	registerClientFn oauth.RegisterClientFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		registerClientRequest, err := ParseRegisterClientRequest(r)

		if err != nil {
			handlers.InvalidRequest("Failed to parse request to register client", w)
			return
		} else {
			credentials, err := registerClientFn(registerClientRequest.Name, registerClientRequest.Scopes,
				registerClientRequest.RedirectUris, registerClientRequest.DefaultRedirectURI)
			if err == nil {
				bytes, err := json.Marshal(wire.ResponseEnvelope{
					Data: RegisterClientResponse{
						ClientId:     credentials.ClientId,
						ClientSecret: credentials.ClientSecret,
					},
				})
				handlers.WriteOrInternalError(w, bytes, err)
			} else {
				log.Printf("Failed to register the client: %+v", err)
				handlers.InternalServerError("An unexpected error occurred", w)
			}
		}

	}
}

//TODO: refactor this bit to the wire package
func ParseRegisterClientRequest(request *http.Request) (*RegisterClientCommand, error) {
	body, err := ioutil.ReadAll(request.Body)
	var registerClientRequest *RegisterClientRequest = nil
	if err == nil {
		registerClientRequest = &RegisterClientRequest{}
		err = json.Unmarshal(body, &registerClientRequest)
	}

	if strings.TrimSpace(registerClientRequest.Command.Name) == "" {
		return nil, fmt.Errorf("Missing required field \"name\"")
	}

	return &registerClientRequest.Command, err
}

type RegisterClientRequest struct {
	Command RegisterClientCommand `json:"command,omitempty"`
}

type RegisterClientCommand struct {
	Name         string    `json:"name,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
	RedirectUris []string `json:"redirect_uris,omitempty"`
	DefaultRedirectURI string `json:"default_redirect_uri,omitempty"`
}

type RegisterClientResponse struct {
	ClientId     string    `json:"clientId,omitempty"`
	ClientSecret string    `json:"clientSecret,omitempty"`
}
