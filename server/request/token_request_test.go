package request

import (
	"net/url"
	"testing"
	"reflect"
)

func TestTokenRequestWithMissingFields(t *testing.T) {
	tokenRequest, parseFailure := ParseTokenRequest(url.Values{})
	if tokenRequest != nil {
		t.Error("Expected parse error, not request")
	} else if !reflect.DeepEqual(parseFailure.MissingFields, []string{"grant_type"}) {
		t.Errorf("Missing fields were %s", parseFailure.MissingFields)
	}
}