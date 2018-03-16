package mypackage

import "testing"


func TestIntAdd(t *testing.T) {

	t.Log("Add for int")
	vl := Sliceint{}


	vl.Add(new(int))
	vl.Add(new(int))
	vl.Add(new(int))

	v := len(vl)
	t.Log("len(vl) : ",v)
	if v != 3 {
		t.FailNow()
	}

}
func TestIntIndexOf( t *testing.T  ) {
	t.Log("Index of int")
	vl := Sliceint{}
	
	v2 := new(int)
	vl.Add(new(int))
	vl.Add(v2)
	vl.Add(new(int))

	v := vl.IndexOf(v2)
	t.Log("IndexOf v2 ==",v)	
	
	if v != 1 {
		t.FailNow()
	}
}
