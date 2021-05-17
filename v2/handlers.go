package v2

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/markamdev/repico/gpio"
	"github.com/sirupsen/logrus"
)

type gpioHandler struct {
	ctrl gpio.Controller
}

func (gh *gpioHandler) addPin(wr http.ResponseWriter, req *http.Request) {
	buffer := make([]byte, 1024)
	n, err := req.Body.Read(buffer)
	if err != nil && err != io.EOF {
		logrus.Errorln("Failed to read body")
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}
	if n == 0 {
		logrus.Warnln("Empty request body")
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	var pinDesc struct {
		Pin       *int    `json:"pin"`
		Direction *string `json:"direction"`
	}
	err = json.Unmarshal(buffer[:n], &pinDesc)
	if err != nil {
		logrus.Error("Unable to unmarshal request data:", err)
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	if pinDesc.Pin == nil || pinDesc.Direction == nil {
		logrus.Error("No proper pin description in body")
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	err = gh.ctrl.ExportPin(*pinDesc.Pin, gpio.StringToDirection(*pinDesc.Direction))
	if err == nil {
		wr.WriteHeader(http.StatusOK)
		return
	}

	if err == gpio.ErrUnknown {
		logrus.Error("GPIO pin exporting failed")
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	logrus.Error("GPIO pin exporting error:", err)
	wr.WriteHeader(http.StatusBadRequest)
	wr.Write([]byte(err.Error()))

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
		logrus.Errorln("Error when listing GPIO pins:", err)
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
		logrus.Errorln("Error when marshalling pin data:", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.Write(buffer)
	wr.WriteHeader(http.StatusOK)
}
