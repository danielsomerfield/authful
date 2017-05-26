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
	} else if !reflect.DeepEqual(parseFailure.MissingFields, []string{"client_id", "response_type"}) {
		t.Errorf("Missing fields were %s", parseFailure.MissingFields)
	}
}

func TestAuthorizeRequestWithOnlyRequiredFields(t *testing.T) {
	authorizationRequest, parseFailure := ParseAuthorizeRequest(url.Values{
		"response_type": []string{"code"},
		"client_id":    []string{"cid"},
	})
	if parseFailure != nil {
		t.Error("No error expected")
	} else if *authorizationRequest != (AuthorizeRequest{
		ResponseType: "code",
		ClientId:     "cid",
	}) {
		t.Errorf("Authorization request looks like this: %+v", authorizationRequest)
	}
}

func TestAuthorizeRequestWithOptionalFields(t *testing.T) {
	authorizationRequest, parseFailure := ParseAuthorizeRequest(url.Values{
		"response_type": []string{"code"},
		"client_id":    []string{"cid"},
		"redirect_uri":    []string{"http://blah"},
		"scope":    []string{"the scope"},
		"state":    []string{"the state"},
	})
	if parseFailure != nil {
		t.Error("No error expected")
	} else if *authorizationRequest != (AuthorizeRequest{
		ResponseType: "code",
		ClientId:     "cid",
		RedirectURI:  "http://blah",
		Scope:        "the scope",
		State:        "the state",
	}) {
		t.Errorf("Authorization request looks like this: %+v", authorizationRequest)
	}
}
