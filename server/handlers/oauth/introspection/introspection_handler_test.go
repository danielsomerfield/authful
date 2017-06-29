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
)

func mockRequestValidation(request http.Request) bool {
	return request.Header.Get("Authorization") == "Bearer "+validBearerToken
}

func mockGetTokenMetaDataFn(token string) *oauth.TokenMetaData {
	if token == activeToken {
		return &oauth.TokenMetaData{
			Token:      token,
			Expiration: time.Time{},
			ClientId:   "",
		}
	}
	return nil
}

var activeToken = "active_token"
var unknownToken = "unknownToken"
var validBearerToken = "valid-bearer-token"

func TestIntrospectionHandler_ValidToken(t *testing.T) {
	post, _ := http.NewRequest("POST", "http://localhost:8080/introspect",
		strings.NewReader(fmt.Sprintf("token=%s", activeToken)))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.Header.Set("Authorization", "Bearer "+validBearerToken)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(NewIntrospectionHandler(mockRequestValidation, mockGetTokenMetaDataFn))
	handler.ServeHTTP(response, post)

	if response.Code != 200 {
		t.Errorf("Unexpected 200 but got %d", response.Code)
		return
	}

	responseJSON := map[string]interface{}{}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
		return
	}

	json.Unmarshal(responseBody, &responseJSON)

	expected := map[string]interface{}{
		"active": true,
	}
	fmt.Printf("======> %+v", responseJSON)

	if !reflect.DeepEqual(responseJSON, expected) {
		t.Errorf("Returned jwt didn't match. \nExpected: %+v. \nWas:      %+v\n", expected,
			responseJSON)
	}
}

func TestIntrospectionHandler_UnknownToken(t *testing.T) {
	post, _ := http.NewRequest("POST", "http://localhost:8080/introspect",
		strings.NewReader(fmt.Sprintf("token=%s", unknownToken)))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.Header.Set("Authorization", "Bearer "+validBearerToken)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(NewIntrospectionHandler(mockRequestValidation, mockGetTokenMetaDataFn))
	handler.ServeHTTP(response, post)

	if response.Code != 200 {
		t.Errorf("Unexpected 200 but got %d", response.Code)
		return
	}

	responseJSON := map[string]interface{}{}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
		return
	}

	json.Unmarshal(responseBody, &responseJSON)

	if responseJSON["active"] != false {
		t.Errorf("Expected active to equal 'false' but it was %s", responseJSON["active"])
	}
}
/*
//TODO: test with expired token
//TODO: test with invalid bearer:
WWW-Authenticate: Bearer realm="example",
                       error="invalid_token",
                       error_description="The access token expired"

 */
//TODO:	test with valid bearer inactive creds
