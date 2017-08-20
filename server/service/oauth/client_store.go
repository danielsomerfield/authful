package oauth

import (
	"fmt"
	"net/url"
	"encoding/base64"
	"github.com/danielsomerfield/authful/util"
)

type Client interface {
	CheckSecret(secret string) bool
	GetScopes() []string
	IsValidRedirectURI(uri string) bool
}

func (client DefaultClient) IsValidRedirectURI(uri string) bool {
	return false
}

func (c Credentials) String() string {
	creds := fmt.Sprintf("%s:%s", url.QueryEscape(c.ClientId), url.QueryEscape(c.ClientSecret))
	return base64.StdEncoding.EncodeToString([]byte(creds))
}

type Credentials struct {
	ClientId     string
	ClientSecret string
}

type ClientLookupFn func(clientId string) (Client, error)
type RegisterClientFn func(name string, scopes []string, urls [] string) (*Credentials, error)

type inMemoryClientStore struct {
	clients map[string]DefaultClient
}

func NewInMemoryClientStore() inMemoryClientStore {
	return inMemoryClientStore{
		clients: map[string]DefaultClient{},
	}
}

func (store inMemoryClientStore) LookupClient(clientId string) (Client, error) {
	return store.clients[clientId], nil
}

func (store inMemoryClientStore) RegisterClient(name string, scopes []string, redirectUris []string) (*Credentials, error) {
	clientId := util.GenerateRandomString(30)
	secret := util.GenerateRandomString(60) //TODO: replace with hash storage
	store.clients[clientId] = DefaultClient{
		name:     name,
		clientId: clientId,
		secret:   secret,
		scopes:   scopes,
		redirectUris: redirectUris,
	}
	return &Credentials{
		ClientId:     clientId,
		ClientSecret: secret,
	}, nil
}

type DefaultClient struct {
	name     string
	clientId string
	secret   string //TODO: replace this with a hash
	scopes   []string
	redirectUris [] string
}

func (client DefaultClient) GetScopes() []string {
	return client.scopes
}

func (client DefaultClient) CheckSecret(secret string) bool {
	return client.secret == secret
}
