package accesscontrol

import (
	"net/http"
	"testing"
	"encoding/base64"
	"github.com/danielsomerfield/authful/server/service/oauth"
)

var validClientId = "valid-client-id"
var validClientSecret = "valid-client-secret"
var invalidClientSecret = "invalid-client-secret"

type MockClient struct {
}

func (MockClient) CheckSecret(secret string) bool {
	return secret == validClientSecret
}

func (mc MockClient) GetScopes() []string {
	return []string{}
}

func (mc MockClient) IsValidRedirectURI(uri string) bool {
	return true
}

func MockClientLookupFn(clientId string) (oauth.Client, error) {
	if clientId == validClientId {
		return MockClient{

		}, nil
	} else {
		return nil, nil
	}
}

var mockClientLookup = MockClientLookupFn

var clientAccessControlFn = NewClientAccessControlFn(mockClientLookup)

func TestIntrospect_FailsWithoutCredentials(t *testing.T) {
	request := http.Request{

	}

	passed, err := clientAccessControlFn(request)

	if err != nil {
		t.Errorf("Unexpected error %+v\n", err)
		return
	}

	if passed {
		t.Error("Succeeded with request that should have failed")
	}
}

func TestIntrospect_FailsWithInvalidSecret(t *testing.T) {
	request := http.Request{
		Header: http.Header{
			"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(validClientId+":"+invalidClientSecret))},
		},
	}

	passed, err := clientAccessControlFn(request)

	if err != nil {
		t.Errorf("Unexpected error %+v\n", err)
		return
	}

	if passed {
		t.Error("Succeeded with request that should have failed")
	}
}

func TestIntrospect_PassesWithValidCredentials(t *testing.T) {
	request := http.Request{
		Header: http.Header{
			"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(validClientId+":"+validClientSecret))},
		},
	}

	passed, err := clientAccessControlFn(request)

	if err != nil {
		t.Errorf("Unexpected error %+v\n", err)
		return
	}

	if !passed {
		t.Error("Failed with request that should have succeeded")
	}
}
