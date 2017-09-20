package authorization

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/wire/oauth"
	oauth2 "github.com/danielsomerfield/authful/server/service/oauth"
	"log"
	"fmt"
	"net/url"
	"strings"
	"github.com/danielsomerfield/authful/common/util"
)

type CodeGenerator func() string
type ErrorPageRenderer func(error string) []byte

func NewAuthorizationHandler(clientLookup oauth2.ClientLookupFn, generator CodeGenerator,
	errorRenderer ErrorPageRenderer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		authorizationRequest, err := oauth.ParseAuthorizeRequest(*r)
		if err != nil {
			log.Printf("Invalid request from client %+v.", err)
			writeROError("invalid_request", w, errorRenderer)
			return
		} else {
			client, _ := clientLookup(authorizationRequest.ClientId)
			if client == nil {
				log.Printf("Request for unknown client %s.", authorizationRequest.ClientId)
				writeROError("unknown_client", w, errorRenderer)
				return
			}

			unknownScopes := util.ElementsNotIn(strings.Fields(authorizationRequest.Scope), client.GetScopes())
			if len(unknownScopes) > 0 {
				log.Printf("Request with invalid scope(s) %+v.", unknownScopes)
				writeROError("invalid_scope", w, errorRenderer)
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
				writeROError("invalid_redirect_uri", w, errorRenderer)
				return
			} else {
				uri, err := appendParam(redirectURI, "code", generator())
				if err == nil {
					http.Redirect(w, r, uri, http.StatusFound)
				} else {
					//TODO: refactor this so it doesn't check twice
					log.Printf("Request with invalid redirect uri %s.", authorizationRequest.RedirectURI)
					writeROError("invalid_redirect_uri", w, errorRenderer)
					return
				}
			}
		}

		//Authenticate RO and ask for approval of request
	}
}

func writeROError(errorCode string, w http.ResponseWriter, errorRenderer ErrorPageRenderer) {
	w.Header().Set("Content-type", "text/html")
	w.Write(errorRenderer(errorCode))
}

func appendParam(uri string, paramName string, paramValue string) (string, error) {
	redirectUri, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	} else if redirectUri.RawQuery != "" {
		return fmt.Sprintf("%s&%s=%s", uri, url.QueryEscape(paramName), url.QueryEscape(paramValue)), nil
	} else {
		return fmt.Sprintf("%s?%s=%s", uri, url.QueryEscape(paramName), url.QueryEscape(paramValue)), nil
	}
}
