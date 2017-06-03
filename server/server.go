package server

import (
	"net/http"
	"fmt"
	"encoding/base64"
	"net/url"
	"log"
	"github.com/danielsomerfield/authful/server/oauth"
	"github.com/danielsomerfield/authful/server/oauth/handlers"
	"time"
)

type AuthServer struct {
	port       int
	httpServer http.Server
}

func NewAuthServer() *AuthServer {
	return &AuthServer{
		port: 8080,
	}
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\":\"ok\"}"))
}

type Credentials struct {
	ClientId     string
	ClientSecret string
}

func (c Credentials) String() string {
	creds := fmt.Sprintf("{%s}:{%s}", url.QueryEscape(c.ClientId), url.QueryEscape(c.ClientSecret))
	return base64.StdEncoding.EncodeToString([]byte(creds))
}

func (server *AuthServer) Start() *Credentials {
	go func() {
		httpServer := http.Server{Addr: fmt.Sprintf(":%v", server.port)}
		err := httpServer.ListenAndServe()
		if err == nil {
			log.Printf("Failed to start up http server %s%n", err)
		} else {
			log.Printf("Server started on port %v\n", server.port)
		}
	}()
	clientId := "CID"
	clientSecret := "Secret"
	credentials := Credentials{
		ClientId:     clientId, //TODO: real random id and secret and store client
		ClientSecret: clientSecret,
	}
	return &credentials
}

func (server *AuthServer) Stop() error {
	return server.httpServer.Shutdown(nil)
}

func defaultTokenGenerator() string {
	return "TODO"
}

var tokenHandlerConfig = handlers.TokenHandlerConfig {
	DefaultTokenExpiration: 3600,
}

type DefaultTokenStore struct {

}

func (tokenStore DefaultTokenStore) StoreToken(token string, tokenMetaData handlers.TokenMetaData) {

}

var tokenStore = DefaultTokenStore{}

func currentTimeFn() time.Time {
	return time.Now()
}

func init() {
	http.HandleFunc("/token", handlers.NewTokenHandler(
		tokenHandlerConfig,
		oauth.DefaultClientLookup,
		defaultTokenGenerator,
		tokenStore,
		currentTimeFn))
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/authorize", oauth.AuthorizeHandler)
}
