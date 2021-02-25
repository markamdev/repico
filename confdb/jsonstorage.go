package confdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// jsonStorage represents object for JSON based application config "database"
type jsonStorage struct {
	path    string
	configs []AppConfig
}

// createJSONStorage creates new JSON based configuration storage
func createJSONStorage(path string) ConfigStorage {
	result := jsonStorage{path: path}
	return &result
}

// Load returns configuration saved in storage under given name
func (js *jsonStorage) Load(name string) (AppConfig, error) {
	buf, err := ioutil.ReadFile(js.path)
	if err != nil {
		return AppConfig{}, fmt.Errorf("Failed to read file with confing: %v", err.Error())
	}
	loaded := []AppConfig{}
	err = json.Unmarshal(buf, &loaded)
	if err != nil {
		return AppConfig{}, fmt.Errorf("Failed to unmarshal file content: %v", err.Error())
	}
	js.configs = loaded
	for _, cfg := range js.configs {
		if cfg.Name == name {
			return cfg, nil
		}
	}

	return AppConfig{}, fmt.Errorf("Failed to load %v config: Not found", name)
}

// Save writes in storage configuration given as 'content'
func (js *jsonStorage) Save(content AppConfig) error {
	if content.Name == "" || len(content.Pins) == 0 {
		return errors.New("Invalid (incomplete) config given. Saving skipped")
	}
	cfgIndex := -1
	for idx, cfg := range js.configs {
		if cfg.Name == content.Name {
			cfgIndex = idx
			break
		}
	}
	if cfgIndex >= 0 {
		// overwrite existing config
		js.configs[cfgIndex] = content
	} else {
		// add new config at the end
		js.configs = append(js.configs, content)
	}
	buffer, err := json.Marshal(js.configs)
	if err != nil {
		return fmt.Errorf("Failed to prepare date for storing: %v", err.Error())
	}
	err = ioutil.WriteFile(js.path, buffer, os.FileMode(os.O_RDWR))
	if err != nil {
		return fmt.Errorf("Failed to store configs in file: %v", err.Error())
	}
	return nil
}
