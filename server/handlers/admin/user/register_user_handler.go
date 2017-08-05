package user

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/handlers"
	wireUser "github.com/danielsomerfield/authful/server/wire/admin/user"
	"github.com/danielsomerfield/authful/server/service/accesscontrol"
	"github.com/danielsomerfield/authful/server/service/admin/user"
	"log"
)


func NewRegisterUserHandler(accessControlFn accesscontrol.ClientAccessControlFn, registerUserFn user.RegisterUserFn) http.HandlerFunc {
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
		user := user.User{
			Username:    command.Username,
			Password:    command.Password,
			AuthMethods: command.AuthMethods,
		}

		err = registerUserFn(user)
		if err != nil {
			log.Printf("Failed with following error: %+v", err)
			handlers.InvalidRequest("Failed to register the user.", w)
		} else {
			w.WriteHeader(http.StatusCreated)
			//TODO: write the register response
		}

	}
}
