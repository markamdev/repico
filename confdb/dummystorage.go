package confdb

import "errors"

// DummyStorage is a stub storage class used (temporary) as a safe solution for not initialized storage
type DummyStorage struct {
}

// Load implements ConfigStorage.Load
func (dm *DummyStorage) Load(name string) (AppConfig, error) {
	return AppConfig{}, errors.New("Message from DummyStorage: storage not initialized")
}

// Save implements ConfigStorage.Save
func (dm *DummyStorage) Save(content AppConfig) error {
	return errors.New("Message from DummyStorage: storage not initialized")
}
