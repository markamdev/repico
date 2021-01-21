package internal

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
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
