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
)

func mockRequestValidation(request http.Request) bool {
	return request.Header.Get("Authorization") == "Bearer "+validBearerToken
}

func mockGetTokenMetaDataFn(token string) *oauth.TokenMetaData {
	return nil
}

func init() {
	http.HandleFunc("/introspect", NewIntrospectionHandler(mockRequestValidation))
	go http.ListenAndServe(":8080", nil)
}

var activeToken = "active_token"
var validBearerToken = "valid-bearer-token"

func TestIntrospectionHandler_ApprovesValidToken(t *testing.T) {
	return
	post, _ := http.NewRequest("POST", "http://localhost:8080/introspect",
		strings.NewReader(fmt.Sprintf("token=%s", activeToken)))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.Header.Set("Authorization", "Bearer "+validBearerToken)
	response, err := http.DefaultClient.Do(post)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
		return
	}

	if response.StatusCode != 200 {
		t.Errorf("Unexpected 200 but got %d", response.StatusCode)
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

/*
//TODO: test with invalid bearer:
WWW-Authenticate: Bearer realm="example",
                       error="invalid_token",
                       error_description="The access token expired"

 */
//TODO:	test with valid bearer inactive creds
