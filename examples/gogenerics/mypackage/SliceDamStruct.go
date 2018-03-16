package mypackage

// SliceDamStruct Describes the a DamStruct slice
type SliceDamStruct []*DamStruct 

// Add a value to slice
func (s *SliceDamStruct) Add( values ...*DamStruct) {
	*s = append(*s, values...)
}
//IndexOf returns the index of 
func ( s *SliceDamStruct) IndexOf( value *DamStruct ) int {
	for i, v := range *s {
		if v == value {
			return i
		}
	}
	return -1
}


