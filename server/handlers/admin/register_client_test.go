package admin

import (
	"github.com/danielsomerfield/authful/server/handlers"
	"encoding/json"
	"testing"
	"net/http"
	"github.com/danielsomerfield/authful/server/service/oauth"
	"os"
	"reflect"
)

var mockClientAccessControlFn = func(request http.Request) (bool, error) {
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
		clientId: clientId,
		clientSecret: clientSecret,
		name : name,
		scopes: scopes,
	}
	return &oauth.Credentials{
		ClientId: clientId,
		ClientSecret: clientSecret,
	}, nil
}

func TestMain(m *testing.M) {
	registeredClients = map[string]registeredClient{}
	os.Exit(m.Run())
}

func TestRegisterClientHandler_registersClientWithValidCredentials(t *testing.T) {


	registerClientRequest := map[string]interface{} {
		"command": map[string]interface{}{
			"name":   "test-client",
			"scopes": []string{"scope-1", "scope-2"},
		},
	}

	body, _ := json.Marshal(registerClientRequest)
	response := handlers.DoEndpointRequest(
		NewRegisterClientHandler(mockClientAccessControlFn, mockRegisterClientFn), string(body))

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
		clientId: createdClientId,
		clientSecret: createdClientSecret,
		name: "test-client",
		scopes: []string{"scope-1", "scope-2"},
	}

	if len(registeredClients) != 1 || reflect.DeepEqual(registeredClients["test-client"], expectedClient) {
		t.Fatalf("Expected client %+v to exist in %+v\n", expectedClient, registeredClients)
	}
}

//Test that without correct scope, request fails
//Test that without credentials, request fails
//Test that registering the same client twice fails
