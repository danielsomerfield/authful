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

type AuthorizeResponse struct {
	MissingFields []string
}

func authorizeHandler(w http.ResponseWriter, req *http.Request) {
	//req.ParseForm()
	//form := req.Form


	//request_type := form.Get("request_type")
	//if request_type == "" {
	//	return AuthorizeResponse{
	//		MissingFields:[]string{"request_type"},
	//	}
	//}

	//client_id 	REQUIRED
	//redirect_uri 	OPTIONAL
	//scope		OPTIONAL
	//state		OPTIONAL
}

func (server *AuthServer) Stop() error {
	return nil
}

func (server *AuthServer) Start() error {
	fmt.Printf("Server started on port %v\n", server.port)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/authorize", authorizeHandler)
	http.ListenAndServe(fmt.Sprintf(":%v", server.port), nil)
	return nil
}
