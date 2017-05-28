package main

import (
	"testing"
	"github.com/danielsomerfield/authful/testutils"
	"io/ioutil"
	"net/http"
	"github.com/danielsomerfield/authful/server"
	"strings"
	"fmt"
	"encoding/json"
)

type TokenResponse struct {
	AccessToken  string        `json:"access_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    int32         `json:"expires_in"`
	RefreshToken string        `json:"refresh_token"`
}

func requestAdminToken(credentials server.Credentials) (*TokenResponse, error) {
	var err error = nil

	var request *http.Request
	if request, err = http.NewRequest("POST", "http://localhost:8080/token",
		strings.NewReader("grant_type=client_credentials")); err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Basic "+credentials.String())

	var response *http.Response
	var body []byte

	response, err = http.DefaultClient.Do(request)

	body, err = ioutil.ReadAll(response.Body)

	tokenResponse := TokenResponse{}
	err = json.Unmarshal(body, &tokenResponse)
	if err == nil {
		return &tokenResponse, nil
	} else {
		return nil, err
	}
}

func TestAuthorize(t *testing.T) {

	var err error = nil
	var credentials *server.Credentials
	var authServer *server.AuthServer

	authServer, credentials, err = testutils.RunServer()
	defer authServer.Stop()

	token, err := requestAdminToken(*credentials)
	fmt.Printf("Token: %+v, %+v", token, err)
	return

	//Register a client and get back client id and secret

	//Register a user

	//Hit the authorization endpoint
	//resp, err = http.Get("http://localhost:8080/authorize?request_type=code?client_id=1234")
	//
	////Ensure that the user is authentcated and prompted for approval
	//if err == nil {
	//	if resp.StatusCode == 200 {
	//		body, err = ioutil.ReadAll(resp.Body)
	//		if err == nil {
	//			print(string(body))
	//		}
	//	} else {
	//		t.Errorf("Expected status code 200 but was %s", resp.StatusCode)
	//	}
	//} else {
	//
	//}
}
