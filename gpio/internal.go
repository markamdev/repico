package gpio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
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

func isAvailable(pin int, dir Direction) (bool, error) {
	elm, ok := usedPins[pin]
	if !ok {
		// pin not configured
		return false, nil
	}
	if elm != dir {
		// pin used BUT set to different direction
		return false, errors.New("Pin configured with different direction")
	}
	// pin used and set to desired direction
	return true, nil
}

func addPin(pin int, dir Direction) error {
	elm, ok := usedPins[pin]
	if ok {
		// already exist in map - check if same direction selected
		if elm == dir {
			// ok, no conflict
			return nil
		} else {
			// pin direction changing temporary not supported
			return errors.New("Pin alrady used with another direction")
		}
	}

	err := exportGPIO(pin, dir)
	if err != nil {
		return err
	}
	usedPins[pin] = dir
	return nil
}

func removePin(pin int) error {
	_, ok := usedPins[pin]
	if !ok {
		// do not return error to allow double removal call
		return nil
	}
	return unexportGPIO(pin)
}

func exportGPIO(pin int, dir Direction) error {
	pinID := []byte(strconv.Itoa(pin))
	err := ioutil.WriteFile(pathGpioExport, pinID, 0770)
	if err != nil {
		return fmt.Errorf("Failed to export pin '%v': %v", string(pinID), err)
	}
	pinValuePath := createDirectionPath(pin)
	var directionString string
	if dir == Input {
		directionString = "in"
	} else {
		directionString = "out"
	}
	// TODO add proper checking for exported GPIO
	time.Sleep(time.Millisecond * 200)
	err = ioutil.WriteFile(pinValuePath, []byte(directionString), 0770)
	if err != nil {
		unexportGPIO(pin)
		return fmt.Errorf("Failed to set direction '%v' for pin '%v': %v", directionString, string(pinID), err)
	}
	return nil
}

func unexportGPIO(pin int) error {
	pinID := []byte(strconv.Itoa(pin))
	err := ioutil.WriteFile(pathGpioUnexport, pinID, 0770)
	if err != nil {
		return fmt.Errorf("Failed to unexport pin '%v': %v", string(pinID), err)
	}
	return err
}

func createValuePath(pin int) string {
	result := pathGpioPinBase + strconv.Itoa(int(pin)) + pathValueSuffix
	return result
}

func createDirectionPath(pin int) string {
	result := pathGpioPinBase + strconv.Itoa(int(pin)) + pathDirectionSuffix
	return result
}

func setValue(pin, value int) error {
	pinPath := createValuePath(pin)
	var err error
	if value == 0 {
		err = ioutil.WriteFile(pinPath, []byte{'0'}, 0700)
	} else {
		err = ioutil.WriteFile(pinPath, []byte{'1'}, 0700)
	}
	return err
}
