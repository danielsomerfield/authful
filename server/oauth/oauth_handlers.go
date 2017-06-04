package oauth

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/request"
	"log"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/wireTypes"
)

func AuthorizeHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	authorizationRequest, err := request.ParseAuthorizeRequest(*req)
	if err != nil {
		InvalidRequest(err.Error(), w)
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

func getClient(clientId string) *Client {
	return nil
}

func InvalidRequest(errorDescription string, w http.ResponseWriter) {
	JsonError("invalid_request", errorDescription, "", http.StatusBadRequest, w)
}

func Unauthorized(errorDescription string, w http.ResponseWriter) {
	JsonError("invalid_client", errorDescription, "", http.StatusUnauthorized, w)
}

func JsonError(errorType string, errorDescription string, errorURI string, httpStatus int, w http.ResponseWriter) {
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

func WriteOrError(w http.ResponseWriter, bytes []byte, err error) {
	if err == nil {
		w.Write(bytes)
	} else {
		log.Printf("Failed with following error: %+v", err)
		JsonError("unknown", "an unexpected error occurred", "",
			http.StatusInternalServerError, w)
	}
}
