package v2

import (
	"github.com/gorilla/mux"
	"github.com/markamdev/repico/gpio"
	"github.com/sirupsen/logrus"
)

func AttachHandlers(handler *mux.Router, controller gpio.Controller) {
	logrus.Traceln("v2.AttachHandlers()")

	hndlr := gpioHandler{ctrl: controller}

	handler.HandleFunc("/gpio", hndlr.addPin).Methods("POST")
	handler.HandleFunc("/gpio", hndlr.getAllPins).Methods("GET")

	handler.HandleFunc("/gpio/{pin:[0-9]+}", hndlr.deletePin).Methods("DELETE")
	handler.HandleFunc("/gpio/{pin:[0-9]+}", hndlr.setPin).Methods("PATCH")
	handler.HandleFunc("/gpio/{pin:[0-9]+}", hndlr.getPin).Methods("GET")
}
