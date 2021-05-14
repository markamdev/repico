package v2

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markamdev/repico/gpio"
	"github.com/sirupsen/logrus"
)

func CreateHandler(controller gpio.Controller) http.Handler {
	logrus.Traceln("v2.CreateHandler()")

	hndlr := gpioHandler{ctrl: controller}

	router := mux.NewRouter()
	router.HandleFunc("/v2/gpio/{pin:[0-9]+}", hndlr.deletePin).Methods("DELETE")
	router.HandleFunc("/v2/gpio/{pin:[0-9]+}", hndlr.setPin).Methods("PATCH")
	router.HandleFunc("/v2/gpio/{pin:[0-9]+}", hndlr.getPin).Methods("GET")
	router.HandleFunc("/v2/gpio", hndlr.addPin).Methods("POST")
	router.HandleFunc("/v2/gpio", hndlr.getAllPins).Methods("GET")

	return router
}
