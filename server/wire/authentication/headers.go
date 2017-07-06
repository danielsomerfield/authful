package authentication

import (
	"regexp"
	"net/http"
	"encoding/base64"
)

type ClientCredentials struct {
	ClientId     string
	ClientSecret string
}

func ParseClientCredentialsBasicHeader(httpRequest http.Request) (*ClientCredentials, error) {
	authHeader := httpRequest.Header.Get("Authorization")
	if authHeader != "" {
		encodedCredentials := regexp.MustCompile("Basic ([a-zA-Z0-9]*)").FindStringSubmatch(string(authHeader))
		if len(encodedCredentials) > 1 {
			credentialBytes, err := base64.RawStdEncoding.DecodeString(encodedCredentials[1])
			if err != nil {
				return nil, err
			}
			clientCredentialsString := string(credentialBytes)
			creds := regexp.MustCompile("(.*):(.*)").FindStringSubmatch(clientCredentialsString)
			if len(creds) > 1 {
				clientCredentials := ClientCredentials{
					ClientId: creds[1],
				}

				if len(creds) > 2 {
					clientCredentials.ClientSecret = creds[2]
				}
				return &clientCredentials, nil
			}
		}

	}
	return nil, nil
}
