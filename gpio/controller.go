package gpio

import "errors"

type controller struct {
	basePath string
}

func (c *controller) SetValue(pin, value int) error {
	return errors.New("not implemented")
}

func (c *controller) GetValue(pin int) (int, error) {
	return 0, errors.New("not implemented")
}

func (c *controller) ExportPin(pin int, mode Direction) error {
	return errors.New("not implemented")
}

func (c *controller) UnexportPin(pin int) error {
	return errors.New("not implemented")
}

func (c *controller) ListExportedPins() (map[string]Direction, error) {
	var result map[string]Direction
	return result, errors.New("not implemented")
}

func CreateController(gpioPath string) Controller {
	return &controller{basePath: gpioPath}
}
