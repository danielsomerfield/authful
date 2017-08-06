package authorization

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/wire/oauth"
	oauth_handlers "github.com/danielsomerfield/authful/server/handlers/oauth"

	oauth2 "github.com/danielsomerfield/authful/server/service/oauth"
	"log"
)

func NewAuthorizationHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()

		authorizationRequest, err := oauth.ParseAuthorizeRequest(*req)
		if err != nil {
			log.Printf("Invalid authorization request due to error: %+v", err)
			oauth_handlers.InvalidRequest(err.Error(), w)
			return
		} else {
			client := getClient(authorizationRequest.ClientId)
			if client == nil {
				//TODO: write back 401 and {"error": "invalid_client"}
				return
			}
		}

		//Reject if the redirect_uri doesn't match one configured with the client

		//Check scopes
		//Redirect to error if there is a scope in the request that's not in the client

		//Authenticate RO and ask for approval of request
	}
}

func getClient(clientId string) oauth2.Client {
	return nil
}
