package user

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/handlers"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"errors"
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

		registerUserRequest, err := ParseRegisterUserRequest(r)
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

func ParseRegisterUserRequest(r *http.Request) (*RegisterUserRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		request := &RegisterUserRequest{}
		err = json.Unmarshal(body, request)
		if err != nil {
			log.Printf("Failed to unmarshal RegisterUserRequest: %s", string(body))
		} else {
			if strings.TrimSpace(request.Command.Username) == "" {
				return nil, errors.New("Missing field 'username'")
			}

			if strings.TrimSpace(request.Command.Password) == "" {
				return nil, errors.New("Missing field 'password'")
			}

			if len(request.Command.AuthMethods) == 0 {
				return nil, errors.New("Missing field 'authMethods'")
			}

			return request, nil
		}
	}

	return nil, err
}

type RegisterUserRequest struct {
	Command RegisterUserCommand `json:"command,omitempty"`
}

type RegisterUserCommand struct {
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	AuthMethods []string    `json:"authMethods"`
}
