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
	"crypto/x509"
	"crypto/tls"
)

type ClientRegistration struct {
	Data struct {
		ClientId     string  `json:"clientId,omitempty"`
		ClientSecret string  `json:"clientSecret,omitempty"`
	} `json:"data,omitempty"`
}

type ClientError struct {
	Type    string  `json:"errorType,omitempty"`
	Message string  `json:"errorMessage,omitempty"`
}

func (ce ClientError) Error() string {
	return fmt.Sprintf("%s : %s", ce.Type, ce.Message)
}

func (ce ClientError) ErrorType() string {
	return ce.Type
}

type APIClient struct {
	httpClient   *http.Client
	host         string
	clientId     string
	clientSecret string
}

func NewAPIClient(host string, clientId string, clientSecret string, caCertFiles []string) (*APIClient, error) {
	caCertPool := x509.NewCertPool()
	for _, file := range caCertFiles {
		caCert, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("Failed to read CA cerficiate at %s", file)
		}
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("Failed to add CA cerficiate at %s", file)
		}
	}

	return &APIClient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			}},
		host:         host,
		clientId:     clientId,
		clientSecret: clientSecret,
	}, nil
}

func (apiClient *APIClient) RegisterClient(clientName string) (*ClientRegistration, error) {
	credentials := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s:%s", url.QueryEscape(apiClient.clientId), url.QueryEscape(apiClient.clientSecret))))

	createClientRequest := map[string]interface{}{
		"command": map[string]string{
			"name": clientName,
		},
	}
	messageBytes, _ := json.Marshal(createClientRequest)

	post, _ := http.NewRequest("POST", fmt.Sprintf("%s/admin/clients", apiClient.host),
		bytes.NewReader(messageBytes))
	post.Header.Set("Content-Type", "application/json")
	post.Header.Set("Authorization", "Basic "+credentials)

	response, e := apiClient.httpClient.Do(post)
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

type UserRegistration struct{}

func (apiClient *APIClient) RegisterUser(username string, password string, authMethods []string) (*UserRegistration, error) {
	credentials := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s:%s", url.QueryEscape(apiClient.clientId), url.QueryEscape(apiClient.clientSecret))))

	createClientRequest := map[string]interface{}{
		"command": map[string]interface{}{
			"username":    username,
			"password":    password,
			"authMethods": authMethods,
		},
	}
	messageBytes, _ := json.Marshal(createClientRequest)

	post, _ := http.NewRequest("POST", fmt.Sprintf("%s/admin/users", apiClient.host),
		bytes.NewReader(messageBytes))
	post.Header.Set("Content-Type", "application/json")
	post.Header.Set("Authorization", "Basic "+credentials)

	response, e := apiClient.httpClient.Do(post)
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
		responseMessage := new(UserRegistration)
		e = json.Unmarshal(body, &responseMessage)
		if e != nil {
			return nil, ClientError{"ClientError", e.Error()}
		}
		return responseMessage, nil
	}
}
