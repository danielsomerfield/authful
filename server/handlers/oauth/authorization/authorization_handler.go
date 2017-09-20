package authorization

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/wire/oauth"
	oauth_handlers "github.com/danielsomerfield/authful/server/handlers/oauth"

	oauth2 "github.com/danielsomerfield/authful/server/service/oauth"
	"log"
	"fmt"
	"net/url"
	"strings"
)

type CodeGenerator func() string
type ErrorPageRenderer func(error string) []byte

func NewAuthorizationHandler(clientLookup oauth2.ClientLookupFn, generator CodeGenerator,
	errorRenderer ErrorPageRenderer) func(http.ResponseWriter, *http.Request) {
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
				w.Header().Set("Content-type", "text/html")
				w.Write(errorRenderer("unknown_client"))
				return
			}

			scopes := strings.Fields(authorizationRequest.Scope)
			unknownScopes := elementsNotIn(scopes, client.GetScopes())
			if len(unknownScopes) > 0 {
				log.Printf("Request with invalid scope(s) %+v.", unknownScopes)
				w.Header().Set("Content-type", "text/html")
				w.Write(errorRenderer("invalid_scope"))
				return
			}

			var redirectURI string
			if authorizationRequest.RedirectURI == "" {
				redirectURI = client.GetDefaultRedirectURI()
			} else if client.IsValidRedirectURI(authorizationRequest.RedirectURI) {
				redirectURI = authorizationRequest.RedirectURI
			}

			if redirectURI == "" {
				log.Printf("Request with invalid redirect uri %s.", authorizationRequest.RedirectURI)
				w.Header().Set("Content-type", "text/html")
				w.Write(errorRenderer("invalid_redirect_uri"))
				return
			} else {
				http.Redirect(w, r, appendParam(redirectURI, "code", generator()), http.StatusFound)
			}
		}

		//Redirect to error if there is a scope in the request that's not in the client

		//Authenticate RO and ask for approval of request
	}
}

//TODO: refactor out
func elementsNotIn(array []string, knownElements []string) []string {
	extraElements := []string{}

	for _, element := range array {
		if !contains(knownElements, element) {
			extraElements = append(extraElements, element)
		}
	}

	return extraElements
}

func contains(array []string, element string) bool {
	for _, e := range array {
		if element == e {
			return true;
		}
	}
	return false
}

func appendParam(uri string, paramName string, paramValue string) string {
	return fmt.Sprintf("%s?%s=%s", uri, url.QueryEscape(paramName), url.QueryEscape(paramValue))
}
