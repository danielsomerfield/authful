package user

import (
	"testing"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/handlers"
	"github.com/danielsomerfield/authful/common/util"
	"fmt"
	"github.com/danielsomerfield/authful/server/service/admin/user"
	"errors"
	"github.com/danielsomerfield/authful/server/service/oauth"
)

var registeredUsers = map[string]user.User{}

var mockRegisterUserFn = func(user user.User) error {
	registeredUsers[user.Username] = user
	return nil
}

func setup() {
	registeredUsers = map[string]user.User{}
}

func TestRegisterUserHandler_registers_valid_user(t *testing.T) {
	setup()
	registerRequest := map[string]interface{}{
		"command": map[string]interface{}{
			"username":    "user1",
			"password":    "user1password",
			"authMethods": []string{"username-password"},
		},
	}
	body, _ := json.Marshal(registerRequest)

	expectedUser := user.User{
		Username:    "user1",
		Password:    "user1password",
		AuthMethods: []string{"username-password"},
	}

	handlers.DoPostEndpointRequest(NewRegisterUserHandler(mockRegisterUserFn), string(body)).
		ThenAssert(func(response *handlers.JSONEndpointResponse) error {
		util.AssertTrue(response.HttpStatus == 201, "Http status is 201", t)
		util.AssertTrue(len(registeredUsers) == 1, "There is 1 registered user", t)
		util.AssertEquals(expectedUser, registeredUsers["user1"], t)
		return nil
	}, t)

	//TODO: test response
}

func TestRegisterUserHandler_malformed_message_fails(t *testing.T) {
	runMalformedMessageTest(map[string]interface{}{
		"password":    "user1password",
		"authMethods": []string{"username-password"},
	}, 400, "invalid_request", t)

	runMalformedMessageTest(map[string]interface{}{
		"username":    "user1",
		"authMethods": []string{"username-password"},
	}, 400, "invalid_request", t)

	runMalformedMessageTest(map[string]interface{}{
		"username": "user1",
		"password": "user1password",
	}, 400, "invalid_request", t)
}

func runMalformedMessageTest(command map[string]interface{}, expectedCode int, expectedErrorType string, t *testing.T) {
	setup()

	registerRequest := map[string]interface{}{
		"command": command,
	}
	body, _ := json.Marshal(registerRequest)

	handlers.DoPostEndpointRequest(NewRegisterUserHandler(mockRegisterUserFn), string(body)).
		ThenAssert(func(response *handlers.JSONEndpointResponse) error {

		errorResponse := response.Json
		errorJson, converted := errorResponse["error"].(map[string]interface{})

		util.AssertTrue(converted, "Error message exists in response", t)

		util.AssertTrue(errorJson != nil, "There is an error in the response", t)

		util.AssertTrue(response.HttpStatus == expectedCode, fmt.Sprintf("Http status is %d", expectedCode), t)
		util.AssertEquals(errorJson["errorType"], expectedErrorType, t)
		util.AssertTrue(len(registeredUsers) == 0, "There should be no registered users", t)
		return nil
	}, t)
}

func TestRegisterUserHandler_access_control_fails(t *testing.T) {
	setup()

	registerRequest := map[string]interface{}{
		"command": map[string]interface{}{
			"username":    "user1",
			"password":    "user1password",
			"authMethods": []string{"username-password"},
		},
	}

	body, _ := json.Marshal(registerRequest)

	lookup := func(clientId string) (oauth.Client, error) {
		return nil, nil
	}

	handlers.DoPostEndpointRequest(NewProtectedHandler(mockRegisterUserFn, lookup), string(body)).
		ThenAssert(func(response *handlers.JSONEndpointResponse) error {

		errorResponse := response.Json
		errorJson, converted := errorResponse["error"].(map[string]interface{})

		util.AssertTrue(converted, "Error message exists in response", t)

		util.AssertTrue(errorJson != nil, "There is an error in the response", t)

		util.AssertTrue(response.HttpStatus == 401, fmt.Sprintf("Http status is 401"), t)
		util.AssertEquals(errorJson["errorType"], "invalid_client", t)
		util.AssertTrue(len(registeredUsers) == 0, "There should be no registered users", t)
		return nil
	}, t)
}

func TestRegisterUserHandler_registration_fails(t *testing.T) {
	setup()

	registerRequest := map[string]interface{}{
		"command": map[string]interface{}{
			"username":    "user1",
			"password":    "user1password",
			"authMethods": []string{"username-password"},
		},
	}

	body, _ := json.Marshal(registerRequest)

	var mockFailingRegisterUserFn = func(user user.User) error {
		return errors.New("Failed for some reason")
	}

	handlers.DoPostEndpointRequest(NewRegisterUserHandler(mockFailingRegisterUserFn), string(body)).
		ThenAssert(func(response *handlers.JSONEndpointResponse) error {

		errorResponse := response.Json
		errorJson, converted := errorResponse["error"].(map[string]interface{})

		util.AssertTrue(converted, "Error message exists in response", t)
		util.AssertTrue(errorJson != nil, "There is an error in the response", t)
		util.AssertTrue(response.HttpStatus == 400, fmt.Sprintf("Http status is 401"), t)
		util.AssertEquals(errorJson["errorType"], "invalid_request", t)
		return nil
	}, t)
}
