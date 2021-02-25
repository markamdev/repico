package confdb

import (
	"fmt"
)

const (
	// StorageTypeJSON represents type of storage based on JSON file
	StorageTypeJSON = "json"
)

// GetStorage returns instance of ConfigStorage with selected type (stype)
// additional params for saver (ex. path for file based saver) should be passed as a second param
func GetStorage(stype, sparams string) (ConfigStorage, error) {
	if stype == StorageTypeJSON {
		return createJSONStorage(sparams), nil
	}
	return nil, fmt.Errorf("No support implemented for '%v' storage", stype)
}
