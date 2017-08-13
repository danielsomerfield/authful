package server

import (
	"log"
	"io/ioutil"
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	oauth_service "github.com/danielsomerfield/authful/server/service/oauth"
	"crypto/x509"
	"crypto/tls"
	"regexp"
	"strings"
)

const RESOURCE_SERVER_PORT = ":8181"
const TEST_CERTIFICATE = "../resources/test/server.crt"

var resourceServer = http.Server{Addr: RESOURCE_SERVER_PORT}
var creds *oauth_service.Credentials = nil

func StartResourceServer() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/test", func(w http.ResponseWriter, request *http.Request) {
			bearerHeader := request.Header.Get("Authorization")
			bearerTokenArray := regexp.MustCompile("Bearer ([a-zA-Z0-9]*)").FindStringSubmatch(string(bearerHeader))
			if len(bearerTokenArray) != 2 || !validateToken(bearerTokenArray[1]) {
				log.Printf("Bad token: size = %d value = %+v\n", len(bearerTokenArray)-1, bearerTokenArray)
				http.Error(w, "Unauthorized", 401)
			}
		})

		mux.HandleFunc("/ping", func(w http.ResponseWriter, request *http.Request) {

		})

		mux.HandleFunc("/error", func(w http.ResponseWriter, request *http.Request) {
			fmt.Printf("=============++> /error")
		})

		resourceServer.Handler = mux
		err := resourceServer.ListenAndServeTLS("../resources/test/server.crt", "../resources/test/server.key")

		if err != nil {
			log.Fatalf("Failed to start resource server because of error %+v", err)
		}
	}()
}

func StopResourceServer() {
	resourceServer.Shutdown(nil)
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

type HealthCheck struct {
	Status string `json:"status"`
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
