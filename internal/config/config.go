package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config boiler.yaml
type Config struct {
	UserVars   []UserVar            `yaml:"vars"`
	Generators map[string]Generator `yaml:"generators"`
}

// FromFile load config from file
func FromFile(configPath string, config *Config) error {
	//log.Println("Reading config file:", configPath)
	// Check for file or ignore
	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	configData, err := ioutil.ReadAll(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		return err
	}
	return nil

}
