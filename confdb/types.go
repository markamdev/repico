package confdb

// ConfigStorage is an interface for future configuration databases
type ConfigStorage interface {
	Load(name string) (AppConfig, error)
	Save(content AppConfig) error
}

// types below (config description) are temporary just a copy of RestConfig
// in future will be replaced with other types

// AppConfig describes configuration of RePiCo application and is used for permanent config storing
type AppConfig struct {
	// Name contains configuration name (can be usefull when debugging and log checking)
	Name string `json:"name"`
	// Pins contains list of pins that should be configured for reading/setting values
	Pins []PinConfig `json:"pins"`
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
