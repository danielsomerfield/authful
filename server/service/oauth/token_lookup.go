package oauth

import "time"

type TokenMetaData struct {
	Token      string
	Expiration time.Time
	ClientId   string
}