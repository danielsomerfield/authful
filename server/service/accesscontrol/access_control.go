package accesscontrol

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/service/oauth"
)

type ClientAccessControlFn func(request http.Request) bool

func NewClientAccessControlFn(clientLookup oauth.ClientLookupFn) ClientAccessControlFn {
	return func(request http.Request) bool {
		//TODO: This will need to support two auth methods: client credentials and token
		//TODO: Implement client credentials first (token can come later)
		//Get the credentials from the request
		authorizationHeader := request.Header.Get("Authorization")
		if authorizationHeader == "" {
			return false
		} else {
			//Look up the client

			//client := clientLookup()
			//Make sure the client has the introspect_token or administrate scope
		}
		return true //TODO: NYI
	}
}
