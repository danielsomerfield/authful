package introspection

import (
	"net/http"
	"testing"
	"fmt"
	"reflect"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"time"
	"github.com/danielsomerfield/authful/server/handlers"
)

func mockRequestValidation(request http.Request) bool {
	return request.Header.Get("Authorization") == "Bearer "+validBearerToken
}

func mockGetTokenMetaDataFn(token string) *oauth.TokenMetaData {
	if token == activeToken {
		return &oauth.TokenMetaData{
			Token:      token,
			Expiration: time.Now().AddDate(1, 0, 0),
			ClientId:   "",
		}
	} else if token == expiredToken {
		return &oauth.TokenMetaData{
			Token:      token,
			Expiration: time.Now().AddDate(-1, 0, 0),
			ClientId:   "",
		}
	}
	return nil
}

var activeToken = "active-token"
var unknownToken = "unknown-token"
var expiredToken = "expired-token"
var validBearerToken = "valid-bearer-token"

func TestIntrospectionHandler_ValidToken(t *testing.T) {

	introspectWithToken(activeToken, validBearerToken).ThenAssert(func(response *handlers.EndpointResponse) error {
		if response.HttpStatus != 200 {
			return fmt.Errorf("Expected 200, but got %d", response.HttpStatus)
		}
		expected := map[string]interface{}{
			"active": true,
		}

		if !reflect.DeepEqual(response.Json, expected) {
			t.Errorf("Returned jwt didn't match. \nExpected: %+v. \nWas:      %+v\n", expected,
				response.Json)
		}
		return nil
	}, t)
}

func TestIntrospectionHandler_UnknownToken(t *testing.T) {
	introspectWithToken(unknownToken, validBearerToken).ThenAssert(func(response *handlers.EndpointResponse) error {
		if response.HttpStatus != 200 {
			return fmt.Errorf("Expected 200, but got %d", response.HttpStatus)
		}
		if response.Json["active"] != false {
			return fmt.Errorf("Expected active to equal 'false' but it was %s", response.Json["active"])
		}
		return nil
	}, t)
}

func TestIntrospectionHandler_ExpiredToken(t *testing.T) {
	introspectWithToken(expiredToken, validBearerToken).ThenAssert(func(response *handlers.EndpointResponse) error {
		if response.HttpStatus != 200 {
			return fmt.Errorf("Expected 200, but got %d", response.HttpStatus)
		}
		if response.Json["active"] != false {
			return fmt.Errorf("Expected active to equal 'false' but it was %s", response.Json["active"])
		}
		return nil
	}, t)
}

/*
//TODO: test with invalid bearer:
//TODO: check for WWW-Authenticate response on denials
WWW-Authenticate: Bearer realm="example",
                       error="invalid_token",
                       error_description="The access token expired"

 */
//TODO:	test with valid bearer inactive creds

func introspectWithToken(tokenToValidate string, callingBearerToken string) *handlers.EndpointResponse {
	body := fmt.Sprintf("token=%s", tokenToValidate)
	headers := map[string]string {
		"Authorization" : "Bearer "+callingBearerToken,
	}
	//post.Header.Set("Authorization", "Bearer "+callingBearerToken)
	return handlers.DoEndpointRequestWithHeaders(
		NewIntrospectionHandler(mockRequestValidation, mockGetTokenMetaDataFn),
		"http://localhost:8080/token",
		body, headers)
}
