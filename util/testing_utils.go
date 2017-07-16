package util

import (
	"testing"
	"net/http"
)

func AssertNoError(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("No error expected, but received %+v\n", err)
	}
}

func AssertStatusCode(r *http.Response, expected int, t *testing.T) {
	if r.StatusCode != expected {
		t.Fatalf("Expected HTTP status %d but received %d", expected, r.StatusCode)
	}
}