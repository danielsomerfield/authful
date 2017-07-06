package accesscontrol

import (
	"net/http"
	"testing"
	"github.com/danielsomerfield/authful/server/service/oauth"
)

var validClientId = "valid-client-id"
var validClientSecret = "valid-client-secret"

type MockClient struct {

}

func (MockClient) CheckSecret(secret string) bool {
	return secret == validClientSecret
}

func (mc MockClient) GetScopes() []string {
	return []string{}
}

func MockClientLookupFn(clientId string) (oauth.Client, error) {
	if clientId == validClientId {
		return MockClient{

		}, nil
	} else {
		return nil, nil
	}
}

var clientAccessControlFn = NewClientAccessControlFn(MockClientLookupFn)

func TestIntrospect_FailsWithoutCredentials(t *testing.T) {
	request := http.Request{

	}

	if clientAccessControlFn(request) {
		t.Error("Succeeded with request that should have failed")
	}
}

//func TestIntrospect_FailsWithInvalidCredentials(t *testing.T) {
//	request := http.Request{
//		Header: http.Header{
//			"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(validClientId+":"+validClientSecret))},
//		},
//	}
//
//	if clientAccessControlFn(request) {
//		t.Error("Failed with request that should have succeeded")
//	}
//}
