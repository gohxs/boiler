package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config boiler.yaml
type Config struct {
	Description string               `yaml:"description"`
	UserVars    []UserVar            `yaml:"vars"`
	Generators  map[string]Generator `yaml:"generators"`
}

// FromFile load config from file
func FromFile(configPath string, config *Config) error {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(configData, config)

}

// SaveFile marshals into file
func SaveFile(configPath string, config *Config) error {

	bdata, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, bdata, os.FileMode(0644))
}
