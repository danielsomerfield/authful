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

func AssertNotNil(pointer interface{}, t *testing.T) {
	if pointer == nil {
		t.Fatal("Received an unexpected nil pointer")
	}
}

func AssertTrue(boolean bool, message string, t *testing.T) {
	if !boolean {
		t.Fatalf("Expected %s to be true but was false", message)
	}
}

func AssertStatusCode(r *http.Response, expected int, t *testing.T) {
	if r.StatusCode != expected {
		t.Fatalf("Expected HTTP status %d but received %d", expected, r.StatusCode)
	}
}