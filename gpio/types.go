package gpio

// Direction defines GPIO pin direction
// Possible values are Unset, Input and Output
type Direction int

const (
	// Invalid - edfault value, not set
	Invalid Direction = iota
	// Input - GPIO pin set as input
	Input
	// Output - GPIO pin set as output
	Output
)

const (
	directionIn  = "in"
	directionOut = "out"
)

func DirectionToString(dr Direction) string {
	switch dr {
	case Input:
		return directionIn
	case Output:
		return directionOut
	default:
		return "-"
	}
}

func StringToDirection(dr string) Direction {
	switch dr {
	case directionIn:
		return Input
	case directionOut:
		return Output
	default:
		return Invalid
	}
}

// Controller is an interface of GPIO controlling object
type Controller interface {
	SetValue(pin, value int) error
	GetValue(pin int) (int, error)
	ExportPin(pin int, mode Direction) error
	UnexportPin(pin int) error
	ListExportedPins() (map[int]Direction, error)
}
