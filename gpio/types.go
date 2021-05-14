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

func DirectionToString(dr Direction) string {
	switch dr {
	case Input:
		return "in"
	case Output:
		return "out"
	default:
		return "-"
	}
}

func StringToDirection(dr string) Direction {
	switch dr {
	case "in":
		return Input
	case "out":
		return Output
	default:
		return Invalid
	}
}

// PinConfig describes GPIO pin: it's number, direction and optional label (alias)
type PinConfig struct {
	Number int
	Mode   Direction
	Label  string
}

// PinState describes state of GPIO pin
type PinState struct {
	Number int
	Alias  string
	State  int
}

// Controller is an interface of GPIO controlling object
type Controller interface {
	SetValue(pin, value int) error
	GetValue(pin int) (int, error)
	ExportPin(pin int, mode Direction) error
	UnexportPin(pin int) error
	ListExportedPins() (map[string]Direction, error)
}
