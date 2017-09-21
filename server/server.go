package server

import (
	"net/http"
	"fmt"
	"log"
	oauthsvc "github.com/danielsomerfield/authful/server/service/oauth"
	"time"
	"github.com/danielsomerfield/authful/util"
	"github.com/danielsomerfield/authful/server/handlers/oauth/token"
	"github.com/danielsomerfield/authful/server/handlers/oauth/authorization"
	"github.com/danielsomerfield/authful/server/handlers/oauth/introspection"
	"github.com/danielsomerfield/authful/server/service/accesscontrol"
	"github.com/danielsomerfield/authful/server/handlers/admin/client"
	adminuser "github.com/danielsomerfield/authful/server/handlers/admin/user"
	usersvc "github.com/danielsomerfield/authful/server/service/admin/user"
	user2 "github.com/danielsomerfield/authful/server/repository/user"
	"github.com/danielsomerfield/authful/server/service/crypto"
	"github.com/danielsomerfield/authful/server/wire/oauth"
	"net/url"
)

type AuthServer struct {
	port        int
	httpServer  http.Server
	running     chan bool
	tlsCertFile string
	tlsKeyFile  string
}

type Config struct {
	Port int
	TLS TLSConfig
}

type TLSConfig struct {
	TLSCertFile string     `yaml:"cert"`
	TLSKeyFile  string     `yaml:"key"`
}

func NewAuthServer(config *Config) *AuthServer {
	return &AuthServer{
		port:        config.Port,
		tlsCertFile: config.TLS.TLSCertFile,
		tlsKeyFile:  config.TLS.TLSKeyFile,
	}
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\":\"ok\"}"))
}

func (server *AuthServer) Start() (*oauthsvc.Credentials, error) {
	log.Printf("Starting server up http server on port %d\n", server.port)

	go func() {
		httpServer := http.Server{Addr: fmt.Sprintf(":%v", server.port)}
		err := httpServer.ListenAndServeTLS(server.tlsCertFile, server.tlsKeyFile)
		if err != nil {
			log.Fatalf("Failed to start up http server %s\n", err)
		} else {
			log.Printf("Server started on port %v\n", server.port)
		}
		server.running <- false
	}()

	server.running = make(chan bool)

	//TODO: make generation of startup credentials a configuration option

	return clientStore.RegisterClient("Root Admin", []string{"administrate", "introspect"}, nil, "")

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

func defaultCodeGenerator() string {
	return util.GenerateRandomString(6)
}

var tokenHandlerConfig = token.TokenHandlerConfig{
	DefaultTokenExpiration: 3600,
}

//Injected Service Dependencies
var tokenStore = oauthsvc.NewInMemoryTokenStore()
var clientStore = oauthsvc.NewInMemoryClientStore()
var userRepository = user2.NewInMemoryUserRepository()

var approvalRequestStore = func(request *oauth.AuthorizeRequest) string {
	return util.GenerateRandomString(6)
}

var approvalLookup = func(approvalType string, requestId string) *url.URL {
	u, _ := url.Parse("/login")
	return u
}

var defaultErrorRenderer = func(code string) []byte { //TODO: build full template renderer
	return []byte(fmt.Sprintf("<html>%s</html>", code))
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
	http.HandleFunc("/authorize", authorization.NewAuthorizationHandler(clientStore.LookupClient,
		defaultErrorRenderer, approvalRequestStore, approvalLookup))
	http.HandleFunc("/introspect", introspection.NewIntrospectionHandler(
		accesscontrol.NewClientAccessControlFnWithScopes(clientStore.LookupClient, "introspect"), tokenStore.GetToken))
	http.HandleFunc("/admin/clients", client.NewProtectedHandler(
		clientStore.RegisterClient, clientStore.LookupClient))
	http.HandleFunc("/admin/users", adminuser.NewProtectedHandler(
		usersvc.NewRegisterUserFn(userRepository.SaveUser, crypto.ScryptHash), clientStore.LookupClient))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

	})

}
