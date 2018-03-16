package mypackage

// Sliceint Describes the a int slice
type Sliceint []*int 

// Add a value to slice
func (s *Sliceint) Add( values ...*int) {
	*s = append(*s, values...)
}
//IndexOf returns the index of 
func ( s *Sliceint) IndexOf( value *int ) int {
	for i, v := range *s {
		if v == value {
			return i
		}
	}
	return -1
}


