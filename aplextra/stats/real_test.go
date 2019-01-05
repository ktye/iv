package stats

import "testing"

func TestReal(t *testing.T) {
	var r Real
	r.Push(1)
	r.Push(2)
	r.Push(3)
	if r.Num != 3 || r.Mean != 2 || r.Min != 1 || r.Max != 3 || r.Std() != 1 {
		t.Fail()
	}
}
