package v2

import (
	"net/http"

	"github.com/markamdev/repico/gpio"
)

type gpioHandler struct {
	ctrl gpio.Controller
}

func (gh *gpioHandler) addPin(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusNotImplemented)
}

func (gh *gpioHandler) deletePin(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusNotImplemented)
}

func (gh *gpioHandler) setPin(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusNotImplemented)
}

func (gh *gpioHandler) getPin(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusNotImplemented)
}

func (gh *gpioHandler) getAllPins(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusNotImplemented)
}
