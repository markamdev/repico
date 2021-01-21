package gpio

// Direction defines GPIO pin direction
// Possible values are Unset, Input and Output
type Direction int

const (
	// Unset - edfault value, not set
	Unset Direction = iota
	// Input - GPIO pin set as input
	Input
	// Output - GPIO pin set as output
	Output
)

// PinConfig describes GPIO pin: it's number, direction and optional label (alias)
type PinConfig struct {
	Number int
	Mode   Direction
	Label  string
}
