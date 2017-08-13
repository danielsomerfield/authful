package authorization

import (
	"testing"
	"github.com/danielsomerfield/authful/server/handlers"
	"fmt"
	"net/url"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"github.com/danielsomerfield/authful/util"
)

var validClientId = "valid-client-id"
var validClientSecret = "valid-client-secret"
var invalidClientId = "invalid-client-id"

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

func MockCodeGenerator() string {
	return "a-code"
}

func TestAuthorizeHandler_successfulAuthorization(t *testing.T) {
	responseType := "code"
	state := "state1"
	redirectUri := "https://example.com/redirect"

	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s&redirect_uri=%s",
		validClientId, responseType, state, url.QueryEscape(redirectUri))

	handlers.DoGetEndpointRequest(NewAuthorizationHandler(MockClientLookupFn, MockCodeGenerator), requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		response.AssertHttpStatusEquals(302)
		response.AssertHasHeader("location", t)
		uri, err := url.Parse(response.GetHeader("location"))
		util.AssertNoError(err, t)

		withoutQuery := fmt.Sprintf("%s://%s%s", uri.Scheme, uri.Host, uri.Path)
		util.AssertEquals("https://example.com/redirect", withoutQuery, t)

		params, err := url.ParseQuery(uri.RawQuery)
		util.AssertNoError(err, t)

		util.AssertEquals("a-code", params.Get("code"), t)

		return nil
	}, t)
}

//func TestAuthorizeHandler_invalidClient(t *testing.T) {
//	responseType := "code"
//	state := "state1"
//	redirectUri := "https://example.com?redirect"
//
//	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s&redirect_uri=%s",
//		invalidClientId, responseType, state, url.QueryEscape(redirectUri))
//
//	handlers.DoGetEndpointRequest(NewAuthorizationHandler(MockClientLookupFn), requestUrl).
//		ThenAssert(func(response *handlers.EndpointResponse) error {
//		response.AssertHttpStatusEquals(302)
//		response.AssertHeaderValue("location", redirectUri, t)
//		return nil
//	}, t)
//}

//TODO: bad redirect url
//TODO: test with no redirect url (redirects to default)
//TODO: invalid scope for client
//TODO: invalid request
//TODO: test for redirect uri with a ? already
