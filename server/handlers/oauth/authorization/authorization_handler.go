package authorization

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/wire/oauth"
	oauth_handlers "github.com/danielsomerfield/authful/server/handlers/oauth"

	oauth2 "github.com/danielsomerfield/authful/server/service/oauth"
	"log"
	"fmt"
	"net/url"
)

type CodeGenerator func() string

func NewAuthorizationHandler(clientLookup oauth2.ClientLookupFn, generator CodeGenerator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		authorizationRequest, err := oauth.ParseAuthorizeRequest(*r)
		if err != nil {
			log.Printf("Invalid authorization request due to error: %+v", err)
			oauth_handlers.InvalidRequest(err.Error(), w) //TODO: make this a redirect to error endpoint
			return
		} else {
			client, _ := clientLookup(authorizationRequest.ClientId)
			if client == nil {
				log.Printf("Request for unknown client %s.", authorizationRequest.ClientId)
				return
			}
		}
		http.Redirect(w, r, appendParam(authorizationRequest.RedirectURI, "code", generator()), http.StatusFound)

		//Reject if the redirect_uri doesn't match one configured with the client

		//Check scopes
		//Redirect to error if there is a scope in the request that's not in the client

		//Authenticate RO and ask for approval of request
	}
}

func appendParam(uri string, paramName string, paramValue string) string {
	return fmt.Sprintf("%s?%s=%s", uri, url.QueryEscape(paramName), url.QueryEscape(paramValue))
}
