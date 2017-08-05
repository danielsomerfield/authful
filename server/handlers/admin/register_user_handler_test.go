package admin

import (
	"testing"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/handlers"
	"github.com/danielsomerfield/authful/util"
)

//type registeredUser struct {
//	userId       string
//	username     string
//	passwordHash string
//	scopes       []string
//}


var registeredUsers = map[string]User{}

var mockRegisterUserFn = func(user User) error {
	registeredUsers[user.username] = user
	return nil
}


func TestRegisterUserHandler_registers_valid_user(t *testing.T) {
	registerRequest := map[string]interface{}{
		"command": map[string]interface{}{
			"username": "user1",
			"password": "user1password",
			"scopes":   []string{"foo"},
		},
	}
	body, _ := json.Marshal(registerRequest)

	expectedUser := User{
		username: "user1",
		password: "user1password",
		scopes: []string{"foo"},
	}

	handlers.DoEndpointRequest(NewRegisterUserHandler(mockRegisterUserFn), string(body)).
		ThenAssert(func(response *handlers.EndpointResponse) error {
			util.AssertTrue(response.HttpStatus == 201, "Http status is 201", t)
			util.AssertTrue(len(registeredUsers) == 1, "There is 1 registered user", t)
			util.AssertEquals(expectedUser, registeredUsers["user1"], t)
			return nil
		}, t)
}

//Duplicate username fails
//Malformed message (missing fields)
//Missing scopes
//Invalid token
