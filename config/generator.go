package config

// Generator configuration
type Generator struct {
	Description string       `yaml:"description"`
	Aliases     []string     `yaml:"aliases"`
	Files       []FileTarget `yaml:"files"`
	Vars        []UserVar    `yaml:"vars"`
}

// FileTarget composed by source, and Target
type FileTarget struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}
