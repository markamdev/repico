package gpio

import (
	"os"

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
	return ErrNotImplemented
}
