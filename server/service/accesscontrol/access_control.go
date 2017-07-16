package accesscontrol

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"github.com/danielsomerfield/authful/server/wire/authentication"
)

type ClientAccessControlFn func(request http.Request) (bool, error)

func NewClientAccessControlFn(clientLookup oauth.ClientLookupFn) ClientAccessControlFn {
	return NewClientAccessControlFnWithScopes(clientLookup, "")
}

func NewClientAccessControlFnWithScopes(clientLookup oauth.ClientLookupFn, requiredScope string) ClientAccessControlFn {
	return func(request http.Request) (bool, error) {
		//TODO: This will need to support two auth methods: client credentials and token
		//TODO: Implement client credentials first (token can come later)
		//Get the credentials from the request

		credentials, err := authentication.ParseClientCredentialsBasicHeader(request)
		if credentials != nil {
			client, err := clientLookup(credentials.ClientId)
			if client != nil {
				return client.CheckSecret(credentials.ClientSecret) &&
					(requiredScope == "" || HasScope(client, requiredScope)), err
			}
		}

		return false, err //TODO: NYI
	}
}

func HasScope(client oauth.Client, scope string) bool {
	if scope == "" {
		return true
	}
	scopes := client.GetScopes()
	for i := range scopes {
		if scopes[i] == scope {
			return true
		}
	}
	return false
}
