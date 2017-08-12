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
)

const RESOURCE_SERVER_PORT = ":8181"
const TEST_CERTIFICATE = "../resources/test/server.crt"

var resourceServer = http.Server{Addr: RESOURCE_SERVER_PORT}

func StartResourceServer(pattern string, handler http.HandlerFunc) {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc(pattern, handler)
		resourceServer.Handler = mux
		err := resourceServer.ListenAndServeTLS("../resources/test/server.crt", "../resources/test/server.key")

		if err != nil {
			log.Fatalf("Failed to start resource server because of error %+v", err)
		}
	}()
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
