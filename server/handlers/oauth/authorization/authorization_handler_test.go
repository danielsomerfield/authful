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
var validRedirect = "https://example.com/redirect"
var defaultRedirect = "https://example.com/defaultRedirect"

type MockClient struct {
}

func (MockClient) CheckSecret(secret string) bool {
	return secret == validClientSecret
}

func (mc MockClient) GetDefaultRedirectURI() string {
	return defaultRedirect
}

func (mc MockClient) GetScopes() []string {
	return []string{}
}

func (mc MockClient) IsValidRedirectURI(uri string) bool {
	return uri == validRedirect
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

func MockErrorPageRenderer(error string) []byte {
	return []byte(fmt.Sprintf("<html>%s</html>", error))
}

func TestAuthorizeHandler_successfulAuthorization(t *testing.T) {
	responseType := "code"
	state := "state1"
	redirectUri := "https://example.com/redirect"

	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s&redirect_uri=%s",
		validClientId, responseType, state, url.QueryEscape(redirectUri))

	handlers.DoGetEndpointRequest(NewAuthorizationHandler(MockClientLookupFn, MockCodeGenerator, MockErrorPageRenderer), requestUrl).
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
		util.AssertEquals("", params.Get("error"), t)

		return nil
	}, t)
}

func TestAuthorizeHandler_invalidClient(t *testing.T) {
	responseType := "code"
	state := "state1"
	redirectUri := "https://example.com/redirect"

	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s&redirect_uri=%s",
		invalidClientId, responseType, state, url.QueryEscape(redirectUri))

	handlers.DoGetEndpointRequest(NewAuthorizationHandler(MockClientLookupFn, MockCodeGenerator, MockErrorPageRenderer), requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		response.AssertHttpStatusEquals(200) //TODO: 200? Seems common, but seems wrong.
		response.AssertHeaderValue("Content-type", "text/html", t)
		response.AssertResponseContent("<html>unknown_client</html>", t)
		return nil
	}, t)
}

func TestAuthorizeHandler_badRedirectURL(t *testing.T) {
	responseType := "code"
	state := "state1"
	redirectUri := "https://example.com/badRedirect"

	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s&redirect_uri=%s",
		validClientId, responseType, state, url.QueryEscape(redirectUri))

	handlers.DoGetEndpointRequest(NewAuthorizationHandler(MockClientLookupFn, MockCodeGenerator, MockErrorPageRenderer), requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		response.AssertHttpStatusEquals(200)
		response.AssertHeaderValue("Content-type", "text/html", t)
		response.AssertResponseContent("<html>invalid_redirect_uri</html>", t)
		return nil
	}, t)
}

func TestAuthorizeHandler_noRedirectURL(t *testing.T) {
	responseType := "code"
	state := "state1"

	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s",
		validClientId, responseType, state)

	handlers.DoGetEndpointRequest(NewAuthorizationHandler(MockClientLookupFn, MockCodeGenerator, MockErrorPageRenderer), requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		response.AssertHttpStatusEquals(302)
		response.AssertHasHeader("location", t)

		uri, err := url.Parse(response.GetHeader("location"))
		util.AssertNoError(err, t)

		withoutQuery := fmt.Sprintf("%s://%s%s", uri.Scheme, uri.Host, uri.Path)
		util.AssertEquals("https://example.com/defaultRedirect", withoutQuery, t)

		params, err := url.ParseQuery(uri.RawQuery)
		util.AssertNoError(err, t)

		util.AssertEquals("a-code", params.Get("code"), t)
		util.AssertEquals("", params.Get("error"), t)

		return nil
	}, t)
}

func TestAuthorizeHandler_invalidScopesRequested(t *testing.T) {
	responseType := "code"
	state := "state1"
	redirectUri := "https://example.com/redirect"
	scope := "invalid-scope"

	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s&redirect_uri=%s&scope=%s",
		validClientId, responseType, state, url.QueryEscape(redirectUri), scope)

	handlers.DoGetEndpointRequest(NewAuthorizationHandler(MockClientLookupFn, MockCodeGenerator, MockErrorPageRenderer), requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		response.AssertHttpStatusEquals(200)
		response.AssertHeaderValue("Content-type", "text/html", t)
		response.AssertResponseContent("<html>invalid_scope</html>", t)
		return nil
	}, t)
}

//TODO: invalid request
//TODO: test for redirect uri with a ? already
