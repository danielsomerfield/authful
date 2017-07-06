package introspection

import (
	"net/http"

	"github.com/danielsomerfield/authful/server/handlers"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"encoding/json"
	"time"
)

type AccessControlFunction func(request http.Request) bool

func NewIntrospectionHandler(validation AccessControlFunction, getTokenMetaData oauth.GetTokenMetaDataFn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		//TODO support for client credentials

		if !validation(*request) {
			handlers.JsonError("invalid_token", "Failed to authenticate.", "",
				http.StatusUnauthorized, w)
			return
		}

		err := request.ParseForm()
		if err != nil {
			handlers.InvalidRequest(err.Error(), w)
			return
		}

		token := request.Form.Get("token")

		if token == "" {
			handlers.InvalidRequest("Missing field 'token'", w)
			return
		} else {
			tokenMetaData, err := getTokenMetaData(token)
			if err != nil {
				handlers.InvalidRequest(err.Error(), w)
				return
			}

			active := tokenMetaData != nil && isCurrent(*tokenMetaData)
			bytes, err := json.Marshal(IntrospectionResponse{
				Active: active,
			})
			handlers.WriteOrError(w, bytes, err)
		}
	}
}

func isCurrent(tokenMetaData oauth.TokenMetaData) bool {
	return time.Now().Before(tokenMetaData.Expiration)
}

type IntrospectionResponse struct {
	Active bool `json:"active"`
}
