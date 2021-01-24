package gpio

import (
	"errors"
	"fmt"
)

// SetGPIO tries to set given value to GPIO port
func SetGPIO(number, value int) error {
	av, err := isAvailable(number, Output)
	if err != nil {
		return fmt.Errorf("Pin %v not available for setting: %v", number, err)
	}
	if !av {
		return fmt.Errorf("Pin %v not configured", number)
	}
	if value < 0 || value > 1 {
		return fmt.Errorf("Invalid value %v for GPIO", value)
	}
	return setValue(number, value)
}

// GetGPIO tries to check current status of given GPIO port
func GetGPIO(number int) (int, error) {
	//return false, errors.New("GPIO state get not implemented")

	av, err := isAvailable(number, Input)
	if err != nil {
		return -1, fmt.Errorf("Pin %v not available for setting: %v", number, err)
	}
	if !av {
		return -1, fmt.Errorf("Pin %v not configured", number)
	}
	return getValue(number)
}

// EnableGPIO enables given GPIO pin in requested direction
func EnableGPIO(number int, dir Direction) error {
	// TODO: add number range verification
	if dir != Input && dir != Output {
		return errors.New("Invalid GPIO direction")
	}
	return addPin(number, dir)
}

// DisableGPIO disables (unexport) given GPIO pin
func DisableGPIO(number int) error {
	// TODO: add number range verification
	return removePin(number)
}
