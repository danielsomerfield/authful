package admin

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/service/accesscontrol"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/wire"
	"io/ioutil"
	"log"
	"github.com/danielsomerfield/authful/server/handlers"
)

func NewRegisterClientHandler(
	clientAccessControlFn accesscontrol.ClientAccessControlFn,
	registerClientFn oauth.RegisterClientFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authorized, err := clientAccessControlFn(*r)

		if !authorized {
			if err != nil {
				handlers.InternalServerError("An unexpected error occurred", w)
			} else {
				handlers.Unauthorized("The requested operation was denied.", w)
			}
			return
		}

		registerClientRequest, err := ParseRegisterClientRequest(r)

		if err != nil {
			handlers.InvalidRequest("Failed to parse request to register client", w)
			return
		} else {
			credentials, err := registerClientFn(registerClientRequest.Name, registerClientRequest.Scopes)
			if err == nil {
				bytes, err := json.Marshal(wire.ResponseEnvelope{
					Data: RegisterClientResponse{
						ClientId:     credentials.ClientId,
						ClientSecret: credentials.ClientSecret,
					},
				})
				handlers.WriteOrError(w, bytes, err)
			} else {
				log.Printf("Failed to register the client: %+v", err)
				handlers.InternalServerError("An unexpected error occurred", w)
			}
		}

	}
}

func ParseRegisterClientRequest(request *http.Request) (*RegisterClientCommand, error) {
	body, err := ioutil.ReadAll(request.Body)
	var registerClientRequest *RegisterClientRequest = nil
	if err == nil {
		registerClientRequest = &RegisterClientRequest{}
		err = json.Unmarshal(body, &registerClientRequest)
	}

	return &registerClientRequest.Command, err
}

type RegisterClientRequest struct {
	Command RegisterClientCommand `json:"command,omitempty"`
}

type RegisterClientCommand struct {
	Name   string    `json:"name,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}

type RegisterClientResponse struct {
	ClientId     string    `json:"clientId,omitempty"`
	ClientSecret string    `json:"clientSecret,omitempty"`
}
