package v2

import (
	"errors"
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

	t.Run("list all pins", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNotImplemented, resRecorder.Code, "GPIO listing failed")
	})

	t.Run("add pin", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v2/gpio", body)
		resRecorder := httptest.NewRecorder()

		hndlr.ServeHTTP(resRecorder, req)

		assert.Equal(t, http.StatusNotImplemented, resRecorder.Code, "GPIO pin adding failed")
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
}

func (cs *controllerStub) SetValue(pin, value int) error {
	return errors.New("not implemented")
}

func (cs *controllerStub) GetValue(pin int) (int, error) {
	return 0, errors.New("not implemented")
}

func (cs *controllerStub) ExportPin(pin int, mode gpio.Direction) error {
	return errors.New("not implemented")
}

func (cs *controllerStub) UnexportPin(pin int) error {
	return errors.New("not implemented")
}

func (cs *controllerStub) ListExportedPins() (map[int]gpio.Direction, error) {
	var result map[int]gpio.Direction
	return result, errors.New("not implemented")
}

type bodyStub struct {
	dataToReturn  []byte
	errorToReturn error
}

func (bs *bodyStub) Read(buffer []byte) (int, error) {
	return len(bs.dataToReturn), bs.errorToReturn
}
