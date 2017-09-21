package authorization

import (
	"testing"
	"github.com/danielsomerfield/authful/server/handlers"
	"fmt"
	"net/url"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"github.com/danielsomerfield/authful/util"
	util2 "github.com/danielsomerfield/authful/common/util"
	oauth2 "github.com/danielsomerfield/authful/server/wire/oauth"
	"os"
)

var validClientId = "valid-client-id"
var validClientSecret = "valid-client-secret"
var invalidClientId = "invalid-client-id"
var defaultRedirect = "https://example.com/defaultRedirect"
var validClient = MockClient{
	validRedirects: []string{"https://example.com/redirect"},
}

type MockClient struct {
	validRedirects []string
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
	return util2.Contains(mc.validRedirects, uri)
}

func MockClientLookupFn(clientId string) (oauth.Client, error) {
	if clientId == validClientId {
		return validClient, nil
	} else {
		return nil, nil
	}
}

func MockErrorPageRenderer(error string) []byte {
	return []byte(fmt.Sprintf("<html>%s</html>", error))
}

var approvalRequests map[string]*oauth2.AuthorizeRequest

func mockApprovalRequestStore(request *oauth2.AuthorizeRequest) string {
	approvalRequests["random-request-id"] = request
	return "random-request-id"
}

func mockApprovalLookup(approvalType string, requestId string) *url.URL {
	url, _ := url.Parse(fmt.Sprintf("https://%s?requestId=%s", approvalType, requestId))
	return url
}

var handler = NewAuthorizationHandler(
	MockClientLookupFn,
	//MockCodeGenerator,
	MockErrorPageRenderer,
	mockApprovalRequestStore,
	mockApprovalLookup)

func TestAuthorizeHandler_successfulAuthorization(t *testing.T) {
	responseType := "code"
	state := "state1"
	redirectUri := "https://example.com/redirect"

	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s&redirect_uri=%s",
		validClientId, responseType, state, url.QueryEscape(redirectUri))

	handlers.DoGetEndpointRequest(handler, requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		response.AssertHttpStatusEquals(302)
		response.AssertHasHeader("location", t)

		uri, err := url.Parse(response.GetHeader("location"))
		util.AssertNoError(err, t)

		withoutQuery := fmt.Sprintf("%s://%s%s", uri.Scheme, uri.Host, uri.Path)
		util.AssertEquals("https://username-password", withoutQuery, t)

		params, err := url.ParseQuery(uri.RawQuery)
		util.AssertNoError(err, t)

		util.AssertEquals("random-request-id", params.Get("requestId"), t)
		util.AssertEquals("", params.Get("error"), t)

		expectedRequest := oauth2.AuthorizeRequest{
			ResponseType: "code",
			ClientId:     validClientId,
			RedirectURI:  redirectUri,
			State:        state,
		}

		util.AssertEquals(1, len(approvalRequests), t)
		util.AssertEquals(expectedRequest, *approvalRequests["random-request-id"], t)

		return nil
	}, t)
}

func TestAuthorizeHandler_invalidClient(t *testing.T) {
	responseType := "code"
	state := "state1"
	redirectUri := "https://example.com/redirect"

	requestUrl := fmt.Sprintf("/authorize?client_id=%s&response_type=%s&state=%s&redirect_uri=%s",
		invalidClientId, responseType, state, url.QueryEscape(redirectUri))

	handlers.DoGetEndpointRequest(handler, requestUrl).
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

	handlers.DoGetEndpointRequest(handler, requestUrl).
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

	handlers.DoGetEndpointRequest(handler, requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		response.AssertHttpStatusEquals(302)
		response.AssertHasHeader("location", t)

		uri, err := url.Parse(response.GetHeader("location"))
		util.AssertNoError(err, t)

		withoutQuery := fmt.Sprintf("%s://%s%s", uri.Scheme, uri.Host, uri.Path)
		util.AssertEquals("https://username-password", withoutQuery, t)

		params, err := url.ParseQuery(uri.RawQuery)
		util.AssertNoError(err, t)

		util.AssertEquals("random-request-id", params.Get("requestId"), t)
		util.AssertEquals("", params.Get("error"), t)

		expectedRequest := oauth2.AuthorizeRequest{
			ResponseType: "code",
			ClientId:     validClientId,
			RedirectURI:  defaultRedirect,
			State:        state,
		}

		util.AssertEquals(1, len(approvalRequests), t)
		util.AssertEquals(expectedRequest, *approvalRequests["random-request-id"], t)
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

	handlers.DoGetEndpointRequest(handler, requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		response.AssertHttpStatusEquals(200)
		response.AssertHeaderValue("Content-type", "text/html", t)
		response.AssertResponseContent("<html>invalid_scope</html>", t)
		return nil
	}, t)
}

func TestAuthorizeHandler_invalidRequested(t *testing.T) {
	responseType := "code"
	state := "state1"
	redirectUri := "https://example.com/redirect"

	requestUrl := fmt.Sprintf("/authorize?response_type=%s&state=%s&redirect_uri=%s", responseType,
		state, url.QueryEscape(redirectUri))

	handlers.DoGetEndpointRequest(handler, requestUrl).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		response.AssertHttpStatusEquals(200)
		response.AssertHeaderValue("Content-type", "text/html", t)
		response.AssertResponseContent("<html>invalid_request</html>", t)
		return nil
	}, t)
}

func TestMain(m *testing.M) {
	approvalRequests = map[string]*oauth2.AuthorizeRequest{}
	result := m.Run()
	os.Exit(result)
}