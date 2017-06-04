package introspection

import (
	"net/http"

)

func NewIntrospectionHandler() func(http.ResponseWriter, *http.Request) {
	return func(http.ResponseWriter, *http.Request) {

	}
}
