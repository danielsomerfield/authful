package server

import (
	"net/http"
	"fmt"
	"log"
	oauth_service "github.com/danielsomerfield/authful/server/service/oauth"
	"time"
	"github.com/danielsomerfield/authful/util"
	"github.com/danielsomerfield/authful/server/handlers/oauth/token"
	"github.com/danielsomerfield/authful/server/handlers/oauth/authorization"
	"github.com/danielsomerfield/authful/server/handlers/oauth/introspection"
)

type AuthServer struct {
	port       int
	httpServer http.Server
}

func NewAuthServer(port int) *AuthServer {
	return &AuthServer{
		port: port,
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
			log.Fatalf("Failed to start up http server %s\n", err)
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

var tokenHandlerConfig = token.TokenHandlerConfig{
	DefaultTokenExpiration: 3600,
}

//Injected Service Dependencies
var tokenStore = oauth_service.NewInMemoryTokenStore()
var clientStore = oauth_service.NewInMemoryClientStore()
var clientAccessControlFn = func(request http.Request) bool {
	//TODO: This will need to support two auth methods: client credentials and token
	//TODO: Implement client credentials first (token can come later)
	//Get the credentials from the request
	//Look up the client
	//Make sure the client has the introspect_token or administrate scope
	return true //TODO: NYI
}

func currentTimeFn() time.Time {
	return time.Now()
}

func init() {
	http.HandleFunc("/token", token.NewTokenHandler(
		tokenHandlerConfig,
		clientStore.LookupClient,
		defaultTokenGenerator,
		tokenStore.StoreToken,
		currentTimeFn))
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/authorize", authorization.AuthorizeHandler)
	http.HandleFunc("/introspect", introspection.NewIntrospectionHandler(clientAccessControlFn, tokenStore.GetToken))
}
