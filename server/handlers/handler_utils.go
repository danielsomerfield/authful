package handlers

import (
	"github.com/danielsomerfield/authful/server/wire"
	"log"
	"net/http"
	"encoding/json"
)

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
	errorMessageJSON, err := json.Marshal(wire.ErrorsResponse{
		Errors: []wire.Error{
			{
				Status:    httpStatus,
				ErrorType: errorType,
				Detail:    errorDescription,
				ErrorURI:  errorURI,
			},
		},
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
