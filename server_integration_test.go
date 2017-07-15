package main

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
	"github.com/danielsomerfield/authful/server"
	"os"
	"log"
	"regexp"
)

var creds *oauth_service.Credentials = nil

func TestMain(m *testing.M) {
	var authServer *server.AuthServer
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

func requestAdminToken(credentials oauth_service.Credentials) (*oauth_wire.TokenResponse, error) {
	var err error = nil

	var request *http.Request
	if request, err = http.NewRequest("POST", "http://localhost:8081/token",
		strings.NewReader("grant_type=client_credentials")); err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Basic "+credentials.String())

	var response *http.Response
	var body []byte

	response, err = http.DefaultClient.Do(request)

	body, err = ioutil.ReadAll(response.Body)

	tokenResponse := oauth_wire.TokenResponse{}
	err = json.Unmarshal(body, &tokenResponse)
	if err == nil {
		return &tokenResponse, nil
	} else {
		return nil, err
	}
}

func TestClientCredentialsEnd2End(t *testing.T) {
	go func() {
		httpServer := http.Server{Addr: ":8181"}
		http.HandleFunc("/test", func(w http.ResponseWriter, request *http.Request) {
			//body, err := ioutil.ReadAll(request.Body)
			//fmt.Printf("/test: body = %+v err = %+v headers = %+v\n", body, err, request.Header)
			bearerHeader := request.Header.Get("Authorization")
			bearerTokenArray := regexp.MustCompile("Bearer ([a-zA-Z0-9]*)").FindStringSubmatch(string(bearerHeader))
			if len(bearerTokenArray) != 2 || !validateToken(bearerTokenArray[1]) {
				log.Printf("Bad token: size = %d value = %+v\n", len(bearerTokenArray)-1, bearerTokenArray)
				http.Error(w, "Unauthorized", 401)
			}

		})
		httpServer.ListenAndServe()

	}()

	ctx := context.Background()
	config := clientcredentials.Config{
		ClientID:     creds.ClientId,
		ClientSecret: creds.ClientSecret,
		TokenURL:     "http://localhost:8081/token",
		Scopes:       []string{},
	}

	resp, err := config.Client(ctx).Get("http://localhost:8181/test")
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	} else if resp.StatusCode != 200 {
		t.Errorf("Expected 200 but was %d\n", resp.StatusCode)
		return
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	}

	//No token
	resp, err = http.Get("http://localhost:8181/test")
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	} else if resp.StatusCode != 401 {
		t.Errorf("Expected 401 but was %d\n", resp.StatusCode)
		return
	}

	//Invalid token
	invalidRequest, _ := http.NewRequest("GET", "http://localhost:8181/test", nil)
	invalidRequest.Header.Set("Authorization", "Bearer 12345")
	resp, err = http.DefaultClient.Do(invalidRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
		return
	} else if resp.StatusCode != 401 {
		t.Errorf("Expected 401 but was %d\n", resp.StatusCode)
		return
	}
}

func validateToken(token string) bool {
	post, _ := http.NewRequest("POST", "http://localhost:8081/introspect",
		strings.NewReader(fmt.Sprintf("token=%s", token)))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.Header.Set("Authorization", "Basic "+creds.String())

	response, err := http.DefaultClient.Do(post)

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

func TestErrorResponse(t *testing.T) {

	response, err := http.Post("http://localhost:8081/token", "application/json", strings.NewReader(""))
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

func TestAuthorize(t *testing.T) {

	_, err := requestAdminToken(*creds)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
	//fmt.Printf("Token: %+v, %+v", token, err)
	return

	//Register a client and get back client id and secret

	//Register a user

	//Hit the authorization endpoint
	//resp, err = http.Get("http://localhost:8081/authorize?request_type=code?client_id=1234")
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

type HealthCheck struct {
	Status string `json:"status"`
}

func WaitForServer(server *server.AuthServer) error {
	var err error = nil
	var resp *http.Response = nil
	var healthcheck HealthCheck
	var body []byte

	for i := 0; i < 25; i++ {
		resp, err = http.Get("http://localhost:8081/health")

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

func RunServer() (*server.AuthServer, *oauth_service.Credentials, error) {
	authServer := server.NewAuthServer(8081)

	var credentials *oauth_service.Credentials
	credentials, err := authServer.Start()

	if err := WaitForServer(authServer); err != nil {
		return nil, nil, err
	}

	return authServer, credentials, err
}
