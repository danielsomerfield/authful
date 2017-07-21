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
	"github.com/danielsomerfield/authful/server/service/accesscontrol"
	"github.com/danielsomerfield/authful/server/handlers/admin"
)

type AuthServer struct {
	port       int
	httpServer http.Server
	running chan bool
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
	log.Printf("Starting server up http server on port %d\n", server.port)

	go func() {
		httpServer := http.Server{Addr: fmt.Sprintf(":%v", server.port)}
		err := httpServer.ListenAndServeTLS("../resources/test/server.crt", "../resources/test/server.key")
		if err != nil {
			log.Fatalf("Failed to start up http server %s\n", err)
		} else {
			log.Printf("Server started on port %v\n", server.port)
		}
		server.running <- false
	}()

	server.running = make(chan bool)

	//TODO: make generation of startup credentials a configuration option

	return clientStore.RegisterClient("Root Admin", []string{"administrate", "introspect"})

}

func (server *AuthServer) Wait() {
	<-server.running
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
	http.HandleFunc("/introspect", introspection.NewIntrospectionHandler(
		accesscontrol.NewClientAccessControlFnWithScopes(clientStore.LookupClient, "introspect"), tokenStore.GetToken))
	http.HandleFunc("/admin/clients", admin.NewRegisterClientHandler(
		accesscontrol.NewClientAccessControlFnWithScopes(clientStore.LookupClient, "administrate"), clientStore.RegisterClient))
}
