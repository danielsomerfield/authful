package testutils

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
	"github.com/danielsomerfield/authful/server"
	"net/http"
)

type HealthCheck struct {
	Status string `json:"status"`
}

func WaitForServer(server *server.AuthServer) error {

	var err error = nil
	var resp *http.Response = nil
	var healthcheck HealthCheck
	var body []byte

	for i := 0; i < 1000; i++ {
		resp, err = http.Get("http://localhost:8080/health")

		if err == nil {
			if resp.StatusCode == 200 {
				body, err = ioutil.ReadAll(resp.Body)
				if err == nil {
					err = json.Unmarshal(body, &healthcheck)
					return err
				}
			} else {
				err = fmt.Errorf("Expected status code 200 but was %s", resp.StatusCode)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	return err
}

func RunServer() (*server.AuthServer, *server.Credentials, error) {
	authServer := server.NewAuthServer()

	var credentials *server.Credentials
	credentials = authServer.Start()
	err := WaitForServer(authServer)
	return authServer, credentials, err
}
