package config

// UserVar user defined vars
type UserVar struct {
	Name     string   `yaml:"name"`
	Default  string   `yaml:"default"`
	Flag     string   `yaml:"flag"`
	Question string   `yaml:"question"`
	Choices  []string `yaml:"choices"`
}
