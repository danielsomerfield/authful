package oauth

import "errors"

type ClientLookupFn func(clientId string, clientSecret string) (*Client, error)

func DefaultClientLookup(clientId string, clientSecret string) (*Client, error) {
	return nil, errors.New("NYI")
}
