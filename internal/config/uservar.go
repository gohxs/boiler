package config

import (
	yaml "gopkg.in/yaml.v2"
)

// Vars
type Vars []UserVar

// UserVar user defined vars
type UserVar struct {
	Name     string
	Default  string `yaml:"default"`
	Question string `yaml:"question"`
}

// UnmarshalYAML implementation for yaml decoder, for ordered key
func (vr *Vars) UnmarshalYAML(unmarshal func(interface{}) error) error {
	orderm := yaml.MapSlice{}
	unmarshal(&orderm) // Ordered
	// Unmarshal Var
	m := map[string]UserVar{}
	unmarshal(&m)

	for _, v := range orderm {
		lvar := m[v.Key.(string)]
		lvar.Name = v.Key.(string)
		*vr = append(*vr, lvar)
	}

	return nil
}
