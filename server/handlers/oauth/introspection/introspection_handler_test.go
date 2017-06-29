package introspection

import (
	"net/http"
	"testing"
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"time"
	"net/http/httptest"
	"bytes"
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

	introspectWithToken(activeToken, validBearerToken).thenAssert(func(response *TokenResponse) error {
		if response.httpStatus != 200 {
			return fmt.Errorf("Expected 200, but got %d", response.httpStatus)
		}
		expected := map[string]interface{}{
			"active": true,
		}

		if !reflect.DeepEqual(response.json, expected) {
			t.Errorf("Returned jwt didn't match. \nExpected: %+v. \nWas:      %+v\n", expected,
				response.json)
		}
		return nil
	}, t)
}

func TestIntrospectionHandler_UnknownToken(t *testing.T) {
	introspectWithToken(unknownToken, validBearerToken).thenAssert(func(response *TokenResponse) error {
		if response.httpStatus != 200 {
			return fmt.Errorf("Expected 200, but got %d", response.httpStatus)
		}
		if response.json["active"] != false {
			return fmt.Errorf("Expected active to equal 'false' but it was %s", response.json["active"])
		}
		return nil
	}, t)
}

func TestIntrospectionHandler_ExpiredToken(t *testing.T) {
	introspectWithToken(expiredToken, validBearerToken).thenAssert(func(response *TokenResponse) error {
		if response.httpStatus != 200 {
			return fmt.Errorf("Expected 200, but got %d", response.httpStatus)
		}
		if response.json["active"] != false {
			return fmt.Errorf("Expected active to equal 'false' but it was %s", response.json["active"])
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

func introspectWithToken(tokenToValidate string, callingBearerToken string) * TokenResponse {

	body := fmt.Sprintf("token=%s", tokenToValidate)
	post, _ := http.NewRequest("POST", "http://localhost:8080/introspect",
		strings.NewReader(body))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.Header.Set("Authorization", "Bearer "+ callingBearerToken)


	response := httptest.NewRecorder()
	handler := http.HandlerFunc(NewIntrospectionHandler(mockRequestValidation, mockGetTokenMetaDataFn))
	handler.ServeHTTP(response, post)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &TokenResponse{
			err: err,
		}
	}

	var jwt map[string]interface{}
	decoder := json.NewDecoder(bytes.NewBuffer(responseBody))
	decoder.UseNumber()
	decoder.Decode(&jwt)
	if err != nil {
		return &TokenResponse{
			err: err,
		}
	}
	return &TokenResponse{
		json:       jwt,
		httpStatus: response.Code,
		err:        nil,
	}
}

type TokenResponse struct {
	json       map[string]interface{}
	httpStatus int
	err        error
}

func (rs *TokenResponse) thenAssert(test func(response *TokenResponse) error, t *testing.T) error {
	if rs.err != nil {
		t.Errorf("Request failed: %+v", rs.err)
		return rs.err
	}

	err := test(rs)
	if err != nil {
		t.Errorf("Assertion failed: %+v", err)
		return rs.err
	}
	return nil
}