package gpio

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var usedPins = map[int]Direction{}

// functions in this file are only used internally so param validation is omited

const (
	pathGpioBase        = "/sys/class/gpio"
	pathGpioExport      = "/sys/class/gpio/export"
	pathGpioUnexport    = "/sys/class/gpio/unexport"
	pathGpioPinBase     = "/sys/class/gpio/gpio"
	pathDirectionSuffix = "/direction"
	pathValueSuffix     = "/value"
)

func isExported(pin string) bool {
	pinPath := pathGpioPinBase + pin
	logrus.Trace("isExported():", pinPath)
	_, err := os.Stat(pinPath)
	if err != nil {
		logrus.Traceln("isExported() Stat() error:", err)
		return false
	}
	return true
}

func exportPin(pin string) error {
	fExport, err := os.OpenFile(pathGpioExport, os.O_WRONLY, 0755)
	if err != nil {
		logrus.Traceln("exportPin() export opening failed:", err)
		return ErrUnknown
	}
	defer fExport.Close()

	_, err = fExport.WriteString(pin)
	if err != nil {
		logrus.Traceln("exportPin() export error:", err)
		return ErrUnknown
	}

	return nil
}

func setDirection(pin, dir string) error {
	dirPath := pathGpioPinBase + pin + pathDirectionSuffix
	logrus.Traceln("setDirection(): ", dirPath)

	fDir, err := os.OpenFile(dirPath, os.O_WRONLY, 0755)
	if err != nil {
		return ErrUnknown
	}
	defer fDir.Close()

	_, err = fDir.WriteString(dir)
	if err != nil {
		logrus.Traceln("setDirection() writing error:", err)
		return ErrUnknown
	}

	return nil
}

func unexportPin(pin string) error {
	fUnexport, err := os.OpenFile(pathGpioUnexport, os.O_WRONLY, 0755)
	if err != nil {
		logrus.Traceln("unexportPin() unexport opening failed:", err)
		return ErrUnknown
	}
	defer fUnexport.Close()

	_, err = fUnexport.WriteString(pin)
	if err != nil {
		logrus.Traceln("exportPin() export error:", err)
		return ErrUnknown
	}

	return nil
}

func isOutput(pin string) (bool, error) {
	dirPath := pathGpioPinBase + pin + pathDirectionSuffix
	fDirection, err := os.OpenFile(dirPath, os.O_RDONLY, 0755)
	if err != nil {
		logrus.Traceln("isOutput() cannot open direction file:", err)
		return false, ErrUnknown
	}
	defer fDirection.Close()

	buffer := make([]byte, 16)
	n, err := fDirection.Read(buffer)
	if err != nil && err != io.EOF {
		logrus.Traceln("isOutput() failed to read direction file:", err)
		return false, ErrUnknown
	}
	dirString := strings.TrimRight(string(buffer[:n]), "\n\r")
	if dirString == "out" {
		return true, nil
	}
	return false, nil
}

func setValue(pin, value string) error {
	valuePath := pathGpioPinBase + pin + pathValueSuffix
	fValue, err := os.OpenFile(valuePath, os.O_WRONLY, 0755)
	if err != nil {
		logrus.Traceln("setValue() cannot open value file:", err)
		return ErrUnknown
	}
	defer fValue.Close()

	_, err = fValue.WriteString(value)
	if err != nil {
		logrus.Traceln("setValue() cannot set GPIO state:", err)
		return ErrUnknown
	}

	return nil
}

func getValue(pin string) (string, error) {
	valuePath := pathGpioPinBase + pin + pathValueSuffix
	fValue, err := os.OpenFile(valuePath, os.O_RDONLY, 0755)
	if err != nil {
		logrus.Traceln("getValue() cannot open value file:", err)
		return "", ErrUnknown
	}

	defer fValue.Close()

	buffer := make([]byte, 16)
	n, err := fValue.Read(buffer)
	if err != nil && err != io.EOF {
		logrus.Traceln("getValue() failed to read value file:", err)
		return "", ErrUnknown
	}

	return strings.TrimRight(string(buffer[:n]), "\r\n"), nil
}
