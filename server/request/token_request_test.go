package request

import (
	"net/url"
	"testing"
	"net/http"
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

//TODO: check valid parse from post body
//TODO: check valid parse from headers
//TODO: check if form isn't parseable

//TODO: support for flows other than client credentials
