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

var tokenHandlerConfig = TokenHandlerConfig {
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

func TestTokenHandler_ClientCredentialsWithValidData(t *testing.T) {
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", validClientId, validClientSecret)
	post, err := http.NewRequest("POST", "http://localhost:8080/token",
		strings.NewReader(body))
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//creds := fmt.Sprintf("%s:%s", validClientId, validClientSecret)
	//post.Header.Set("Authorization", base64.StdEncoding.EncodeToString([]byte(creds)))

	response, err := http.DefaultClient.Do(post)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	var jwt map[string]interface{}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil || response.StatusCode != 200 {
		t.Errorf("Unexpected error %+v, %d, %s", err, response.StatusCode, string(responseBody))
		return
	}

	err = json.Unmarshal(responseBody, &jwt)

	expected := map[string]interface{}{
		"access_token":"mock-token",
		"token_type": "Bearer",
		"expires_in": float64(3600),
	}

	if err != nil || !reflect.DeepEqual(jwt, expected) {
		t.Errorf("Returned jwt didn't match. \nExpected: %+v. \nWas:      %+v\nError:     %v", expected, jwt, err)
	}
}

func TestTokenHandler_UnknownClientInBody(t *testing.T) {
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", unknownClientId, validClientSecret)
	post, err := http.NewRequest("POST", "http://localhost:8080/token",
		strings.NewReader(body))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(post)

	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	if response.StatusCode != 401 {
		t.Errorf("Expected 401, but got %d", response.StatusCode)
		return
	}

	var errorMessage interface{}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}
	json.Unmarshal(data, &errorMessage)

	expected := map[string]interface{}{
		"error":             "invalid_client",
		"error_description": "Invalid client.",
		"error_uri":         "",
	}

	if !reflect.DeepEqual(errorMessage, expected) {
		t.Errorf("Got unexpected error %+v", errorMessage)
		return
	}
}

func TestTokenHandler_IncorrectSecret(t *testing.T) {
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", validClientId, invalidClientSecret)
	post, err := http.NewRequest("POST", "http://localhost:8080/token",
		strings.NewReader(body))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(post)

	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	if response.StatusCode != 401 {
		t.Errorf("Expected 401, but got %d", response.StatusCode)
		return
	}

	var errorMessage interface{}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}
	json.Unmarshal(data, &errorMessage)

	expected := map[string]interface{}{
		"error":             "invalid_client",
		"error_description": "Invalid client.",
		"error_uri":         "",
	}

	if !reflect.DeepEqual(errorMessage, expected) {
		t.Errorf("Got unexpected error %+v", errorMessage)
		return
	}
}

func TestTokenHandler_UnknownGrantType(t *testing.T) {
	body := fmt.Sprintf("grant_type=non_existing&client_id=%s&client_secret=%s", validClientId, validClientSecret)
	post, err := http.NewRequest("POST", "http://localhost:8080/token",
		strings.NewReader(body))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(post)

	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	if response.StatusCode != 400 {
		t.Errorf("Expected 400, but got %d", response.StatusCode)
		return
	}

	var error map[string]string
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}
	json.Unmarshal(data, &error)

	if error["error"] != "unsupported_grant_type" {
		t.Errorf("Got unexpected error %s but should have been 'unsupported_grant_type'", error["error"])
		return
	}
}

//TODO: ensure that header and body methods together is an error
