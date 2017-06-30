package oauth

import (
	"time"
)

type GetTokenMetaDataFn func(token string) *TokenMetaData
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

func (*inMemoryTokenStore) StoreToken(token string, tokenMetaData TokenMetaData) error {
	return nil
}

func (*inMemoryTokenStore) GetToken(token string) *TokenMetaData {
	return nil
}
