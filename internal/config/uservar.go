package config

// Vars
//type Vars []UserVar

// UserVar user defined vars
type UserVar struct {
	Name     string `yaml:"name"`
	Default  string `yaml:"default"`
	Flag     string `yaml:"flag"`
	Question string `yaml:"question"`
}

// UnmarshalYAML implementation for yaml decoder, for ordered key
/*func (vr *UserVar) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	//try string
	var tryStr string
	err = unmarshal(&tryStr)
	if err == nil { // If OK
		vr.Name = tryStr
		return
	}
	// Try map
	tryMap := map[string]interface{}{}
	err = unmarshal(&tryMap)
	if err != nil { // If ERR
		return err
	}

	// Copy map to struct
	typ := reflect.TypeOf(vr).Elem()
	val := reflect.ValueOf(vr).Elem()
	for i := 0; i < typ.NumField(); i++ {
		ftyp := typ.Field(i)
		fname := ftyp.Name

		if tag, ok := ftyp.Tag.Lookup("yaml"); ok {
			tagParts := strings.Split(tag, ",")
			fname = tagParts[0]
		}
		if fieldData, ok := tryMap[fname]; ok {
			val.Field(i).Set(reflect.Indirect(reflect.ValueOf(fieldData)))
		}
	}

	// somehow  copy map to struct
	return nil
}*/
