package oauth

import (
	"testing"
	"net/http"
	"strings"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"reflect"
)

var validClientId = "valid-client-id"
var validClientSecret = "valid-client-secret"
var invalidClientSecret = "invalid-client-secret"

var unknownClientId = "unknown-client-id"

type MockClient struct {
}

func (MockClient) checkSecret(secret string) bool {
	return secret == validClientSecret
}

func mockClientLookup(clientId string) (Client, error) {
	if clientId == validClientId {
		return MockClient{}, nil
	} else {
		return nil, nil
	}
}

func mockTokenGenerator() string {
	return "mock-token"
}

var tokenHandlerConfig = TokenHandlerConfig{
	DefaultTokenExpiration: 3600,
}

func init() {
	http.HandleFunc("/token", NewTokenHandler(tokenHandlerConfig, mockClientLookup, mockTokenGenerator))
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
	json      map[string]interface{}
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

func doTokenEndpointRequestWithBody(grantType string, clientId string, clientSecret string) *TokenResponse {
	body := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s", grantType, clientId, clientSecret)
	post, err := http.NewRequest("POST", "http://localhost:8080/token",
		strings.NewReader(body))

	if post.Header.Set("Content-Type", "application/x-www-form-urlencoded"); err != nil {
		return &TokenResponse{
			err: err,
		}
	}

	//TODO: header case
	//creds := fmt.Sprintf("%s:%s", validClientId, validClientSecret)
	//post.Header.Set("Authorization", base64.StdEncoding.EncodeToString([]byte(creds)))

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
	err = json.Unmarshal(responseBody, &jwt)
	if err != nil {
		return &TokenResponse{
			err: err,
		}
	}
	return &TokenResponse{
		json:      jwt,
		httpStatus: response.StatusCode,
		err:        nil,
	}
}

func TestTokenHandler_ClientCredentialsWithValidData(t *testing.T) {
	doTokenEndpointRequestWithBody("client_credentials", validClientId, validClientSecret).
		thenAssert(func(response *TokenResponse) error {
		if response.httpStatus != 200 {
			return fmt.Errorf("Expected 401, but got %d", response.httpStatus)
		}
		expected := map[string]interface{}{
			"access_token": "mock-token",
			"token_type":   "Bearer",
			"expires_in":   float64(3600),
		}

		if !reflect.DeepEqual(response.json, expected) {
			return fmt.Errorf("Returned jwt didn't match. \nExpected: %+v. \nWas:      %+v\n", expected,
				response.json)
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

//TODO: ensure that header and body methods together is an error
