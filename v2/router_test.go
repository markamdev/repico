package v2

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/markamdev/repico/gpio"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	ctrl := &controllerStub{}
	body := &bodyStub{}

	hndlr := CreateHandler(ctrl)

	assert.NotEqual(t, nil, hndlr, "Returned handler should not be nil")

	t.Run("list all pins - no content", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNoContent, resRecorder.Code, "GPIO listing failed")
	})

	t.Run("list all pins - 2 pins returned", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()
		ctrl.mapToReturn = map[string]gpio.Direction{"1": gpio.Input, "2": gpio.Output}

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusOK, resRecorder.Code, "GPIO listing failed")

		type listItem struct {
			Pin       string
			Direction string
		}
		gpioList := []listItem{}
		json.Unmarshal(resRecorder.Body.Bytes(), &gpioList)
		assert.Equal(t, 2, len(gpioList), "Expected two elements but got %d:", len(gpioList))
	})

	t.Run("add pin - empty body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		body.dataToReturn = []byte{}

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusBadRequest, resRecorder.Code, "Invalid GPIO pin adding status")
	})

	t.Run("add pin - invalid body (no direction)", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		body.dataToReturn = []byte("{ \"pin\" : 1 }")

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusBadRequest, resRecorder.Code, "Invalid GPIO pin adding status")
	})

	t.Run("add pin - invalid body (no pin)", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		body.dataToReturn = []byte("{ \"direction\" : \"out\" }")

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusBadRequest, resRecorder.Code, "Invalid GPIO pin adding status")
	})

	t.Run("add pin - correct case", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		body.dataToReturn = []byte("{ \"pin\" : 1, \"direction\" : \"out\" }")

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusOK, resRecorder.Code, "GPIO pin adding failed")
	})

	t.Run("delete pin", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v2/gpio/1", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNotImplemented, resRecorder.Code, "GPIO pin deletion failed")
	})

	t.Run("get pin", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v2/gpio/1", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNotImplemented, resRecorder.Code, "GPIO reading failed")
	})

	t.Run("set pin", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/v2/gpio/1", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNotImplemented, resRecorder.Code, "GPIO pin setting failed")
	})
}

type controllerStub struct {
	valueToReturn int
	errorToReturn error
	mapToReturn   map[string]gpio.Direction
}

func (cs *controllerStub) SetValue(pin, value int) error {
	return cs.errorToReturn
}

func (cs *controllerStub) GetValue(pin int) (int, error) {
	return 0, cs.errorToReturn
}

func (cs *controllerStub) ExportPin(pin int, mode gpio.Direction) error {
	return cs.errorToReturn
}

func (cs *controllerStub) UnexportPin(pin int) error {
	return cs.errorToReturn
}

func (cs *controllerStub) ListExportedPins() (map[string]gpio.Direction, error) {
	return cs.mapToReturn, cs.errorToReturn
}

type bodyStub struct {
	dataToReturn []byte
}

func (bs *bodyStub) Read(buffer []byte) (int, error) {
	copy(buffer, bs.dataToReturn)
	return len(bs.dataToReturn), io.EOF
}
