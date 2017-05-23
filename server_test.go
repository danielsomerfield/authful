package main

import (
	"net/http"
	"testing"
	"io/ioutil"
	"encoding/json"
	"authful/server"
	"fmt"
	"time"
)

func TestMain(t *testing.T) {
	authServer := server.NewAuthServer()

	go authServer.Start()

	defer authServer.Stop()
	healthcheck, err := waitForServer(authServer)

	if err != nil {
		t.Error(err)
	} else if healthcheck.Status != "ok" {
		t.Errorf("Message should have been 'ok' but was %s", healthcheck.Status)
	}
}

func waitForServer(server *server.AuthServer) (HealthCheck, error) {

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
					return healthcheck, err
				}
			} else {
				err = fmt.Errorf("Expected status code 200 but was %s", resp.StatusCode)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	return healthcheck, err
}

type HealthCheck struct {
	Status string `json:"status"`
}
