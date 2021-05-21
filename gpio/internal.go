package gpio

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
