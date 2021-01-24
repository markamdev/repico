package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/markamdev/repico/gpio"
)

const (
	// expected request URI is in format: /v1/gpio/<number>
	apiPrefix       = "/v1/gpio/"
	listURIPrefix   = "/v1/gpio/list"
	aliasURIPrefix  = "/v1/gpio/alias"
	numberURIPrefix = "/v1/gpio/number"
	configURIPrefix = "/v1/gpio/config"
)

type myHandler struct {
	handlers map[string]func(resp http.ResponseWriter, req *http.Request)
}

// ServeHTTP method for serving HTTP requests
func (m myHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// TODO remove log below when server fully implemented
	log.Printf("Request received. Type: %+v, URI: %v", req.Method, req.RequestURI)

	// check if URI is supported
	for uri, handler := range m.handlers {
		if strings.HasPrefix(req.RequestURI, uri) {
			log.Println("URI match:", uri)
			handler(resp, req)
			return
		}
	}
	log.Println("Invalid URI:", req.RequestURI)
	resp.WriteHeader(http.StatusNotFound)
}

func createHandler() myHandler {
	result := myHandler{}
	result.handlers = make(map[string]func(resp http.ResponseWriter, req *http.Request))
	// new approach and URIs
	result.handlers[listURIPrefix] = notSupported
	result.handlers[aliasURIPrefix] = notSupported
	result.handlers[numberURIPrefix] = setGPIO
	result.handlers[configURIPrefix] = notSupported

	return result
}

func notSupported(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusNotImplemented)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte("{ \"message\" : \"Function not yet implemented\"}"))
}

func setGPIO(resp http.ResponseWriter, req *http.Request) {
	// check HTTP method first
	if req.Method != http.MethodPut {
		// currently not supported
		resp.WriteHeader(http.StatusNotImplemented)
		return
	}
	// parse URI to check pin number
	pinString := strings.TrimPrefix(req.RequestURI, numberURIPrefix+"/")
	val, err := strconv.ParseInt(pinString, 10, 32)
	if err != nil {
		log.Println("Invalid GPIO number", pinString)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	// 1kB of body data should be enough
	buffer := make([]byte, 1024)
	len, err := req.Body.Read(buffer)
	// EOF does not mean - only end of data reached)
	if err != nil && err != io.EOF {
		log.Println("Error while reading request body:", err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	// empty buffer -> no need to go further
	if len == 0 {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	dt := ReqData{}
	err = json.Unmarshal(buffer[:len], &dt)
	if err != nil {
		log.Println("Failed to unmarshal data:", buffer[:len], string(buffer[:len]), "Error:", err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = gpio.SetGPIO(int(val), dt.State)
	if err != nil {
		log.Println("Setting GPIO value failed:", err.Error())
		// TODO FIXME not all errors here are InteralError - some can be BadRequest
		resp.WriteHeader(http.StatusInternalServerError)
	} else {
		resp.WriteHeader(http.StatusOK)
	}
}

/*
func (m myHandler) setGPIOValue(pin, value int) error {
	log.Printf("Setting pin %v to value %v", pin, value)
	return gpio.SetGPIO(pin, value)
}
*/
