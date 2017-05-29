package oauth

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/request"
	"log"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/wireTypes"
	"fmt"
)

func TokenHandler(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		log.Printf("Failed with following error: %+v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	_, parseError := request.ParseTokenRequest(req.Form)
	if parseError != nil {
		invalidRequest(formatParseError(parseError), w)
		return
	}
	//Check that all scopes are known
	//Create the token in the backend
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(wireTypes.TokenResponse{
		AccessToken: "TODO",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	})
	writeOrError(w, bytes, err)
}

func AuthorizeHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	values := req.Form

	authorizationRequest, parseError := request.ParseAuthorizeRequest(values)
	if parseError != nil {
		invalidRequest(formatParseError(parseError), w)
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

type Client struct {
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

func formatParseError(error *request.ParseError) string {
	return fmt.Sprintf("The following fields are required: %s", error.MissingFields)
}