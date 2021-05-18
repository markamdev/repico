package v2

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/markamdev/repico/gpio"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	ctrl := &controllerStub{}
	body := &bodyStub{}

	hndlr := CreateHandler(ctrl)

	assert.NotEqual(t, nil, hndlr, "Returned handler should not be nil")

	t.Run("list all pins - no content", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNoContent, resRecorder.Code)
	})

	t.Run("list all pins - 2 pins returned", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()
		ctrl.mapToReturn = map[string]gpio.Direction{"1": gpio.Input, "2": gpio.Output}

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusOK, resRecorder.Code)

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

		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
	})

	t.Run("add pin - invalid body (no direction)", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		body.dataToReturn = []byte("{ \"pin\" : 1 }")

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
	})

	t.Run("add pin - invalid body (no pin)", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		body.dataToReturn = []byte("{ \"direction\" : \"out\" }")

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
	})

	t.Run("add pin - correct case", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		body.dataToReturn = []byte("{ \"pin\" : 1, \"direction\" : \"out\" }")

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusOK, resRecorder.Code)
	})

	t.Run("delete pin - invalid pin number", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v2/gpio/-10", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNotFound, resRecorder.Code)
	})

	t.Run("delete pin - invalid path (no number)", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v2/gpio/xyz", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNotFound, resRecorder.Code)
	})

	t.Run("delete pin - unexported pin", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v2/gpio/2", body)
		resRecorder := httptest.NewRecorder()

		ctrl.errorToReturn = gpio.ErrNotExported

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
	})

	t.Run("delete pin - correct case", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v2/gpio/2", body)
		resRecorder := httptest.NewRecorder()

		ctrl.errorToReturn = nil

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusOK, resRecorder.Code)
	})

	t.Run("get pin - invalid pin number", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v2/gpio/-10", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNotFound, resRecorder.Code)
	})

	t.Run("get pin - invalid path (no number)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v2/gpio/xyz", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNotFound, resRecorder.Code)
	})

	t.Run("get pin - unexported pin", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v2/gpio/2", body)
		resRecorder := httptest.NewRecorder()

		ctrl.errorToReturn = gpio.ErrNotExported

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
	})

	t.Run("get pin - correct case", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v2/gpio/2", body)
		resRecorder := httptest.NewRecorder()

		ctrl.errorToReturn = nil
		ctrl.valueToReturn = 1

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusOK, resRecorder.Code)

		buffer := make([]byte, 1024)
		n, err := resRecorder.Body.Read(buffer)
		expectedErrors := []error{nil, io.EOF}
		assert.Contains(t, expectedErrors, err, "Response body reading error")

		var respData struct {
			Pin   int `json:"pin"`
			Value int `json:"value"`
		}
		err = json.Unmarshal(buffer[:n], &respData)
		assert.NoError(t, err)
		assert.Equal(t, 2, respData.Pin)
		assert.Equal(t, ctrl.valueToReturn, respData.Value)
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
	return cs.valueToReturn, cs.errorToReturn
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
