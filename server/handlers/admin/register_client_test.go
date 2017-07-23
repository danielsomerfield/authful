package admin

import (
	"github.com/danielsomerfield/authful/server/handlers"
	"encoding/json"
	"testing"
	"net/http"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"reflect"
)

var mockSucceedingClientAccessControlFn = func(request http.Request) (bool, error) {
	return true, nil
}

var mockFailingClientAccessControlFn = func(request http.Request) (bool, error) {
	return false, nil
}

type registeredClient struct {
	clientId     string
	clientSecret string
	name         string
	scopes       []string
}

var registeredClients = map[string]registeredClient{}

var mockRegisterClientFn = func(name string, scopes []string) (*oauth.Credentials, error) {
	clientId := name + "-id"
	clientSecret := name + "-secret"
	registeredClients[clientId] = registeredClient{
		clientId:     clientId,
		clientSecret: clientSecret,
		name:         name,
		scopes:       scopes,
	}
	return &oauth.Credentials{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}, nil
}

func setup() {
	registeredClients = map[string]registeredClient{}
}

func TestRegisterClientHandler_registersClientWithValidCredentials(t *testing.T) {

	setup()
	registerClientRequest := map[string]interface{}{
		"command": map[string]interface{}{
			"name":   "test-client",
			"scopes": []string{"scope-1", "scope-2"},
		},
	}

	body, _ := json.Marshal(registerClientRequest)
	response := handlers.DoEndpointRequest(
		NewRegisterClientHandler(mockSucceedingClientAccessControlFn, mockRegisterClientFn), string(body))

	if response.HttpStatus != 200 {
		t.Fatalf("Expected 200 but got %d\n", response.HttpStatus)
	}

	dataField, converted := response.Json["data"].(map[string]interface{})

	if !converted {
		t.Fatalf("Failed to convert to expected type.")
	}

	createdClientId := "test-client-id"
	createdClientSecret := "test-client-secret"
	if dataField["clientId"] != createdClientId && dataField["clientSecret"] != createdClientSecret {
		t.Fatalf("Received unexpected payload %+v\n", response.Json)
	}

	expectedClient := registeredClient{
		clientId:     createdClientId,
		clientSecret: createdClientSecret,
		name:         "test-client",
		scopes:       []string{"scope-1", "scope-2"},
	}

	if len(registeredClients) != 1 || reflect.DeepEqual(registeredClients["test-client"], expectedClient) {
		t.Fatalf("Expected client %+v to exist in %+v\n", expectedClient, registeredClients)
	}
}

func TestRegisterClientHandler_registerReturnsErrorWithFailingAuthorization(t *testing.T) {

	setup()
	registerClientRequest := map[string]interface{}{
		"command": map[string]interface{}{
			"name":   "test-client",
			"scopes": []string{"scope-1", "scope-2"},
		},
	}

	body, _ := json.Marshal(registerClientRequest)
	response := handlers.DoEndpointRequest(
		NewRegisterClientHandler(mockFailingClientAccessControlFn, mockRegisterClientFn), string(body))

	if response.HttpStatus != 401 {
		t.Fatalf("Expected 401 but got %d\n", response.HttpStatus)
	}

	theError, converted := response.Json["error"].(map[string]interface{})

	if !converted {
		t.Fatalf("Failed to convert to expected type.")
	}

	if !converted {
		t.Fatalf("Failed to convert to expected type.")
	}

	if theError["status"] != 401 && theError["errorType"] != "invalid_client" {
		t.Fatalf("Unexpected error: %+v\n", theError)
	}

	if len(registeredClients) != 0 {
		t.Fatalf("Expected no clients to be registered: %+v\n", registeredClients)
	}
}

func TestRegisterClientHandler_registerReturnsErrorWithNoProvidedName(t *testing.T) {

	setup()
	registerClientRequest := map[string]interface{}{
		"command": map[string]interface{}{
			"scopes": []string{"scope-1", "scope-2"},
		},
	}

	body, _ := json.Marshal(registerClientRequest)
	response := handlers.DoEndpointRequest(
		NewRegisterClientHandler(mockSucceedingClientAccessControlFn, mockRegisterClientFn), string(body))

	if response.HttpStatus != 400 {
		t.Fatalf("Expected 400 but got %d\n", response.HttpStatus)
	}

	if len(registeredClients) != 0 {
		t.Fatalf("Expected no clients to be registered: %+v\n", registeredClients)
	}
}

//Test that registering the same client twice fails
