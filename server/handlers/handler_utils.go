package handlers

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/danielsomerfield/authful/common/wire"
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
		Error: wire.Error{
			Status:    httpStatus,
			ErrorType: errorType,
			Detail:    errorDescription,
			ErrorURI:  errorURI,
		},
	})
	if err == nil {
		w.Write(errorMessageJSON)
	} else {
		log.Printf("Failed to write error message: %+v", err)
	}
}

func WriteOrError(w http.ResponseWriter, bytes []byte, err error, errorCode int, errorType string, errorDescription string, errorURI string) {
	if err == nil {
		w.Write(bytes)
	} else {
		log.Printf("Failed with following error: %+v", err)
		JsonError(errorType, errorDescription, errorURI,
			errorCode, w)
	}
}

func WriteOrInternalError(w http.ResponseWriter, bytes []byte, err error) {
	WriteOrError(w, bytes, err, http.StatusInternalServerError, "unknown", "an unexpected error occurred", "")
}
