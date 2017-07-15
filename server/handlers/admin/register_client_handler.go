package admin

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/service/accesscontrol"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/wire"
	"io/ioutil"
	"log"
)

func NewRegisterClientHandler(
	clientAccessControlFn accesscontrol.ClientAccessControlFn,
	registerClientFn oauth.RegisterClientFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authorized, err := clientAccessControlFn(*r)

		if !authorized {
			if err != nil {
				InternalServerError("An unexpected error occurred", w)
			} else {
				Unauthorized("The requested operation was defined.", w)
			}
			return
		}

		registerClientRequest, err := ParseRegisterClientRequest(r)

		if err != nil {
			InvalidRequest("Failed to parse request to register client", w)
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
				WriteOrError(w, bytes, err)
			} else {
				log.Printf("Failed to register the client: %+v", err)
				InternalServerError("An unexpected error occurred", w)
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

//TODO: refactor these with oauth_handler_utils
func InvalidRequest(errorDescription string, w http.ResponseWriter) {
	JsonError("invalid_request", errorDescription, "", http.StatusBadRequest, w)
}

func Unauthorized(errorDescription string, w http.ResponseWriter) {
	JsonError("invalid_client", errorDescription, "", http.StatusUnauthorized, w)
}

func InternalServerError(errorDescription string, w http.ResponseWriter) {
	JsonError("server_error", errorDescription, "", http.StatusInternalServerError, w)
}

func JsonError(errorType string, errorDescription string, errorURI string, httpStatus int, w http.ResponseWriter) {
	w.WriteHeader(httpStatus)
	errorMessageJSON, err := json.Marshal(wire.ErrorsResponse{})
	if err == nil {
		w.Write(errorMessageJSON)
	} else {
		log.Printf("Failed to write error message: %+v", err)
	}
}

func WriteOrError(w http.ResponseWriter, bytes []byte, err error) {
	if err == nil {
		w.Write(bytes)
	} else {
		log.Printf("Failed with following error: %+v", err)
		JsonError("unknown", "an unexpected error occurred", "",
			http.StatusInternalServerError, w)
	}
}
