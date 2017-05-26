package request

import (
	"testing"
	"net/url"
	"reflect"
)

func TestAuthorizeRequestWithMissingFields(t *testing.T)  {
	authorizationRequest, parseFailure := ParseAuthorizeRequest(url.Values{})
	if authorizationRequest != nil {
		t.Error("Expected parse error, not request")
	} else if !reflect.DeepEqual(parseFailure.missingFields, []string{"request_type", "client_id",}){
		t.Errorf("Missing fields were %s", parseFailure.missingFields)
	}
}


func TestAuthorizeRequestWithRequiredFields(t *testing.T)  {
	authorizationRequest, parseFailure := ParseAuthorizeRequest(url.Values{
		"request_type":[]string{"code"},
		"client_id":[]string{"cid"},
	})
	if parseFailure != nil {
		t.Error("No error expected")
	} else if *authorizationRequest != (AuthorizeRequest{
		requestType:"code",
		clientId:"cid",
	}) {
		t.Errorf("Authorization request looks like this: %+v", authorizationRequest)
	}

}