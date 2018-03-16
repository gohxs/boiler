package mypackage

import "testing"


func TestDamStructAdd(t *testing.T) {

	t.Log("Add for DamStruct")
	vl := SliceDamStruct{}


	vl.Add(new(DamStruct))
	vl.Add(new(DamStruct))
	vl.Add(new(DamStruct))

	v := len(vl)
	t.Log("len(vl) : ",v)
	if v != 3 {
		t.FailNow()
	}

}
func TestDamStructIndexOf( t *testing.T  ) {
	t.Log("Index of DamStruct")
	vl := SliceDamStruct{}
	
	v2 := new(DamStruct)
	vl.Add(new(DamStruct))
	vl.Add(v2)
	vl.Add(new(DamStruct))

	v := vl.IndexOf(v2)
	t.Log("IndexOf v2 ==",v)	
	
	if v != 1 {
		t.FailNow()
	}
}
