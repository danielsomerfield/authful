package authorization

import (
	"testing"
	"github.com/danielsomerfield/authful/server/handlers"
	"fmt"
	"net/url"
	"github.com/danielsomerfield/authful/server/service/oauth"
)

var validClientId = "valid-client-id"
var validClientSecret = "valid-client-secret"
var invalidClientSecret = "invalid-client-secret"

type MockClient struct {

}

func (MockClient) CheckSecret(secret string) bool {
	return secret == validClientSecret
}

func (mc MockClient) GetScopes() []string {
	return []string{}
}

func MockClientLookupFn(clientId string) (oauth.Client, error) {
	if clientId == validClientId {
		return MockClient{

		}, nil
	} else {
		return nil, nil
	}
}

func TestAuthorizeHandler_successfulAuthorization(t *testing.T) {
	responseType := "code"
	state := "state1"
	redirectUri := "https://example.com?redirect"

	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s&redirect_uri=%s",
		validClientId, responseType, state, url.QueryEscape(redirectUri))

	handlers.DoGetEndpointRequest(NewAuthorizationHandler(MockClientLookupFn), requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		fmt.Printf("response: %+v", response)
		response.AssertHttpStatusEquals(302)
		return nil
	}, t)
}

//TODO: invalid/non-existent client id
//TODO: bad redirect url
//TODO: test with no redirect url (redirects to default)
//TODO: invalid scope for client
//TODO: invalid request
