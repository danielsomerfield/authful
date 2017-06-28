package introspection

import (
	"net/http"

	"github.com/danielsomerfield/authful/server/handlers"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"encoding/json"
)

type RequestValidationFn func(request http.Request) bool
type GetTokenMetaDataFn func(token string) *oauth.TokenMetaData

func NewIntrospectionHandler(validation RequestValidationFn, getTokenMetaData GetTokenMetaDataFn) func(http.ResponseWriter, *http.Request) {
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
			tokenMetaData := getTokenMetaData(token)
			if tokenMetaData != nil {
				bytes, err := json.Marshal(IntrospectionResponse{
					Active: true,
				})
				handlers.WriteOrError(w, bytes, err)
			}
		}
	}
}

type IntrospectionResponse struct {
	Active bool `json:"active"`
}
