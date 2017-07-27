package config

// Generator configuration
type Generator struct {
	Files       []FileTarget `yaml:"files"`
	Flags       []string     `yaml:"flags"`
	Ext         string       `yaml:"ext"`
	Description string       `yaml:"description"`
	Vars        []UserVar    `yaml:"vars"`
	//Target string       `yaml:"target"`
	//Source string       `yaml:"source"`
}

type FileTarget struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

/*func (f *FileTarget) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s := []string{}
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	if len(s) != 2 {
		return errors.New("Must have 2 elements")
	}

	f.Source = s[0]
	f.Target = s[1]

	return nil
}*/
