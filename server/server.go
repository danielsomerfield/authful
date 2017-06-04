package server

import (
	"net/http"
	"fmt"
	"log"
	oauth_handlers "github.com/danielsomerfield/authful/server/handlers/oauth"
	oauth_service "github.com/danielsomerfield/authful/server/service/oauth"
	"time"
	"github.com/danielsomerfield/authful/util"
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

func (server *AuthServer) Start() (*oauth_service.Credentials, error) {
	go func() {
		httpServer := http.Server{Addr: fmt.Sprintf(":%v", server.port)}
		err := httpServer.ListenAndServe()
		if err == nil {
			log.Printf("Failed to start up http server %s%n", err)
		} else {
			log.Printf("Server started on port %v\n", server.port)
		}
	}()

	//TODO: make generation of startup credentials a configuration option

	return clientStore.RegisterClient("Root Admin", []string{"administrate"})

}

func (server *AuthServer) Stop() error {
	return server.httpServer.Shutdown(nil)
}

func defaultTokenGenerator() string {
	return util.GenerateRandomString(25)
}

var tokenHandlerConfig = oauth_handlers.TokenHandlerConfig{
	DefaultTokenExpiration: 3600,
}

type DefaultTokenStore struct {
}

func (tokenStore DefaultTokenStore) StoreToken(token string, tokenMetaData oauth_handlers.TokenMetaData) error {
	return nil
}

var tokenStore = DefaultTokenStore{}
var clientStore = oauth_service.NewInMemoryClientStore()

func currentTimeFn() time.Time {
	return time.Now()
}

func init() {
	http.HandleFunc("/token", oauth_handlers.NewTokenHandler(
		tokenHandlerConfig,
		clientStore.LookupClient,
		defaultTokenGenerator,
		tokenStore.StoreToken,
		currentTimeFn))
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/authorize", oauth_handlers.AuthorizeHandler)
}
