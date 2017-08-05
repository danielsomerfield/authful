package admin

import (
	"net/http"
	"github.com/danielsomerfield/authful/server/handlers"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
)

//TODO: move this to the service area
type RegisterUserFn func(user User) error
type User struct {
	username string
	password string
	scopes   []string
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
			username: command.Username,
			password: command.Password,
			scopes:   command.Scopes,
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
			fmt.Printf("wire: %s", string(body))
			fmt.Printf("request: %+v", request)

			return request, nil
		}
	}

	return nil, err
}

type RegisterUserRequest struct {
	Command RegisterUserCommand `json:"command,omitempty"`
}

type RegisterUserCommand struct {
	Username string 	`json:"username,omitempty"`
	Password string 	`json:"password,omitempty"`
	Scopes   []string	`json:"scopes"`
}


