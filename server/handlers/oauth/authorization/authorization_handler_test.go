package authorization

import (
	"testing"
	"github.com/danielsomerfield/authful/server/handlers"
	"fmt"
	"net/url"
)

var validClientId = "12345"

func mockClientLookup(clientId string)  {

}

func TestAuthorizeHandler_successfulAuthorization(t *testing.T) {
	clientId := "12345"
	responseType := "code"
	state := "state1"
	redirectUri := "https://example.com?redirect"

	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s&redirect_uri=%s",
		clientId, responseType, state, url.QueryEscape(redirectUri))

	handlers.DoGetEndpointRequest(NewAuthorizationHandler(), requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		response.AssertHttpStatusEquals(302)
		return nil
	}, t)
}

//TODO: invalid client id
//TODO: bad redirect url
//TODO: test with no redirect url (redirects to default)
//TODO: invalid scope for client
//TODO: invalid request
