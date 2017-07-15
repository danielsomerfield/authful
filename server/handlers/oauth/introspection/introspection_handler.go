package introspection

import (
	"net/http"

	"github.com/danielsomerfield/authful/server/service/oauth"
	"encoding/json"
	"time"
	"github.com/danielsomerfield/authful/server/service/accesscontrol"
	oauth_handlers "github.com/danielsomerfield/authful/server/handlers/oauth"
)

func NewIntrospectionHandler(accessControlFn accesscontrol.ClientAccessControlFn, getTokenMetaData oauth.GetTokenMetaDataFn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		//TODO support for client credentials

		allowed, err := accessControlFn(*request)
		if !allowed {
			oauth_handlers.JsonOAuthError("invalid_token", "Failed to authenticate.", "",
				http.StatusUnauthorized, w)
			return
		}

		if err != nil {
			oauth_handlers.InvalidRequest(err.Error(), w)
			return
		}

		err = request.ParseForm()
		if err != nil {
			oauth_handlers.InvalidRequest(err.Error(), w)
			return
		}

		token := request.Form.Get("token")

		if token == "" {
			oauth_handlers.InvalidRequest("Missing field 'token'", w)
			return
		} else {
			tokenMetaData, err := getTokenMetaData(token)
			if err != nil {
				oauth_handlers.InvalidRequest(err.Error(), w)
				return
			}

			active := tokenMetaData != nil && isCurrent(*tokenMetaData)
			bytes, err := json.Marshal(IntrospectionResponse{
				Active: active,
			})
			oauth_handlers.WriteOrJsonOAuthError(w, bytes, err)
		}
	}
}

func isCurrent(tokenMetaData oauth.TokenMetaData) bool {
	return time.Now().Before(tokenMetaData.Expiration)
}

type IntrospectionResponse struct {
	Active bool `json:"active"`
}
