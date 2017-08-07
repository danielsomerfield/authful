package handlers

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/service/accesscontrol"
	"github.com/danielsomerfield/authful/server/service/oauth"
)

func Protect(handler http.HandlerFunc, scopes string, fn oauth.ClientLookupFn) http.HandlerFunc {
	acl := accesscontrol.NewClientAccessControlFnWithScopes(fn, scopes)
	return func(w http.ResponseWriter, r *http.Request) {
		ok, err := acl(*r)

		if err != nil {
			InternalServerError("An unexpected error occurred", w)
		} else if ok {
			handler(w, r)
		} else {
			Unauthorized("The requested operation was denied.", w)
		}
	}
}
