package oauth

import (
	"testing"
	"net/http"
	"strings"
	"fmt"
)

var validClientId = "valid-client-id"
var validClientSecret = "value-client-secret"

func mockClientLookup (clientId string, clientSecret string) (*Client, error) {
	if clientId == validClientId && clientSecret == validClientSecret {
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
	_, err := http.Post("http://localhost:8080/token", "application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s", validClientId, validClientSecret)))

	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

}
