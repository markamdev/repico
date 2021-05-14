package v2

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markamdev/repico/gpio"
)

func CreateHandler(controller gpio.Controller) http.Handler {

	router := mux.NewRouter()
	router.HandleFunc("/v2/gpio/{pin:[0-9]+}", defaultStub).Methods("GET", "DELETE", "PATCH")
	router.HandleFunc("/v2/gpio", defaultStub).Methods("GET", "POST")

	return &gpioRouter{ctrl: controller, handler: router}
}

type gpioRouter struct {
	ctrl    gpio.Controller
	handler *mux.Router
}

func (gr *gpioRouter) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	gr.handler.ServeHTTP(wr, req)
}

const (
	errorMessage = "{ \"error\":\"endpoint %s for method %s not implemented\" }"
)

func defaultStub(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusNotImplemented)
	msg := fmt.Sprintf(errorMessage, req.RequestURI, req.Method)
	wr.Write([]byte(msg))
}
