package client

import (
	"github.com/danielsomerfield/authful/server/handlers"
	"encoding/json"
	"testing"
	"net/http"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"reflect"
	"github.com/danielsomerfield/authful/util"
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
	redirectUris  []string
}

var registeredClients = map[string]registeredClient{}

var mockRegisterClientFn = func(name string, scopes []string, redirectUris []string) (*oauth.Credentials, error) {
	clientId := name + "-id"
	clientSecret := name + "-secret"
	registeredClients[clientId] = registeredClient{
		clientId:     clientId,
		clientSecret: clientSecret,
		name:         name,
		scopes:       scopes,
		redirectUris: redirectUris,
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
			"name":          "test-client",
			"scopes":        []string{"scope-1", "scope-2"},
			"redirect_uris": []string{"http://example.com/loggedIn"},
		},
	}

	body, _ := json.Marshal(registerClientRequest)
	response := handlers.DoPostEndpointRequest(
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
		redirectUris:  []string{"http://example.com/loggedIn"},
	}

	if len(registeredClients) != 1 || !reflect.DeepEqual(registeredClients[createdClientId], expectedClient) {
		t.Fatalf("Expected client %+v to exist in %+v\n", expectedClient, registeredClients)
	}
}

func TestRegisterClientHandler_registerReturnsErrorWithFailingAuthorization(t *testing.T) {

	setup()
	registerClientRequest := map[string]interface{}{
		"command": map[string]interface{}{
			"name":          "test-client",
			"scopes":        []string{"scope-1", "scope-2"},
			"redirect_uris": []string{"http://example.com/loggedIn"},
		},
	}

	body, _ := json.Marshal(registerClientRequest)
	handlers.DoPostEndpointRequest(
		NewRegisterClientHandler(mockFailingClientAccessControlFn, mockRegisterClientFn), string(body)).
		ThenAssert(func(r *handlers.JSONEndpointResponse) error {

		r.AssertHttpStatusEquals(401)

		theError, converted := r.Json["error"].(map[string]interface{})
		util.AssertTrue(converted, "Failed to convert to expected type.", t)

		if theError["status"] != 401 && theError["errorType"] != "invalid_client" {
			t.Fatalf("Unexpected error: %+v\n", theError)
		}

		if len(registeredClients) != 0 {
			t.Fatalf("Expected no clients to be registered: %+v\n", registeredClients)
		}

		return nil
	}, t)

}

func TestRegisterClientHandler_registerReturnsErrorWithNoProvidedName(t *testing.T) {

	setup()
	registerClientRequest := map[string]interface{}{
		"command": map[string]interface{}{
			"scopes":        []string{"scope-1", "scope-2"},
			"redirect_uris": []string{"http://example.com/loggedIn"},
		},
	}

	body, _ := json.Marshal(registerClientRequest)
	handlers.DoPostEndpointRequest(
		NewRegisterClientHandler(mockSucceedingClientAccessControlFn, mockRegisterClientFn), string(body)).
		ThenAssert(func(response *handlers.JSONEndpointResponse) error {

		response.AssertHttpStatusEquals(400)

		if len(registeredClients) != 0 {
			t.Fatalf("Expected no clients to be registered: %+v\n", registeredClients)
		}
		return nil
	}, t)

}

//Test that registering the same client twice fails
