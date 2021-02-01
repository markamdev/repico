package internal

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/markamdev/repico/gpio"
)

// Params stores all application parameters
type Params struct {
	Pins []int
}

const (
	flagNamePins = "pins"
)

// ReadConfig checks command line params and system environment variables
// to configure application
func ReadConfig() (Params, error) {
	var result Params
	// read GPIO numbers (not a board pin numbers!!)
	var pinList string
	flag.StringVar(&pinList, flagNamePins, "", "List of comma separated GPIO (out) numbers managed by application")
	flag.Parse()
	if len(pinList) == 0 {
		// at the moment empty pin list is not supported
		return Params{}, errors.New("No managed pin defined")
	}
	pinStrings := strings.Split(pinList, ",")
	result.Pins = make([]int, len(pinStrings))
	for idx, name := range pinStrings {
		pinNum, converr := strconv.Atoi(name)
		if converr != nil {
			return Params{}, fmt.Errorf("Invalid pin number at: %d", idx)
		}
		result.Pins[idx] = pinNum
	}

	// no error till now
	return result, nil
}

// PinConfig contains configuration of sigle GPIO pin
type PinConfig struct {
	// Number is a GPIO number according to board spec
	Number int `json:"number"`
	// Aliast is an optional name for this GPIO pin (can be used for reading/setting pin state)
	Alias string `json:"alias,omitempty"`
	// Direction can be set only to "in" and "out" - other values are rejected
	Direction string `json:"direction"`
}

// RestConfig contains description of configuration applied usint REST API (/config URI)
type RestConfig struct {
	// Name contains configuration name (can be usefull when debugging and log checking)
	Name string `json:"name"`
	// Pins contains list of pins that should be configured for reading/setting values
	Pins []PinConfig `json:"pins"`
}

// ValidateConfig verifies if provided app configuration is correct, returns error if not
func ValidateConfig(conf RestConfig) error {
	errs := make([]string, 0)

	if len(conf.Pins) == 0 {
		return errors.New("Empty pins configuration is not allowed")
	}

	var msg string
	for pos, pinDesc := range conf.Pins {
		if pinDesc.Number < 0 {
			msg = fmt.Sprintf("Invalid pin number %v at position %v", pinDesc.Number, pos)
			errs = append(errs, msg)
		}
		if pinDesc.Direction != "in" && pinDesc.Direction != "out" {
			msg = fmt.Sprintf("Invalid pin direction %v at position %v", pinDesc.Direction, pos)
			errs = append(errs, msg)
		}
		regMatch, _ := regexp.Match(`^[a-zA-Z]+[0-9]*`, []byte(pinDesc.Alias))
		if !regMatch {
			msg = fmt.Sprintf("Invalid alias '%v' at position %v", pinDesc.Alias, pos)
			errs = append(errs, msg)
		}
	}

	if len(errs) > 0 {
		errString := strings.Join(errs, ";")
		return errors.New(errString)
	}
	// no error found
	return nil
}

// ApplyConfig tries to apply GPIO pins configuration based on provided data
func ApplyConfig(conf RestConfig) error {
	for _, pin := range conf.Pins {
		var dir gpio.Direction
		if pin.Direction == "in" {
			dir = gpio.Input
		} else {
			// assumes that config has been already validated
			dir = gpio.Output
		}
		err := gpio.EnableGPIO(pin.Number, dir)
		if err != nil {
			return fmt.Errorf("Failed to set config for pin %v with direction %v. Error: %v",
				pin.Number, dir, err.Error())
		}
	}
	return nil
}
