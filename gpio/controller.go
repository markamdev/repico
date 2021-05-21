package gpio

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

type controller struct {
	basePath string
}

func (c *controller) SetValue(pin, value int) error {
	return ErrNotImplemented
}

func (c *controller) GetValue(pin int) (int, error) {
	return 0, ErrNotImplemented
}

func (c *controller) ExportPin(pin int, mode Direction) error {
	if mode != Input && mode != Output {
		return ErrInvalidDirection
	}
	if pin < 0 {
		return ErrInvalidPin
	}
	pinString := strconv.Itoa(pin)

	if isExported(pinString) {
		return ErrAlreadyExported
	}

	err := exportPin(pinString)
	if err != nil {
		return ErrUnknown
	}

	err = setDirection(pinString, DirectionToString(mode))
	if err != nil {
		unexportPin(pinString)
	}

	return err
}

func (c *controller) UnexportPin(pin int) error {
	return ErrNotImplemented
}

func (c *controller) ListExportedPins() (map[int]Direction, error) {
	var result map[int]Direction
	return result, ErrNotImplemented
}

func CreateController(gpioPath string) Controller {
	logrus.Traceln("gpio.CreateController()")
	return &controller{basePath: gpioPath}
}
