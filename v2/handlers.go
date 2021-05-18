package v2

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	logrus.Error("No proper pin description in body")
	if pinDesc.Pin == nil || pinDesc.Direction == nil {
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
	params := mux.Vars(req)
	pin, err := strconv.Atoi(params["pin"])
	if err != nil || pin < 0 {
		logrus.Errorln("Invalid pin number")
		wr.WriteHeader(http.StatusNotFound)
		return
	}

	err = gh.ctrl.UnexportPin(pin)
	if err != nil {
		logrus.Errorf("Failed to unexport pin '%d': %v\n", pin, err.Error())
		if err == gpio.ErrNotExported {
			wr.WriteHeader(http.StatusBadRequest)
		} else {
			wr.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	wr.WriteHeader(http.StatusOK)
}

func (gh *gpioHandler) setPin(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusNotImplemented)
}

func (gh *gpioHandler) getPin(wr http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	pin, err := strconv.Atoi(params["pin"])
	if err != nil || pin < 0 {
		logrus.Errorln("Invalid pin number")
		wr.WriteHeader(http.StatusNotFound)
		return
	}

	val, err := gh.ctrl.GetValue(pin)
	if err != nil {
		logrus.Errorf("Failed to get pin '%d' value: %v\n", pin, err.Error())
		if err == gpio.ErrNotExported {
			wr.WriteHeader(http.StatusBadRequest)
		} else {
			wr.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	logrus.Tracef("Pin: %d Value: %d", pin, val)

	respData := struct {
		Pin   int `json:"pin"`
		Value int `json:"value"`
	}{
		Pin:   pin,
		Value: val,
	}
	buffer, err := json.Marshal(respData)
	if err != nil {
		logrus.Errorln("Failed to marshal result:", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wr.WriteHeader(http.StatusOK)
	wr.Write(buffer)
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
