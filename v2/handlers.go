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
	logrus.Debugln("addPin() handler")
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

	var pinDesc pinConfigPointer
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

	switch err {
	case nil:
		wr.WriteHeader(http.StatusOK)
	case gpio.ErrNotImplemented:
		wr.WriteHeader(http.StatusNotImplemented)
	case gpio.ErrAlreadyExported:
		fallthrough
	case gpio.ErrInvalidDirection:
		fallthrough
	case gpio.ErrInvalidPin:
		logrus.Warning("GPIO pin exporting error:", err)
		wr.WriteHeader(http.StatusBadRequest)
		wr.Write([]byte(err.Error()))
	default:
		logrus.Error("GPIO pin exporting failed:", err)
		wr.WriteHeader(http.StatusInternalServerError)
	}

}

func (gh *gpioHandler) deletePin(wr http.ResponseWriter, req *http.Request) {
	logrus.Debugln("deletePin() handler")
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
			wr.Write([]byte(err.Error()))
		} else {
			wr.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	wr.WriteHeader(http.StatusOK)
}

func (gh *gpioHandler) setPin(wr http.ResponseWriter, req *http.Request) {
	logrus.Debugln("setPin() handler")
	params := mux.Vars(req)
	pin, err := strconv.Atoi(params["pin"])
	if err != nil || pin < 0 {
		logrus.Errorln("Invalid pin number")
		wr.WriteHeader(http.StatusNotFound)
		return
	}

	buffer := make([]byte, 1024)
	n, err := req.Body.Read(buffer)
	if err != nil && err != io.EOF {
		logrus.Errorln("Failed to read request body:", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}
	if n == 0 {
		logrus.Errorln("Empty body in request")
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	requestData := pinValuePointer{}
	err = json.Unmarshal(buffer[:n], &requestData)
	if err != nil {
		logrus.Errorln("Failed to unmarshall request data:", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	if requestData.Value == nil {
		logrus.Debug("Incomplete request data")
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	err = gh.ctrl.SetValue(pin, *requestData.Value)
	if err == nil {
		wr.WriteHeader(http.StatusOK)
		return
	}

	if err == gpio.ErrInvalidDirection || err == gpio.ErrNotExported {
		logrus.Warnln("Failed to set pin value due to configuration error:", err)
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.Errorln("Failed to set pin value:", err)
	wr.WriteHeader(http.StatusInternalServerError)
}

func (gh *gpioHandler) getPin(wr http.ResponseWriter, req *http.Request) {
	logrus.Debugln("getPin() handler")
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

	respData := pinValue{Pin: pin, Value: val}
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
	logrus.Debugln("getAllPins() handler")
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

	result := make([]pinConfig, 0, len(pins))
	for k, v := range pins {
		result = append(result, pinConfig{Pin: k, Direction: gpio.DirectionToString(v)})
	}

	buffer, err := json.Marshal(result)
	if err != nil {
		logrus.Errorln("Error when marshalling pin data:", err)
		if err == gpio.ErrNotImplemented {
			wr.WriteHeader(http.StatusNotImplemented)
			return
		}
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.Write(buffer)
	wr.WriteHeader(http.StatusOK)
}

// TODO re-think this and maybe unify structures used in code
type pinConfigPointer struct {
	Pin       *int    `json:"pin"`
	Direction *string `json:"direction"`
}

type pinConfig struct {
	Pin       int    `json:"pin"`
	Direction string `json:"direction"`
}

type pinValue struct {
	Pin   int `json:"pin"`
	Value int `json:"value"`
}

type pinValuePointer struct {
	Pin   *int `json:"pin"`
	Value *int `json:"value"`
}
