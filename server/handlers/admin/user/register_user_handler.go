package user

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/handlers"
	wireUser "github.com/danielsomerfield/authful/server/wire/admin/user"
)

//TODO: move this to the service area
type RegisterUserFn func(user User) error
type User struct {
	username    string
	password    string
	authMethods []string
}

func NewRegisterUserHandler(registerUserFn RegisterUserFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		registerUserRequest, err := wireUser.ParseRegisterUserRequest(r)
		if err != nil {
			handlers.InvalidRequest("Failed to parse request to register user", w)
			return
		}

		command := registerUserRequest.Command
		user := User{
			username:    command.Username,
			password:    command.Password,
			authMethods: command.AuthMethods,
		}

		registerUserFn(user)

		w.WriteHeader(http.StatusCreated)
	}
}
