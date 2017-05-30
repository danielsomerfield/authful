package oauth

import "errors"

type ClientLookupFn func(clientId string) (Client, error)

func DefaultClientLookup(clientId string) (Client, error) {
	return nil, errors.New("NYI")
}

type DefaultClient struct {

}

func (client *DefaultClient) checkSecret(secret string) bool {
	return false
}