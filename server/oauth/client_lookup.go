package oauth

import (
	"fmt"
	"net/url"
	"encoding/base64"
	"math/rand"
)

type Client interface {
	CheckSecret(secret string) bool
	GetScopes() []string
}

func (c Credentials) String() string {
	creds := fmt.Sprintf("{%s}:{%s}", url.QueryEscape(c.ClientId), url.QueryEscape(c.ClientSecret))
	return base64.StdEncoding.EncodeToString([]byte(creds))
}

type Credentials struct {
	ClientId     string
	ClientSecret string
}

type ClientLookupFn func(clientId string) (Client, error)

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

func (store inMemoryClientStore) RegisterClient(name string, scopes []string) (*Credentials, error) {
	clientId := generateRandomString(30)
	secret := generateRandomString(60) //TODO: replace with hash storage
	store.clients[clientId] = DefaultClient{
		name:     name,
		clientId: clientId,
		secret:   secret,
		scopes:   scopes,
	}
	return &Credentials{
		ClientId:     clientId,
		ClientSecret: secret,
	}, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type DefaultClient struct {
	name     string
	clientId string
	secret   string //TODO: replace this with a hash
	scopes   []string
}

func (client DefaultClient) GetScopes() []string {
	return client.scopes
}

func (client DefaultClient) CheckSecret(secret string) bool {
	return client.secret == secret
}
