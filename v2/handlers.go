package v2

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/markamdev/repico/gpio"
	"github.com/markamdev/repico/server"
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
		server.WriteMessage(wr, http.StatusInternalServerError, "unexpected internal error")
		return
	}
	if n == 0 {
		logrus.Warnln("Empty request body")
		server.WriteMessage(wr, http.StatusBadRequest, "empty request body")
		return
	}

	var pinDesc pinConfigPointer
	err = json.Unmarshal(buffer[:n], &pinDesc)
	if err != nil {
		logrus.Error("Unable to unmarshal request data:", err)
		server.WriteMessage(wr, http.StatusBadRequest, "invalid request body")
		return
	}

	if pinDesc.Pin == nil || pinDesc.Direction == nil {
		logrus.Error("No proper pin description in body")
		server.WriteMessage(wr, http.StatusBadRequest, "incorrect pin description")
		return
	}

	err = gh.ctrl.ExportPin(*pinDesc.Pin, gpio.StringToDirection(*pinDesc.Direction))

	switch err {
	case nil:
		wr.WriteHeader(http.StatusOK)
	case gpio.ErrNotImplemented:
		server.WriteMessage(wr, http.StatusNotImplemented, "not implemented")
	case gpio.ErrAlreadyExported:
		fallthrough
	case gpio.ErrInvalidDirection:
		fallthrough
	case gpio.ErrInvalidPin:
		logrus.Warning("GPIO pin exporting error:", err)
		server.WriteMessage(wr, http.StatusBadRequest, err.Error())
	default:
		logrus.Error("GPIO pin exporting failed:", err)
		server.WriteMessage(wr, http.StatusInternalServerError, "unexpected internal error")
	}

}

func (gh *gpioHandler) deletePin(wr http.ResponseWriter, req *http.Request) {
	logrus.Debugln("deletePin() handler")
	params := mux.Vars(req)
	pin, err := strconv.Atoi(params["pin"])
	if err != nil || pin < 0 {
		logrus.Errorln("Invalid pin number:", err)
		server.WriteMessage(wr, http.StatusBadRequest, "invalid pin selection")
		return
	}

	err = gh.ctrl.UnexportPin(pin)
	if err != nil {
		logrus.Errorf("Failed to unexport pin '%d': %v\n", pin, err.Error())
		if err == gpio.ErrNotExported {
			server.WriteMessage(wr, http.StatusBadRequest, err.Error())
		} else {
			server.WriteMessage(wr, http.StatusInternalServerError, "unexpected internal error")
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
		server.WriteMessage(wr, http.StatusBadRequest, "invalid pin selection")
		return
	}

	buffer := make([]byte, 1024)
	n, err := req.Body.Read(buffer)
	if err != nil && err != io.EOF {
		logrus.Errorln("Failed to read request body:", err)
		server.WriteMessage(wr, http.StatusInternalServerError, "body reading error")
		return
	}
	if n == 0 {
		logrus.Errorln("Empty body in request")
		server.WriteMessage(wr, http.StatusBadRequest, "empty request body")
		return
	}

	requestData := pinValuePointer{}
	err = json.Unmarshal(buffer[:n], &requestData)
	if err != nil {
		logrus.Errorln("Failed to unmarshall request data:", err)
		server.WriteMessage(wr, http.StatusInternalServerError, "unexpected internal error")
		return
	}

	if requestData.Value == nil {
		logrus.Debug("Incomplete request data")
		server.WriteMessage(wr, http.StatusBadRequest, "invalid incomplete request data")
		return
	}

	err = gh.ctrl.SetValue(pin, *requestData.Value)
	if err == nil {
		wr.WriteHeader(http.StatusOK)
		return
	}

	if err == gpio.ErrInvalidDirection {
		logrus.Warnln("Invalid pin direction")
		server.WriteMessage(wr, http.StatusBadRequest, "invalid pin direction")
		return
	}
	if err == gpio.ErrNotExported {
		logrus.Warnln("Pin not exported")
		server.WriteMessage(wr, http.StatusBadRequest, "pin not exported")
		return
	}

	logrus.Errorln("Failed to set pin value:", err)
	server.WriteMessage(wr, http.StatusInternalServerError, "unexpected internal error")
}

func (gh *gpioHandler) getPin(wr http.ResponseWriter, req *http.Request) {
	logrus.Debugln("getPin() handler")
	params := mux.Vars(req)
	pin, err := strconv.Atoi(params["pin"])
	if err != nil || pin < 0 {
		logrus.Errorln("Invalid pin number")
		server.WriteMessage(wr, http.StatusBadRequest, "invalid pin selection")
		return
	}

	val, err := gh.ctrl.GetValue(pin)
	if err != nil {
		logrus.Errorf("Failed to get pin '%d' value: %v\n", pin, err.Error())
		if err == gpio.ErrNotExported {
			server.WriteMessage(wr, http.StatusBadRequest, "pin not exported")
		} else {
			server.WriteMessage(wr, http.StatusInternalServerError, "unexpected internal error")
		}
		return
	}
	logrus.Tracef("Pin: %d Value: %d", pin, val)

	respData := pinValue{Pin: pin, Value: val}
	buffer, err := json.Marshal(respData)
	if err != nil {
		logrus.Errorln("Failed to marshal result:", err)
		server.WriteMessage(wr, http.StatusInternalServerError, "unexpected internal error")
		return
	}

	server.WriteResponse(wr, http.StatusOK, buffer)
}

func (gh *gpioHandler) getAllPins(wr http.ResponseWriter, req *http.Request) {
	logrus.Debugln("getAllPins() handler")
	pins, err := gh.ctrl.ListExportedPins()
	if err != nil {
		logrus.Errorln("Error when listing GPIO pins:", err)
		server.WriteMessage(wr, http.StatusInternalServerError, "unexpected internal error")
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
			server.WriteMessage(wr, http.StatusNotImplemented, "not implemented")
			return
		}
		server.WriteMessage(wr, http.StatusInternalServerError, "unexpected internal error")
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}
	server.WriteResponse(wr, http.StatusOK, buffer)
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
