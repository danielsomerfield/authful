package handlers

import (
	"strings"
	"net/http/httptest"
	"io/ioutil"
	"bytes"
	"testing"
	"net/http"
	"encoding/json"
)

type JSONEndpointResponse struct {
	Json       map[string]interface{}
	HttpStatus int
	Err        error
}

func DoPostEndpointRequest(underTest http.HandlerFunc, body string) *JSONEndpointResponse {
	return DoPostEndpointRequestWithHeaders(underTest, body, map[string]string{})
}


func DoPostEndpointRequestWithHeaders(underTest http.HandlerFunc, body string, headers map[string]string) *JSONEndpointResponse {

	post, _ := http.NewRequest("POST", "",
		strings.NewReader(body))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for name, value := range headers {
		post.Header.Set(name, value)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(underTest)
	handler.ServeHTTP(response, post)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &JSONEndpointResponse{
			Err: err,
		}
	}

	var jwt map[string]interface{}
	decoder := json.NewDecoder(bytes.NewBuffer(responseBody))
	decoder.UseNumber()
	decoder.Decode(&jwt)
	if err != nil {
		return &JSONEndpointResponse{
			Err: err,
		}
	}
	return &JSONEndpointResponse{
		Json:       jwt,
		HttpStatus: response.Code,
		Err:        nil,
	}

}

func DoGetEndpointRequestWithHeaders(urlstring string, underTest http.HandlerFunc, body string,
	headers map[string]string) *JSONEndpointResponse {

	post, _ := http.NewRequest("Get", urlstring,
		strings.NewReader(body))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for name, value := range headers {
		post.Header.Set(name, value)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(underTest)
	handler.ServeHTTP(response, post)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &JSONEndpointResponse{
			Err: err,
		}
	}

	var jwt map[string]interface{}
	decoder := json.NewDecoder(bytes.NewBuffer(responseBody))
	decoder.UseNumber()
	decoder.Decode(&jwt)
	if err != nil {
		return &JSONEndpointResponse{
			Err: err,
		}
	}
	return &JSONEndpointResponse{
		Json:       jwt,
		HttpStatus: response.Code,
		Err:        nil,
	}

}

func (rs *JSONEndpointResponse) ThenAssert(test func(response *JSONEndpointResponse) error, t *testing.T) error {
	if rs.Err != nil {
		t.Errorf("Request failed: %+v", rs.Err)
		return rs.Err
	}

	err := test(rs)
	if err != nil {
		t.Errorf("Assertion failed: %+v", err)
		return rs.Err
	}
	return nil
}
