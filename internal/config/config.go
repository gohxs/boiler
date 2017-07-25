package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config boiler.yaml
type Config struct {
	Vars map[string]Var `yaml:"vars"`
}

type Var struct {
	Question string `yaml:"question"`
}

// FromFile load config from file
func FromFile(configPath string) (*Config, error) {
	// Check for file or ignore
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	configData, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	config := Config{}
	yaml.Unmarshal(configData, &config)

	return &config, nil

}
