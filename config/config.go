package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/markamdev/repico/confdb"
	"github.com/markamdev/repico/gpio"
)

const (
	configFilePath = "repico.confdb"
)

var configStorage confdb.ConfigStorage

// LoadDefault reads all application configs (if any stored) ant returns "default" one
func LoadDefault() (RestConfig, error) {
	// temporary hardcoded loading config from JSON
	// TODO storage selection and saving for future save operations
	str, err := confdb.GetStorage(confdb.StorageTypeJSON, configFilePath)
	if err != nil {
		return RestConfig{}, err
	}
	configStorage = str
	// TODO replace with loading last storage
	cnf, err := str.Load("default")
	if err != nil {
		return RestConfig{}, err
	}

	result := RestConfig{}
	result.Name = cnf.Name
	for _, pin := range cnf.Pins {
		pinCfg := PinConfig{}
		pinCfg.Alias = pin.Alias
		pinCfg.Direction = pin.Direction
		result.Pins = append(result.Pins, pinCfg)
	}
	return result, nil
}

// Store stores gpio pin configuration in storage
func Store(cfg RestConfig) error {
	appCfg := confdb.AppConfig{}
	appCfg.Name = cfg.Name
	for _, pin := range cfg.Pins {
		pinCfg := confdb.PinConfig{}
		pinCfg.Alias = pin.Alias
		pinCfg.Direction = pin.Direction
		appCfg.Pins = append(appCfg.Pins, pinCfg)
	}
	return configStorage.Save(appCfg)
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
