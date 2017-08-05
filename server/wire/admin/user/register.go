package user

import (
	"log"
	"strings"
	"errors"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

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

