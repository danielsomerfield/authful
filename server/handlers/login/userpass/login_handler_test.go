package userpass

import (
	"testing"
	"net/http"
	"github.com/danielsomerfield/authful/server/handlers"
)

var handler = func(w http.ResponseWriter, r *http.Request) {}

func TestLoginHandler_loginSuccess(t *testing.T) {
	handlers.DoGetEndpointRequest(handler, "/login").
		ThenAssert(func(response *handlers.EndpointResponse) error {
			return nil
	}, t)
}
