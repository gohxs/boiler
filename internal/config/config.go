package config

import (
	"io/ioutil"

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
	//log.Println("Reading config file:", configPath)
	// Check for file or ignore
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		return err
	}
	return nil

}
