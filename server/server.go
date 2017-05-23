package server

import (
	"net/http"
	"fmt"
)

type AuthServer struct {
	port int
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

func (server *AuthServer) Stop() error {
	return nil
}

func (server *AuthServer) Start() error {
	fmt.Printf("Server started on port %v\n", server.port)
	http.HandleFunc("/health", healthHandler)
	http.ListenAndServe(fmt.Sprintf(":%v", server.port), nil)
	return nil
}
