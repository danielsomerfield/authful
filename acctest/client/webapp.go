package main

import (
	"github.com/danielsomerfield/authful/client/admin"
	"log"
	"net/http"
	"fmt"
	"math/rand"
)

func main() {
	clientRegistration := register()

	mux := http.NewServeMux()
	mux.HandleFunc("/show", func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("token")
		if cookie == nil {
			//Redirect to the auth server
			state := string(rand.Int())
			http.Redirect(w, r, fmt.Sprintf(
				"https://auth-server/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
				clientRegistration.Data.ClientId, "https://webapp/callback", "read", state), 302)
		} else {
			//Got the cookie, now make the request
			invalidRequest, _ := http.NewRequest("GET", "https://resource-server:8080/", nil)
			invalidRequest.Header.Set("Authorization", "Bearer 12345")
			//resp, err = CreateHttpsClient().Do(invalidRequest)
		}
	})

}

func CreateHttpsClient() *http.Client {
	return nil
}

func register() *admin.ClientRegistration {
	clientId := "TODO"
	clientSecret := "TODO"
	caCerts := []string{
		"/cacerts.crt",
	}
	apiClient, err := admin.NewAPIClient("https://auth-server:8080", clientId, clientSecret, caCerts)
	if err != nil {
		log.Fatalf("Failed to get api client")
	}
	client, err := apiClient.RegisterClient("webapp",
		[]string{},
		[]string{},
		"")

	if err != nil {
		log.Fatalf("Failed to get api client")
	}

	return client
}
