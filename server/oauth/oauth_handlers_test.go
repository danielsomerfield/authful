package oauth

import (
	"testing"
	"net/http"
	"strings"
	"fmt"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"reflect"
)

var validClientId = "valid-client-id"
var validClientSecret = "valid-client-secret"

var unknownClientId = "unknown-client-id"

func mockClientLookup(clientId string) (*Client, error) {
	if clientId == validClientId {
		return &Client{}, nil
	} else {
		return nil, nil
	}
}

func init() {
	http.HandleFunc("/token", NewTokenHandler(mockClientLookup))
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

func TestTokenHandler_ClientCredentials(t *testing.T) {
	post, err := http.NewRequest("POST", "http://localhost:8080/token", strings.NewReader("grant_type=client_credentials"))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	token := fmt.Sprintf("%s:%s", validClientId, validClientSecret)
	post.Header.Set("Authorization", base64.StdEncoding.EncodeToString([]byte(token)))

	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	//TODO check that token is returned

}

func TestTokenHandler_UnknownClientInBody(t *testing.T) {
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", unknownClientId, validClientSecret)
	post, err := http.NewRequest("POST", "http://localhost:8080/token",
		strings.NewReader(body))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//token := fmt.Sprintf("%s:%s", validClientId, validClientSecret)
	//post.Header.Set("Authorization", base64.StdEncoding.EncodeToString([]byte(token)))

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
		"error_description": "The client was not known.",
		"error_uri":         "",
	}

	if !reflect.DeepEqual(errorMessage, expected) {
		t.Errorf("Got unexpected error %+v", errorMessage)
		return
	}
}

//func TestTokenHandler_IncorrectSecret(t *testing.T) {
//
//}

//TODO: ensure that header and body methods together is an error
