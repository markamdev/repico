package v2

import (
	"encoding/json"
	"log"
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
	pins, err := gh.ctrl.ListExportedPins()
	if err != nil {
		log.Println("Error when listing GPIO pins:", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(pins) == 0 {
		wr.WriteHeader(http.StatusNoContent)
		return
	}

	type pinData struct {
		Pin       string
		Direction string
	}
	result := make([]pinData, 0, len(pins))
	for k, v := range pins {
		result = append(result, pinData{Pin: k, Direction: gpio.DirectionToString(v)})
	}

	buffer, err := json.Marshal(result)
	if err != nil {
		log.Println("Error when marshalling pin data:", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.Write(buffer)
	wr.WriteHeader(http.StatusOK)
}
