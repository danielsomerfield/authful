package oauth

import (
	"fmt"
	"net/url"
	"encoding/base64"
	util "github.com/danielsomerfield/authful/common/util"
)

type Client interface {
	CheckSecret(secret string) bool
	GetScopes() []string
	IsValidRedirectURI(uri string) bool
	GetDefaultRedirectURI() string
}

//TODO: implement these two functions
func (client DefaultClient) IsValidRedirectURI(uri string) bool {
	return util.Contains(client.redirectUris, uri)
}

func (client DefaultClient) GetDefaultRedirectURI() string {
	return ""
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
type RegisterClientFn func(name string, scopes []string, urls [] string, defaultRedirectURI string) (*Credentials, error)

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

func (store inMemoryClientStore) RegisterClient(name string, scopes []string, redirectUris []string,
	defaultRedirectURI string) (*Credentials, error) {
	clientId := util.GenerateRandomString(30)
	secret := util.GenerateRandomString(60) //TODO: replace with hash storage
	store.clients[clientId] = DefaultClient{
		name:               name,
		clientId:           clientId,
		secret:             secret,
		scopes:             scopes,
		redirectUris:       redirectUris,
		defaultRedirectURI: defaultRedirectURI,
	}
	return &Credentials{
		ClientId:     clientId,
		ClientSecret: secret,
	}, nil
}

type DefaultClient struct {
	name               string
	clientId           string
	secret             string //TODO: replace this with a hash
	scopes             []string
	redirectUris       [] string
	defaultRedirectURI string
}

func (client DefaultClient) GetScopes() []string {
	return client.scopes
}

func (client DefaultClient) CheckSecret(secret string) bool {
	return client.secret == secret
}
