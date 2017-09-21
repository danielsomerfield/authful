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
type ApprovalRequestStore func(request *oauth.AuthorizeRequest) string
type ApprovalLookup func(approvalType string, requestId string) *url.URL

func NewAuthorizationHandler(
	clientLookup oauth2.ClientLookupFn,
	errorRenderer ErrorPageRenderer,
	approvalRequestStore ApprovalRequestStore,
	approvalLookup ApprovalLookup) http.HandlerFunc {
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
				if err == nil {
					populateDefaultUri(authorizationRequest, client)
					requestId := approvalRequestStore(authorizationRequest)
					loginRedirect := approvalLookup("username-password", requestId)
					http.Redirect(w, r, loginRedirect.String(), http.StatusFound)
				} else {
					//TODO: refactor this so it doesn't check twice
					log.Printf("Request with invalid redirect uri %s.", authorizationRequest.RedirectURI)
					writeROError("invalid_redirect_uri", w, errorRenderer)
					return
				}
			}
		}
	}
}
func populateDefaultUri(request *oauth.AuthorizeRequest, client oauth2.Client) {
	if request.RedirectURI == "" {
		request.RedirectURI = client.GetDefaultRedirectURI()
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
