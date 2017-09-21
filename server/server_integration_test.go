package server

import (
	"testing"
	"io/ioutil"
	"net/http"
	"strings"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/clientcredentials"
	"context"
	oauth_service "github.com/danielsomerfield/authful/server/service/oauth"
	oauth_wire "github.com/danielsomerfield/authful/server/wire/oauth"
	"os"
	"log"
	"github.com/danielsomerfield/authful/util"
	"golang.org/x/oauth2"
	"github.com/danielsomerfield/authful/client/admin"
	"net/url"
)

func TestAuthorize(t *testing.T) {
	scope := "scope1"
	redirectURI := "https://localhost:8181/ping"

	_, err := requestAdminToken(*creds)
	util.AssertNoError(err, t)

	//Register a client and get back client id and secret
	apiClient := createAPIClient(t)
	registration, err := apiClient.RegisterClient(
		"test-authorize",
		[]string{scope},
		[]string{redirectURI},
		"")
	util.AssertNoError(err, t)
	util.AssertNotNil(registration, t)

	//Register a user
	_, err = apiClient.RegisterUser("username-1", "password-1", []string{"username-password"})
	util.AssertNoError(err, t)

	state := util.GenerateRandomString(5)

	//Hit the authorization endpoint
	httpsClient := CreateHttpsClient()
	authorizeUrl := fmt.Sprintf(
		"https://localhost:8081/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		registration.Data.ClientId, url.QueryEscape(redirectURI), scope, state)
	resp, err := httpsClient.Get(authorizeUrl)
	util.AssertNoError(err, t)

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		print(string(body))
	}

	//TODO: make sure this is actually the login endpoint
	util.AssertStatusCode(resp, 200, t)

}

func TestClientCredentialsEnd2End(t *testing.T) {

	config := clientcredentials.Config{
		ClientID:     creds.ClientId,
		ClientSecret: creds.ClientSecret,
		TokenURL:     "https://localhost:8081/token",
		Scopes:       []string{},
	}

	httpsClient := CreateHttpsClient()
	httpC := config.Client(context.WithValue(context.Background(), oauth2.HTTPClient, httpsClient))

	resp, err := httpC.Get("https://localhost:8181/test")
	util.AssertNoError(err, t)
	util.AssertStatusCode(resp, 200, t)

	_, err = ioutil.ReadAll(resp.Body)
	util.AssertNoError(err, t)

	//No token
	httpsClient = CreateHttpsClient()
	resp, err = httpsClient.Get("https://localhost:8181/test")
	util.AssertNoError(err, t)
	util.AssertStatusCode(resp, 401, t)

	//Invalid token
	invalidRequest, _ := http.NewRequest("GET", "https://localhost:8181/test", nil)
	invalidRequest.Header.Set("Authorization", "Bearer 12345")
	resp, err = CreateHttpsClient().Do(invalidRequest)

	util.AssertNoError(err, t)
	util.AssertStatusCode(resp, 401, t)
}

func TestScopeRequirements(t *testing.T) {

	apiClient := createAPIClient(t)

	//Register a client with no scopes using admin credentials
	registration, err := apiClient.RegisterClient(
		"test-client",
		[]string{},
		[]string{},
		"")
	util.AssertNoError(err, t)
	util.AssertNotNil(registration, t)

	//TODO: validate the client was registered

	apiClientInsufficientCreds, err := admin.NewAPIClient("https://localhost:8081",
		registration.Data.ClientId,
		registration.Data.ClientSecret,
		[]string{"../resources/test/server.crt"},
	)
	util.AssertNoError(err, t)

	_, err = apiClientInsufficientCreds.RegisterClient(
		"test-client-2",
		[]string{},
		[]string{},
		"")
	util.AssertNotNil(err, t)
	util.AssertTrue(err.(admin.ClientError).ErrorType() == "invalid_client", "Expected to equal \"invalid_client\"", t)
}

func TestErrorResponse(t *testing.T) {

	response, err := CreateHttpsClient().Post("https://localhost:8081/token", "application/json", strings.NewReader(""))
	util.AssertNoError(err, t)

	util.AssertStatusCode(response, 400, t)

	body, err := ioutil.ReadAll(response.Body)
	util.AssertNoError(err, t)

	errorResponse := oauth_wire.ErrorResponse{} //TODO: remove dependency on "production" type

	if err := json.Unmarshal(body, &errorResponse); err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	expected := oauth_wire.ErrorResponse{
		Error:            "invalid_request",
		ErrorDescription: "The following fields are required: [grant_type]",
		ErrorURI:         "",
	}

	if errorResponse != expected {
		t.Errorf("Unexpected response %+v", errorResponse)
	}
}

func createAPIClient(t *testing.T) (*admin.APIClient) {
	apiClient, err := admin.NewAPIClient("https://localhost:8081",
		creds.ClientId,
		creds.ClientSecret,
		[]string{TEST_CERTIFICATE},
	)
	util.AssertNoError(err, t)
	return apiClient
}

func requestAdminToken(credentials oauth_service.Credentials) (*oauth_wire.TokenResponse, error) {
	var err error = nil

	var request *http.Request
	if request, err = http.NewRequest("POST", "https://localhost:8081/token",
		strings.NewReader("grant_type=client_credentials")); err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Basic "+credentials.String())

	var response *http.Response
	var body []byte

	httpsClient := CreateHttpsClient()

	response, err = httpsClient.Do(request)

	body, err = ioutil.ReadAll(response.Body)

	tokenResponse := oauth_wire.TokenResponse{}
	err = json.Unmarshal(body, &tokenResponse)
	if err == nil {
		return &tokenResponse, nil
	} else {
		return nil, err
	}
}

func TestMain(m *testing.M) {
	var authServer *AuthServer
	var err error

	authServer, creds, err = RunServer()
	go StartResourceServer()

	if err != nil {
		log.Fatalf("Unexpected error on startup: %+v", err)
		return
	}
	result := m.Run()
	err = authServer.Stop()
	StopResourceServer()

	if err != nil {
		log.Fatalf("Unexpected error on stop: %+v", err)
		return
	}

	os.Exit(result)
}
