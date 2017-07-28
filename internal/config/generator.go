package config

// Generator configuration
type Generator struct {
	Aliases     []string     `yaml:"aliases"`
	Files       []FileTarget `yaml:"files"`
	Flags       []string     `yaml:"flags"`
	Description string       `yaml:"description"`
	Vars        []UserVar    `yaml:"vars"`
	//Ext         string       `yaml:"ext"`
	//Target string       `yaml:"target"`
	//Source string       `yaml:"source"`
}

// FileTarget composed by source, and Target
type FileTarget struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}
