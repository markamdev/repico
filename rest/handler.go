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
	pinString := strings.TrimPrefix(req.RequestURI, apiPrefix)
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
	err = m.setGPIOValue(int(val), dt.State)
	if err != nil {
		log.Println("Setting GPIO value failed:", err.Error())
		// TODO FIXME not all errors here are InteralError - some can be BadRequest
		resp.WriteHeader(http.StatusInternalServerError)
	} else {
		resp.WriteHeader(http.StatusOK)
	}
}

func (m myHandler) setGPIOValue(pin, value int) error {
	log.Printf("Setting pin %v to value %v", pin, value)
	return gpio.SetGPIO(pin, value)
}
