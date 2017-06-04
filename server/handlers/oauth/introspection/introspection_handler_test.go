package introspection

import (
	"net/http"
	"testing"
	"fmt"
	"strings"
)

func init() {
	http.HandleFunc("/introspect", NewIntrospectionHandler())
	go http.ListenAndServe(":8080", nil)
}

var activeToken = "active_token"
var validClientId = "valid-client-id"

func TestIntrospectionHandler_ApprovesValidToken(t *testing.T) {
	post, _ := http.NewRequest("POST", "http://localhost:8080/introspect",
		strings.NewReader(fmt.Sprintf("token=%s", activeToken)))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.Header.Set("Authorization", "Bearer " + validClientId)
	//http.DefaultClient.Do(post)
}
