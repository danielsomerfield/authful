package oauth

import "errors"

type ClientStore interface {
	LookupClient(clientId string) (Client, error)
	RegisterClient(clientId string, client Client) error
}

type InMemoryClientStore struct {

}

func (InMemoryClientStore) LookupClient(clientId string) (Client, error) {
	return nil, errors.New("NYI")
}

func (InMemoryClientStore) RegisterClient(clientId string, client Client) error {
	return errors.New("NYI")
}

type DefaultClient struct {

}

func (client *DefaultClient) checkSecret(secret string) bool {
	return false
}