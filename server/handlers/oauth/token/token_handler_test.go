package token

import (
	"testing"
	"net/http"
	"fmt"
	"encoding/json"
	"reflect"
	"time"
	"errors"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"net/http/httptest"
	"github.com/danielsomerfield/authful/server/handlers"
)

var validClientId = "valid-client-id"
var validClientSecret = "valid-client-secret"
var invalidClientSecret = "invalid-client-secret"
var unknownClientId = "unknown-client-id"

type MockClient struct {
	scope []string
}

func (MockClient) CheckSecret(secret string) bool {
	return secret == validClientSecret
}

func (mc MockClient) GetScopes() []string {
	return mc.scope
}

func mockTokenGenerator() string {
	return "mock-token"
}

var tokenHandlerConfig = TokenHandlerConfig{
	DefaultTokenExpiration: 3600,
}

type MockTokenStore struct {
	storedTokens map[string]oauth.TokenMetaData
}

var mockTokenStore = MockTokenStore{
	storedTokens: map[string]oauth.TokenMetaData{},
}

func (m MockTokenStore) StoreToken(token string, clientMetaData oauth.TokenMetaData) error {
	m.storedTokens[token] = clientMetaData
	return nil
}

func LookupClientFn(clientId string) (oauth.Client, error) {
	if clientId == validClientId {
		return MockClient{
			scope: []string{"scope1", "scope2"},
		}, nil
	} else {
		return nil, nil
	}
}

var mockNow = time.Now()

func mockCurrentTimeFn() time.Time {
	return mockNow
}

func TestTokenHandler_RejectsGetRequest(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost:8080/token", nil)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(NewTokenHandler(tokenHandlerConfig, LookupClientFn, mockTokenGenerator, mockTokenStore.StoreToken, mockCurrentTimeFn))
	handler.ServeHTTP(response, request)

	if response.Code != 400 {
		t.Errorf("Expected 400 but got %d", response.Code)
	}
}

func doTokenEndpointRequestWithBodyAndScope(grantType string, clientId string, clientSecret string, scope string) *handlers.EndpointResponse {
	body := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s", grantType, clientId, clientSecret)
	if scope != "" {
		body = body + "&scope=" + scope
	}
	return handlers.DoEndpointRequest(
		NewTokenHandler(tokenHandlerConfig, LookupClientFn, mockTokenGenerator, mockTokenStore.StoreToken, mockCurrentTimeFn),
		"http://localhost:8080/token",
		body)
}

func doTokenEndpointRequestWithBody(grantType string, clientId string, clientSecret string) *handlers.EndpointResponse {
	return doTokenEndpointRequestWithBodyAndScope(grantType, clientId, clientSecret, "")
}

func TestTokenHandler_ClientCredentialsWithValidData(t *testing.T) {
	doTokenEndpointRequestWithBodyAndScope("client_credentials", validClientId, validClientSecret, "scope1 scope2").
		ThenAssert(func(response *handlers.EndpointResponse) error {
		if response.HttpStatus != 200 {
			return fmt.Errorf("Expected 200, but got %d", response.HttpStatus)
		}
		expected := map[string]interface{}{
			"access_token": "mock-token",
			"token_type":   "Bearer",
			"expires_in":   json.Number("3600"),
			"scope":        "scope1 scope2",
		}

		if !reflect.DeepEqual(response.Json, expected) {
			return fmt.Errorf("Returned jwt didn't match. \nExpected: %+v. \nWas:      %+v\n", expected,
				response.Json)
		}

		if tokenMetaData, ok := mockTokenStore.storedTokens["mock-token"]; ok {
			expectedTokenMetaData := oauth.TokenMetaData{
				Token:      "mock-token",
				Expiration: mockNow.Add(time.Duration(tokenHandlerConfig.DefaultTokenExpiration) * time.Second),
				ClientId:   validClientId,
			}
			if tokenMetaData != expectedTokenMetaData {
				return fmt.Errorf("Token meta data didn't match. \nExpected: %+v. \nWas:      %+v\n",
					expectedTokenMetaData,
					tokenMetaData)
			}
		} else {
			return errors.New("The token did not get stored.")
		}

		return nil
	}, t)
}

func TestTokenHandler_UnknownClientInBody(t *testing.T) {
	doTokenEndpointRequestWithBody("client_credentials", unknownClientId, validClientSecret).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		if response.HttpStatus != 401 {
			return fmt.Errorf("Expected 401, but got %d", response.HttpStatus)
		}

		expected := map[string]interface{}{
			"error":             "invalid_client",
			"error_description": "Invalid client.",
			"error_uri":         "",
		}

		if !reflect.DeepEqual(response.Json, expected) {
			return fmt.Errorf("Returned data didn't match. \nExpected: %+v. \nWas:      %+v\n", expected,
				response.Json)
		}

		return nil
	}, t)
}

func TestTokenHandler_IncorrectSecret(t *testing.T) {
	doTokenEndpointRequestWithBody("client_credentials", validClientId, invalidClientSecret).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		if response.HttpStatus != 401 {
			return fmt.Errorf("Expected 401, but got %d", response.HttpStatus)
		}

		expected := map[string]interface{}{
			"error":             "invalid_client",
			"error_description": "Invalid client.",
			"error_uri":         "",
		}

		if !reflect.DeepEqual(response.Json, expected) {
			return fmt.Errorf("Returned data didn't match. \nExpected: %+v. \nWas:      %+v\n", expected,
				response.Json)
		}

		return nil
	}, t)
}

func TestTokenHandler_UnknownGrantType(t *testing.T) {
	doTokenEndpointRequestWithBody("not real", validClientId, validClientSecret).
		ThenAssert(func(response *handlers.EndpointResponse) error {
		if response.HttpStatus != 400 {
			return fmt.Errorf("Expected 400, but got %d", response.HttpStatus)
		}

		if response.Json["error"] != "unsupported_grant_type" {
			return fmt.Errorf("Got unexpected error %s but should have been 'unsupported_grant_type'",
				response.Json["error"])
		}
		return nil
	}, t)
}

func TestTokenHandler_ClientCredentialsWithInvalidScope(t *testing.T) {
	doTokenEndpointRequestWithBodyAndScope("client_credentials", validClientId, validClientSecret, "badscope").
		ThenAssert(func(response *handlers.EndpointResponse) error {
		if response.HttpStatus != 400 {
			return fmt.Errorf("Expected 400, but got %d", response.HttpStatus)
		}

		if response.Json["error"] != "invalid_scope" {
			return fmt.Errorf("Got unexpected error %s but should have been 'unsupported_grant_type'",
				response.Json["error"])
		}
		return nil
	}, t)
}
