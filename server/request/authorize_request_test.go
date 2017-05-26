package request

import (
	"testing"
	"net/url"
	"reflect"
)

func TestAuthorizeRequestWithMissingFields(t *testing.T) {
	authorizationRequest, parseFailure := ParseAuthorizeRequest(url.Values{})
	if authorizationRequest != nil {
		t.Error("Expected parse error, not request")
	} else if !reflect.DeepEqual(parseFailure.MissingFields, []string{"client_id", "request_type"}) {
		t.Errorf("Missing fields were %s", parseFailure.MissingFields)
	}
}

func TestAuthorizeRequestWithOnlyRequiredFields(t *testing.T) {
	authorizationRequest, parseFailure := ParseAuthorizeRequest(url.Values{
		"request_type": []string{"code"},
		"client_id":    []string{"cid"},
	})
	if parseFailure != nil {
		t.Error("No error expected")
	} else if *authorizationRequest != (AuthorizeRequest{
		RequestType: "code",
		ClientId:    "cid",
	}) {
		t.Errorf("Authorization request looks like this: %+v", authorizationRequest)
	}
}

func TestAuthorizeRequestWithOptionalFields(t *testing.T) {
	authorizationRequest, parseFailure := ParseAuthorizeRequest(url.Values{
		"request_type": []string{"code"},
		"client_id":    []string{"cid"},
		"redirect_uri":    []string{"http://blah"},
	})
	if parseFailure != nil {
		t.Error("No error expected")
	} else if *authorizationRequest != (AuthorizeRequest{
		RequestType: "code",
		ClientId:    "cid",
		RedirectURI:    "http://blah",
	}) {
		t.Errorf("Authorization request looks like this: %+v", authorizationRequest)
	}
}
