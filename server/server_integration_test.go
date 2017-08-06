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
	"time"
	"os"
	"log"
	"regexp"
	"github.com/danielsomerfield/authful/util"
	"crypto/tls"
	"crypto/x509"
	"golang.org/x/oauth2"
	"github.com/danielsomerfield/authful/client/admin"
	"net/url"
)

var creds *oauth_service.Credentials = nil

const TEST_CERTIFICATE = "../resources/test/server.crt"

func TestAuthorize(t *testing.T) {

	_, err := requestAdminToken(*creds)
	util.AssertNoError(err, t)

	//Register a client and get back client id and secret
	apiClient := createAPIClient(t)
	registration, err := apiClient.RegisterClient("test-authorize")
	util.AssertNoError(err, t)
	util.AssertNotNil(registration, t)

	//Register a user
	_, err = apiClient.RegisterUser("username-1", "password-1", []string{"username-password"})
	util.AssertNoError(err, t)

	redirectURI := url.QueryEscape("https://localhost:8080/redirect_location")
	scope := "scope1"
	state := util.GenerateRandomString(5)

	//Hit the authorization endpoint
	httpsClient := CreateHttpsClient()
	authorizeUrl := fmt.Sprintf(
		"https://localhost:8081/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		registration.Data.ClientId, redirectURI, scope, state)
	_, err = httpsClient.Get(authorizeUrl)
	util.AssertNoError(err, t)
	//util.AssertStatusCode(resp, 302, t)

	////Ensure that the user is authenticated and prompted for approval
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

func TestClientCredentialsEnd2End(t *testing.T) {

	go func() {
		httpServer := http.Server{Addr: ":8181"}
		http.HandleFunc("/test", func(w http.ResponseWriter, request *http.Request) {
			bearerHeader := request.Header.Get("Authorization")
			bearerTokenArray := regexp.MustCompile("Bearer ([a-zA-Z0-9]*)").FindStringSubmatch(string(bearerHeader))
			if len(bearerTokenArray) != 2 || !validateToken(bearerTokenArray[1]) {
				log.Printf("Bad token: size = %d value = %+v\n", len(bearerTokenArray)-1, bearerTokenArray)
				http.Error(w, "Unauthorized", 401)
			}

		})
		err := httpServer.ListenAndServeTLS("../resources/test/server.crt", "../resources/test/server.key")

		if err != nil {
			t.Fatalf("Failed to start resource server because of error %+v", err)
		}
	}()

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
	registration, err := apiClient.RegisterClient("test-client")
	util.AssertNoError(err, t)
	util.AssertNotNil(registration, t)

	//TODO: validate the client was registered

	apiClientInsufficientCreds, err := admin.NewAPIClient("https://localhost:8081",
		registration.Data.ClientId,
		registration.Data.ClientSecret,
		[]string{"../resources/test/server.crt"},
	)
	util.AssertNoError(err, t)

	_, err = apiClientInsufficientCreds.RegisterClient("test-client-2")
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

func validateToken(token string) bool {
	post, _ := http.NewRequest("POST", "https://localhost:8081/introspect",
		strings.NewReader(fmt.Sprintf("token=%s", token)))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.Header.Set("Authorization", "Basic "+creds.String())

	response, err := CreateHttpsClient().Do(post)

	if err != nil {
		log.Printf("Failed to execute introspection request: %+v\n", err)
		return false
	} else if response.StatusCode != 200 {
		log.Printf("Request to introspection endpoint failed with status code %d\n", response.StatusCode)
		return false
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("Failed to read http body: %+v\n", err)
			return false
		}
		responseJSON := map[string]interface{}{}
		err = json.Unmarshal(body, &responseJSON)
		validated := responseJSON["active"]
		log.Printf("Validated: %b", validated)
		return validated == true
	}
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

type HealthCheck struct {
	Status string `json:"status"`
}

func WaitForServer(server *AuthServer) error {
	var err error = nil
	var resp *http.Response = nil
	var healthcheck HealthCheck
	var body []byte

	httpClient := CreateHttpsClient()
	for i := 0; i < 25; i++ {
		resp, err = httpClient.Get("https://localhost:8081/health")

		if err == nil {
			if resp.StatusCode == 200 {
				body, err = ioutil.ReadAll(resp.Body)
				if err == nil {
					err = json.Unmarshal(body, &healthcheck)
					return err
				}
				return nil
			} else {
				err = fmt.Errorf("Expected status code 200 but was %s", resp.StatusCode)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	return err
}

func RunServer() (*AuthServer, *oauth_service.Credentials, error) {
	authServer := NewAuthServer(8081)

	var credentials *oauth_service.Credentials
	credentials, err := authServer.Start()

	if err := WaitForServer(authServer); err != nil {
		return nil, nil, err
	}

	return authServer, credentials, err
}

func TestMain(m *testing.M) {
	var authServer *AuthServer
	var err error

	authServer, creds, err = RunServer()
	if err != nil {
		log.Fatalf("Unexpected error on startup: %+v", err)
		return
	}
	result := m.Run()
	err = authServer.Stop()

	if err != nil {
		log.Fatalf("Unexpected error on stop: %+v", err)
		return
	}

	os.Exit(result)
}

func CreateHttpsClient() *http.Client {
	caCert, err := ioutil.ReadFile(TEST_CERTIFICATE)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		}}
}
