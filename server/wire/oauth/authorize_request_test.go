package oauth

import (
	"testing"
	"net/url"
	"net/http"
)

func TestAuthorizeRequestWithMissingFields(t *testing.T) {
	authorizationRequest, err := ParseAuthorizeRequest(http.Request{
		Form: url.Values{},
	})
	if authorizationRequest != nil {
		t.Error("Expected parse error, not request")
	} else if err.Error() != "The following fields are required: [client_id response_type]" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestAuthorizeRequestWithOnlyRequiredFields(t *testing.T) {

	authorizationRequest, parseFailure := ParseAuthorizeRequest(http.Request{
		Form: url.Values{
			"response_type": []string{"code"},
			"client_id":     []string{"cid"},
		},
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
	authorizationRequest, parseFailure := ParseAuthorizeRequest(http.Request{
		Form:url.Values{
			"response_type": []string{"code"},
			"client_id":     []string{"cid"},
			"redirect_uri":  []string{"http://blah"},
			"scope":         []string{"the scope"},
			"state":         []string{"the state"},
		},
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
