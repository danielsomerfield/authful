package util

import (
	"testing"
	"net/http"
	"runtime"
	"log"
	"reflect"
)

func AssertNoError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("No error expected, but received %+v\n", err)
		PrintStackTrace()
	}
}

func AssertNotNil(pointer interface{}, t *testing.T) {
	if pointer == nil {
		t.Errorf("Received an unexpected nil pointer")
		PrintStackTrace()
	}
}

func AssertTrue(boolean bool, message string, t *testing.T) {
	if !boolean {
		t.Errorf("Expected \"%s\" to be true but was false", message)
		PrintStackTrace()
	}
}

func AssertFalse(boolean bool, message string, t *testing.T) {
	if boolean {
		t.Errorf("Expected \"%s\" to be false but was true", message)
		PrintStackTrace()
	}
}

func AssertEquals(expectedValue interface{}, actualValue interface{}, t *testing.T) {
	if !reflect.DeepEqual(expectedValue, actualValue) {
		t.Errorf("Expected value +%v but was +%v", expectedValue, actualValue)
		PrintStackTrace()
	}
}

func AssertStatusCode(r *http.Response, expected int, t *testing.T) {
	if r.StatusCode != expected {
		t.Errorf("Expected HTTP status %d but received %d", expected, r.StatusCode)
		PrintStackTrace()
	}
}

func PrintStackTrace() {
	var stack [4096]byte
	runtime.Stack(stack[:], false)
	log.Printf("%s\n", stack[:])
}
