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
	result.handlers[numberURIPrefix] = handleNumber
	result.handlers[configURIPrefix] = notSupported

	return result
}

func notSupported(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusNotImplemented)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte("{ \"message\" : \"Function not yet implemented\"}"))
}

func handleNumber(resp http.ResponseWriter, req *http.Request) {
	// this function should support
	// - GET and PUT for single number
	// - GET and PATCH for multiple numbers

	if req.RequestURI == numberURIPrefix {
		// multi-number case
		handleMultiNumber(resp, req)
	} else {
		handleSingleNumber(resp, req)
	}
}

func handleSingleNumber(resp http.ResponseWriter, req *http.Request) {
	// check HTTP method first
	if req.Method != http.MethodPut && req.Method != http.MethodGet {
		// for single number URI only GET and PUT are allowed
		log.Println("Invalid method:", req.Method)
		resp.WriteHeader(http.StatusMethodNotAllowed)
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

	var execError error
	// handle PUT
	if req.Method == http.MethodPut {
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
		execError = gpio.SetGPIO(int(val), dt.State)
		if execError == nil {
			resp.WriteHeader(http.StatusOK)
			return
		}
	} else {
		var pinVal int
		pinVal, execError = gpio.GetGPIO(int(val))
		if execError == nil {
			result := PinNumberData{Number: int(val), State: pinVal}
			content, err := json.Marshal(result)
			if err == nil {
				resp.WriteHeader(http.StatusOK)
				resp.Header().Set("Content-Type", "application/json")
				resp.Write(content)
				return
			}

			log.Println("Failed to marshal GET result")
			execError = err
		}
	}

	if execError != nil {
		log.Println("GPIO operation failed:", execError.Error())
		// TODO FIXME not all errors here are InteralError - some can be BadRequest
		// - add new error type to rest/types.go
		// - convert recevied error into some HTTP message here
		resp.WriteHeader(http.StatusInternalServerError)
	}

}

func handleMultiNumber(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPatch && req.Method != http.MethodGet {
		// only GET and PATCH are supported
		log.Println("Invalid method:", req.Method)
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if req.Method == http.MethodPatch {
		// read configuration from BODY
		buffer := make([]byte, 2048)
		len, err := req.Body.Read(buffer)
		if err != nil && err != io.EOF {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len == 0 {
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
		pinConfig := make([]PinNumberData, 0)
		err = json.Unmarshal(buffer[:len], &pinConfig)
		if err != nil {
			log.Println("Unmarshal error:", err.Error())
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte(err.Error()))
			return
		}
		log.Println("Unmarshalled request data:", pinConfig)
		for _, item := range pinConfig {
			log.Println("- setting", item.Number, "to", item.State)
			err = gpio.SetGPIO(item.Number, item.State)
			if err != nil {
				log.Println("Error while setting pin in handleMultiNumber:", err.Error())
				resp.WriteHeader(http.StatusInternalServerError)
				// TODO add info about failed setting
				return
			}
		}
		resp.WriteHeader(http.StatusOK)
		return
	}

	resp.WriteHeader(http.StatusNotImplemented)
}
