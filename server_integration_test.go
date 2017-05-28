package main

import (
	"testing"
	"github.com/danielsomerfield/authful/testutils"
	"io/ioutil"
	"net/http"
	"github.com/danielsomerfield/authful/server"
	"strings"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/wireTypes"
	"fmt"
)

func requestAdminToken(credentials server.Credentials) (*wireTypes.TokenResponse, error) {
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

	tokenResponse := wireTypes.TokenResponse{}
	err = json.Unmarshal(body, &tokenResponse)
	if err == nil {
		return &tokenResponse, nil
	} else {
		return nil, err
	}
}

func TestErrorResponse(t *testing.T) {
	authServer, _, _ := testutils.RunServer()
	defer authServer.Stop()
	response, err := http.Post("http://localhost:8080/token", "application/json", strings.NewReader(""))
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	if response.StatusCode != 400 {
		t.Errorf("Expected 400 but got %s", response.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	errorResponse := wireTypes.ErrorResponse{}

	if err := json.Unmarshal(body, &errorResponse); err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	expected := wireTypes.ErrorResponse{
		Error :"invalid_request",
		ErrorDescription: "The following fields are required: [grant_type]",
		ErrorURI: "",
	}

	if errorResponse != expected {
		t.Errorf("Unexpected response %+v", errorResponse)
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
