package admin

import (
	"fmt"
	"net/url"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/danielsomerfield/authful/common/wire"
)

type ClientRegistration struct {
	ClientId     string
	ClientSecret string
}

type ClientError struct {
	errorType    string
	errorMessage string
}

func (ce ClientError) Error() string {
	return fmt.Sprintf("%s : %s", ce.errorType, ce.errorMessage)
}

func (ce ClientError) ErrorType() string {
	return ce.errorType
}

func RegisterClient(clientId string, clientSecret string, clientName string, httpClient *http.Client) (*ClientRegistration, error) {
	credentials := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s:%s", url.QueryEscape(clientId), url.QueryEscape(clientSecret))))

	createClientRequest := map[string]interface{}{
		"command": map[string]string{
			"name": clientName,
		},
	}
	messageBytes, _ := json.Marshal(createClientRequest)

	post, _ := http.NewRequest("POST", "https://localhost:8081/admin/clients",
		bytes.NewReader(messageBytes))
	post.Header.Set("Content-Type", "application/json")
	post.Header.Set("Authorization", "Basic "+credentials)

	response, e := httpClient.Do(post)
	if e != nil {
		return nil, ClientError{"ClientError", e.Error()}
	}

	body, e := ioutil.ReadAll(response.Body)

	if e != nil {
		return nil, ClientError{"ClientError", e.Error()}
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		errorResponse := wire.ErrorsResponse{}
		json.Unmarshal(body, &errorResponse)
		e := errorResponse.Error
		return nil, ClientError{e.ErrorType, e.Detail}
	} else {
		responseMessage := new(ClientRegistration)
		e = json.Unmarshal(body, &responseMessage)
		if e != nil {
			return nil, ClientError{"ClientError", e.Error()}
		}
		return responseMessage, nil
	}
}