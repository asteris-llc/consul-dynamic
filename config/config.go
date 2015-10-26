package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Data map[string]interface{}
}

const (
	ERR_No_Config = "No consul configuration path"
	ERR_Not_Exist = "Consul configuration path not found: %s"
)

func ReadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf(ERR_No_Config)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf(ERR_Not_Exist, path)
	}

	cdata, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rval := Config{
		Data: make(map[string]interface{}),
	}

	if err := json.Unmarshal(cdata, &rval.Data); err != nil {
		return nil, err
	}

	return &rval, nil
}

func (c *Config) Write(path string) error {
	if path == "" {
		return fmt.Errorf(ERR_No_Config)
	}

	data, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, data, 0640); err != nil {
		return err
	}

	return nil
}
