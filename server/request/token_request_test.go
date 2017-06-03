package request

import (
	"net/url"
	"testing"
	"net/http"
	"encoding/base64"
)

func TestTokenRequestWithMissingFields(t *testing.T) {
	req := http.Request{
		Form: url.Values{},
	}
	tokenRequest, err := ParseTokenRequest(req)
	if tokenRequest != nil {
		t.Error("Expected parse error, not request")
	} else if err.Error() != "The following fields are required: [grant_type]" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestTokenRequestWithAllClientCredentialsFields(t *testing.T) {
	req := http.Request{
		Form: url.Values{
			"grant_type":    []string{"client_credentials"},
			"scope":         []string{"foo bar"},
			"client_id":     []string{"the-client-id"},
			"client_secret": []string{"the-client-secret"},
		},
	}
	tokenRequest, err := ParseTokenRequest(req)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
		return
	}

	expected := TokenRequest{
		GrantType:    GRANT_TYPE_CLIENT_CREDENTIALS,
		Scope:        "foo bar",
		ClientId:     "the-client-id",
		ClientSecret: "the-client-secret",
	}

	if *tokenRequest != expected {
		t.Errorf("Unmatching token request: %+v", tokenRequest)
		return
	}
}

func TestTokenRequestWithAllClientCredentialsInHeaders(t *testing.T) {
	token := base64.StdEncoding.EncodeToString([]byte("the-client-id:the-client-secret"))
	req := http.Request{
		Form: url.Values{
			"grant_type": []string{"client_credentials"},
			"scope":      []string{"foo bar"},
		},
		Header: map[string][]string{
			"Authorization": {"Basic " + token},
		},
	}
	tokenRequest, err := ParseTokenRequest(req)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
		return
	}

	expected := TokenRequest{
		GrantType:    GRANT_TYPE_CLIENT_CREDENTIALS,
		Scope:        "foo bar",
		ClientId:     "the-client-id",
		ClientSecret: "the-client-secret",
	}

	if *tokenRequest != expected {
		t.Errorf("Unmatching token request: %+v", tokenRequest)
		return
	}
}

func TestTokenRequestWithBearerInHeadersAndBodyFails(t *testing.T) {
	token := base64.StdEncoding.EncodeToString([]byte("the-client-id:the-client-secret"))
	req := http.Request{
		Form: url.Values{
			"grant_type":    []string{"client_credentials"},
			"scope":         []string{"foo bar"},
			"client_id":     []string{"the-client-id"},
			"client_secret": []string{"the-client-secret"},
		},
		Header: map[string][]string{
			"Authorization": {"Bearer " + token},
		},
	}
	_, err := ParseTokenRequest(req)
	if err != ERR_INVALID_CLIENT {
		t.Errorf("Expected %s but got %+v", ERR_INVALID_CLIENT, err)
		return
	}

}

//TODO: check if form isn't parseable

//TODO: support for flows other than client credentials

//TODO: ensure that header and body methods together is an error

//TODO: ensure that header with multiple : fails
