package gpio

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

type controller struct {
	basePath string
}

func (c *controller) SetValue(pin, value int) error {
	logrus.Traceln("gpio.controller.SetValue()")
	if value != 0 && value != 1 {
		return ErrInvalidValue
	}

	pinString := strconv.Itoa(pin)

	if !isExported(pinString) {
		return ErrNotExported
	}

	out, err := isOutput(pinString)
	if err != nil {
		logrus.Errorln("Failed to check mode:", err)
		return err
	}
	if !out {
		return ErrInvalidDirection
	}

	valueString := strconv.Itoa(value)
	return setValue(pinString, valueString)
}

func (c *controller) GetValue(pin int) (int, error) {
	logrus.Traceln("gpio.controller.GetValue()")
	pinString := strconv.Itoa(pin)

	if !isExported(pinString) {
		return -1, ErrNotExported
	}

	valString, err := getValue(pinString)
	if err != nil {
		logrus.Errorln("Failed to get value:", err)
		return -1, err
	}

	valInt, err := strconv.Atoi(valString)
	if err != nil {
		logrus.Errorln("Failed to convert value:", err)
		return -1, ErrUnknown
	}

	return valInt, nil
}

func (c *controller) ExportPin(pin int, mode Direction) error {
	logrus.Traceln("gpio.controller.ExportPin()")
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

	ErrInvalidValue := exportPin(pinString)
	if ErrInvalidValue != nil {
		return ErrUnknown
	}

	ErrInvalidValue = setDirection(pinString, DirectionToString(mode))
	if ErrInvalidValue != nil {
		unexportPin(pinString)
	}

	return ErrInvalidValue
}

func (c *controller) UnexportPin(pin int) error {
	logrus.Traceln("gpio.controller.UnexportPin()")
	pinString := strconv.Itoa(pin)
	if !isExported(pinString) {
		return ErrNotExported
	}

	return unexportPin(pinString)
}

func (c *controller) ListExportedPins() (map[int]Direction, error) {
	logrus.Traceln("gpio.controller.ListExportedPins()")
	result := map[int]Direction{}

	pins, err := listExported()
	if err != nil {
		return map[int]Direction{}, ErrUnknown
	}
	logrus.Debug("Currently detected pins:", pins)

	for _, pin := range pins {
		isOut, err := isOutput(pin)
		if err != nil {
			logrus.Warn("Error while checking direction for one of pins:", err)
			return map[int]Direction{}, ErrUnknown
		}
		pinInt, _ := strconv.Atoi(pin)
		if isOut {
			result[pinInt] = Output
		} else {
			result[pinInt] = Input
		}
	}
	return result, nil
}

func CreateController(gpioPath string) Controller {
	logrus.Traceln("gpio.CreateController()")
	return &controller{basePath: gpioPath}
}
