package main

import (
	"testing"
	"github.com/danielsomerfield/authful/testutils"
	"io/ioutil"
	"net/http"
)

func TestAuthorize(t *testing.T) {

	var resp *http.Response = nil
	var err error = nil
	var body []byte

	err = testutils.RunServer()
	defer testutils.StopServer()

	resp, err = http.Get("http://localhost:8080/authorize?request_type=code?client_id=1234")

	if err == nil {
		if resp.StatusCode == 200 {
			body, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				print(string(body))
			}
		} else {
			t.Errorf("Expected status code 200 but was %s", resp.StatusCode)
		}
	} else {

	}
}
