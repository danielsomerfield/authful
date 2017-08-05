package user

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/handlers"
	wireUser "github.com/danielsomerfield/authful/server/wire/admin/user"
	"github.com/danielsomerfield/authful/server/service/accesscontrol"
)

//TODO: move this to the service area
type RegisterUserFn func(user User) error
type User struct {
	username    string
	password    string
	authMethods []string
}

func NewRegisterUserHandler(accessControlFn accesscontrol.ClientAccessControlFn, registerUserFn RegisterUserFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authorized, err := accessControlFn(*r)

		if !authorized {
			if err != nil {
				handlers.InternalServerError("An unexpected error occurred", w)
			} else {
				handlers.Unauthorized("The requested operation was denied.", w)
			}
			return
		}

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
