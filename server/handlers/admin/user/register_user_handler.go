package user

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/handlers"
	wireUser "github.com/danielsomerfield/authful/server/wire/admin/user"
	"github.com/danielsomerfield/authful/server/service/admin/user"
	"log"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/wire"
	"github.com/danielsomerfield/authful/server/service/oauth"
)

func NewProtectedHandler(registerUserFn user.RegisterUserFn, lookup oauth.ClientLookupFn) http.HandlerFunc {
	return handlers.Protect(NewRegisterUserHandler(registerUserFn), "administrate", lookup)
}

func NewRegisterUserHandler(registerUserFn user.RegisterUserFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
			bytes, err := json.Marshal(wire.ResponseEnvelope{
				Data: RegisterUserResponse{},
			})
			handlers.WriteOrInternalError(w, bytes, err)
		}
	}
}

//TODO: fill out the data here
type RegisterUserResponse struct {
}
