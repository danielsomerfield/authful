package user

import (
	"testing"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/handlers"
	"github.com/danielsomerfield/authful/util"
	"fmt"
)

var registeredUsers = map[string]User{}

var mockRegisterUserFn = func(user User) error {
	registeredUsers[user.username] = user
	return nil
}

func setup() {
	registeredUsers = map[string]User{}
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

	expectedUser := User{
		username:    "user1",
		password:    "user1password",
		authMethods: []string{"username-password"},
	}

	handlers.DoEndpointRequest(NewRegisterUserHandler(mockRegisterUserFn), string(body)).
		ThenAssert(func(response *handlers.EndpointResponse) error {
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
		"username":    "user1",
		"password":    "user1password",
	}, 400, "invalid_request", t)
}

func runMalformedMessageTest(command map[string]interface{}, expectedCode int, expectedErrorType string, t *testing.T) {
	setup()

	registerRequest := map[string]interface{}{
		"command": command,
	}
	body, _ := json.Marshal(registerRequest)

	handlers.DoEndpointRequest(NewRegisterUserHandler(mockRegisterUserFn), string(body)).
		ThenAssert(func(response *handlers.EndpointResponse) error {

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

//TODO: Failed access control
