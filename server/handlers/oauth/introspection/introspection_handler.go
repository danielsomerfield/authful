package introspection

import (
	"net/http"

	"github.com/danielsomerfield/authful/server/handlers"
	"github.com/danielsomerfield/authful/server/service/oauth"
)

type RequestValidationFn func(request http.Request) bool
type GetTokenMetaDataFn func(token string) oauth.TokenMetaData

func NewIntrospectionHandler(validation RequestValidationFn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		//TODO support for client credentials

		if !validation(*request) {
			handlers.JsonError("invalid_token", "Failed to authenticate.", "",
				http.StatusUnauthorized, w)
			return
		}
	}
}
