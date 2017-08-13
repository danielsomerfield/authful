package handlers

import (
	"strings"
	"net/http/httptest"
	"io/ioutil"
	"bytes"
	"testing"
	"net/http"
	"encoding/json"
	"github.com/danielsomerfield/authful/util"
	"fmt"
)

type JSONEndpointResponse struct {
	Json       map[string]interface{}
	HttpStatus int
	Err        error
	t		   *testing.T
}

//TODO: refactor the harness for JSON and non JSON cases

func (r *JSONEndpointResponse) AssertHttpStatusEquals(httpStatus int) {
	util.AssertEquals(r.HttpStatus, httpStatus, r.t)
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

func DoGetEndpointRequest(underTest http.HandlerFunc, urlstring string) *EndpointResponse {

	request, _ := http.NewRequest("Get", urlstring, nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(underTest)
	handler.ServeHTTP(response, request)

	return &EndpointResponse{
		httpStatus: response.Code,
		err:        nil,
	}

}

type EndpointResponse struct {
	httpStatus int
	err        error
	t          *testing.T
}

func (r *EndpointResponse) AssertHttpStatusEquals(httpStatus int) {
	util.AssertTrue(r.httpStatus == httpStatus, fmt.Sprintf("Expected http status %d but found %s", httpStatus, r.httpStatus), r.t)
}

func (r *EndpointResponse) ThenAssert(test func(response *EndpointResponse) error, t *testing.T) error {
	r.t = t
	if r.err != nil {
		t.Errorf("Request failed: %+v", r.err)
		return r.err
	}

	err := test(r)
	if err != nil {
		t.Errorf("Assertion failed: %+v", err)
		return r.err
	}
	return nil
}

func (rs *JSONEndpointResponse) ThenAssert(test func(response *JSONEndpointResponse) error, t *testing.T) error {
	rs.t = t
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
