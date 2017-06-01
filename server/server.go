package server

import (
	"net/http"
	"fmt"
	"encoding/base64"
	"net/url"
	"log"
	"github.com/danielsomerfield/authful/server/oauth"
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
	clientId     string
	clientSecret string
}

func (c Credentials) String() string {
	creds := fmt.Sprintf("{%s}:{%s}", url.QueryEscape(c.clientId), url.QueryEscape(c.clientSecret))
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
	return &Credentials{
		clientId:     "CID", //TODO: real random id and secret and store client
		clientSecret: "Secret",
	}
}

func (server *AuthServer) Stop() error {
	return server.httpServer.Shutdown(nil)
}

func defaultTokenGenerator() string {
	return "TODO"
}

var tokenHandlerConfig = oauth.TokenHandlerConfig {
	DefaultTokenExpiration: 3600,
}

func init() {
	http.HandleFunc("/token", oauth.NewTokenHandler(tokenHandlerConfig, oauth.DefaultClientLookup, defaultTokenGenerator))
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/authorize", oauth.AuthorizeHandler)
}
