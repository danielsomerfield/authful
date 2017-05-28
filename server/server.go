package server

import (
	"net/http"
	"fmt"
	"github.com/danielsomerfield/authful/server/request"
	"encoding/base64"
	"net/url"
	"log"
	"encoding/json"
	"github.com/danielsomerfield/authful/server/wireTypes"
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

func authorizeHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	values := req.Form

	authorizationRequest, err := request.ParseAuthorizeRequest(values)
	if err != nil {
		//TODO: handle error
		return
	} else {
		client := getClient(authorizationRequest.ClientId)
		if client == nil {
			http.Error(w, formatError(err), http.StatusBadRequest)
			return;
		}
	}

	//Get the client
	//Reject if client doesn't exist
	//Reject if the redirect_uri doesn't match one configured with the client

	//Check scopes
	//Redirect to error if there is a scope in the request that's not in the client

	//Identify RO and ask for approval of request

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

func getClient(clientId string) *Client {
	return nil
}

type Client struct {
}

type Credentials struct {
	clientId     string
	clientSecret string
}

func formatError(error *request.ParseError) string {
	return fmt.Sprintf("The following fields are required: %s", error.MissingFields)
}

func (server *AuthServer) Stop() error {
	return nil
}

func (c Credentials) String() string {
	creds := fmt.Sprintf("{%s}:{%s}", url.QueryEscape(c.clientId), url.QueryEscape(c.clientSecret))
	return base64.StdEncoding.EncodeToString([]byte(creds))
}

func tokenHandler(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		log.Printf("Failed with following error: %+v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	//TODO:
	//Check for client_credentials type
	//Parse the request type
	//Check that all scopes are known
	//Create the token in the backend
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(wireTypes.TokenResponse{
		AccessToken: "TODO",
		TokenType: "Bearer",
		ExpiresIn: 3600,
	})
	writeOrError(w, bytes, err)
}



func writeOrError(w http.ResponseWriter, bytes []byte, err error) {
	if err == nil {
		w.Write(bytes)
	} else {
		log.Printf("Failed with following error: %+v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func (server *AuthServer) Start() *Credentials {
	log.Printf("Server started on port %v\n", server.port)
	http.HandleFunc("/token", tokenHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/authorize", authorizeHandler)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", server.port), nil)
		if err == nil {
			log.Printf("Failed to start up http server %s%n", err)
		}
	}()
	return &Credentials{
		clientId:     "CID", //TODO: real random id and secret and store client
		clientSecret: "Secret",
	}
}
