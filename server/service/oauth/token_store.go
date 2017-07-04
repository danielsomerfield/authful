package oauth

import (
	"time"
)

type GetTokenMetaDataFn func(token string) (*TokenMetaData, error)
type StoreTokenMetaDataFn func(token string, tokenMetaData TokenMetaData) error

type TokenMetaData struct {
	Token      string
	Expiration time.Time
	ClientId   string
}

type inMemoryTokenStore struct {
	tokenMetaData map[string]TokenMetaData
}

func NewInMemoryTokenStore() inMemoryTokenStore {
	return inMemoryTokenStore{
		tokenMetaData: map[string]TokenMetaData{},
	}
}

func (store *inMemoryTokenStore) StoreToken(token string, tokenMetaData TokenMetaData) error {
	store.tokenMetaData[token] = tokenMetaData
	return nil
}

func (store *inMemoryTokenStore) GetToken(token string) (*TokenMetaData, error) {
	tokenMetaData := store.tokenMetaData[token]
	return &tokenMetaData, nil
}
