package rest

import (
	"log"
	"net/http"
	"strings"
)

const (
	// expected request URI is in format: /v1/gpio/<number>
	apiPrefix = "/v1/gpio/"
)

type myHandler struct{}

// ServeHTTP method for serving HTTP requests
func (m myHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// TODO remove log below when server fully implemented
	log.Printf("Request received. Type: %+v, URI: %v", req.Method, req.RequestURI)
	// skip further processing on incorrect URI
	if !strings.HasPrefix(req.RequestURI, apiPrefix) {
		log.Println("Invalid URI:", req.RequestURI)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	switch req.Method {
	case http.MethodGet:
		// currently not supported
		resp.WriteHeader(http.StatusNotImplemented)
	case http.MethodPut:
		m.setGPIO(resp, req)
	default:
		log.Println("Unsupported method:", req.Method)
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (m myHandler) setGPIO(resp http.ResponseWriter, req *http.Request) {

}
