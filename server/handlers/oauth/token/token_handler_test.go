package token

import (
	"testing"
	"net/http"
	"strings"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"time"
	"bytes"
	"errors"
	"github.com/danielsomerfield/authful/server/service/oauth"
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

func init() {
	http.HandleFunc("/token", NewTokenHandler(tokenHandlerConfig, LookupClientFn, mockTokenGenerator, mockTokenStore.StoreToken, mockCurrentTimeFn))
	go http.ListenAndServe(":8080", nil)
}

func TestTokenHandler_RejectsGetRequest(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/token")
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}
}

type TokenResponse struct {
	json       map[string]interface{}
	httpStatus int
	err        error
}

func (rs *TokenResponse) thenAssert(test func(response *TokenResponse) error, t *testing.T) error {
	if rs.err != nil {
		t.Errorf("Request failed: %+v", rs.err)
		return rs.err
	}

	err := test(rs)
	if err != nil {
		t.Errorf("Assertion failed: %+v", err)
		return rs.err
	}
	return nil
}

func doTokenEndpointRequestWithBodyAndScope(grantType string, clientId string, clientSecret string, scope string) *TokenResponse {

	body := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s", grantType, clientId, clientSecret)
	if scope != "" {
		body = body + "&scope=" + scope
	}
	post, _ := http.NewRequest("POST", "http://localhost:8080/token",
		strings.NewReader(body))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(post)
	if err != nil {
		return &TokenResponse{
			err: err,
		}
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &TokenResponse{
			err: err,
		}
	}

	var jwt map[string]interface{}
	decoder := json.NewDecoder(bytes.NewBuffer(responseBody))
	decoder.UseNumber()
	decoder.Decode(&jwt)
	if err != nil {
		return &TokenResponse{
			err: err,
		}
	}
	return &TokenResponse{
		json:       jwt,
		httpStatus: response.StatusCode,
		err:        nil,
	}

}

func doTokenEndpointRequestWithBody(grantType string, clientId string, clientSecret string) *TokenResponse {
	return doTokenEndpointRequestWithBodyAndScope(grantType, clientId, clientSecret, "")
}

func TestTokenHandler_ClientCredentialsWithValidData(t *testing.T) {
	doTokenEndpointRequestWithBodyAndScope("client_credentials", validClientId, validClientSecret, "scope1 scope2").
		thenAssert(func(response *TokenResponse) error {
		if response.httpStatus != 200 {
			return fmt.Errorf("Expected 200, but got %d", response.httpStatus)
		}
		expected := map[string]interface{}{
			"access_token": "mock-token",
			"token_type":   "Bearer",
			"expires_in":   json.Number("3600"),
			"scope":        "scope1 scope2",
		}

		if !reflect.DeepEqual(response.json, expected) {
			return fmt.Errorf("Returned jwt didn't match. \nExpected: %+v. \nWas:      %+v\n", expected,
				response.json)
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
		thenAssert(func(response *TokenResponse) error {
		if response.httpStatus != 401 {
			return fmt.Errorf("Expected 401, but got %d", response.httpStatus)
		}

		expected := map[string]interface{}{
			"error":             "invalid_client",
			"error_description": "Invalid client.",
			"error_uri":         "",
		}

		if !reflect.DeepEqual(response.json, expected) {
			return fmt.Errorf("Returned data didn't match. \nExpected: %+v. \nWas:      %+v\n", expected,
				response.json)
		}

		return nil
	}, t)
}

func TestTokenHandler_IncorrectSecret(t *testing.T) {
	doTokenEndpointRequestWithBody("client_credentials", validClientId, invalidClientSecret).
		thenAssert(func(response *TokenResponse) error {
		if response.httpStatus != 401 {
			return fmt.Errorf("Expected 401, but got %d", response.httpStatus)
		}

		expected := map[string]interface{}{
			"error":             "invalid_client",
			"error_description": "Invalid client.",
			"error_uri":         "",
		}

		if !reflect.DeepEqual(response.json, expected) {
			return fmt.Errorf("Returned data didn't match. \nExpected: %+v. \nWas:      %+v\n", expected,
				response.json)
		}

		return nil
	}, t)
}

func TestTokenHandler_UnknownGrantType(t *testing.T) {
	doTokenEndpointRequestWithBody("not real", validClientId, validClientSecret).
		thenAssert(func(response *TokenResponse) error {
		if response.httpStatus != 400 {
			return fmt.Errorf("Expected 400, but got %d", response.httpStatus)
		}

		if response.json["error"] != "unsupported_grant_type" {
			return fmt.Errorf("Got unexpected error %s but should have been 'unsupported_grant_type'",
				response.json["error"])
		}
		return nil
	}, t)
}

func TestTokenHandler_ClientCredentialsWithInvalidScope(t *testing.T) {
	doTokenEndpointRequestWithBodyAndScope("client_credentials", validClientId, validClientSecret, "badscope").
		thenAssert(func(response *TokenResponse) error {
		if response.httpStatus != 400 {
			return fmt.Errorf("Expected 400, but got %d", response.httpStatus)
		}

		if response.json["error"] != "invalid_scope" {
			return fmt.Errorf("Got unexpected error %s but should have been 'unsupported_grant_type'",
				response.json["error"])
		}
		return nil
	}, t)
}
