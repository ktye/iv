package iv

import (
	"strings"
	"testing"
)

func TestFile(t *testing.T) {
	type testCase struct {
		cond bool
		msg  string
	}

	f, err := ParseFile(strings.NewReader(file1), "file1", 0, false, "")
	tests := []testCase{
		{err == nil, "err==nil"},
		{f.Begin == "let's start", "begin is wrong"},
		{len(f.Rules) == 1, "number of rules"},
		{f.Rules[0][0] == "x>0", "condition of rule#1 is wrong"},
		{f.Rules[0][1] == "this is the first rule", "rule#1 is wrong"},
		{f.Uniform == true, "uniform should be true"},
	}
	for _, tc := range tests {
		if tc.cond == false {
			t.Fatalf("%s", tc.msg)
		}
	}

	f, err = ParseFile(strings.NewReader(file2), "file2", 0, false, "")
	tests = []testCase{
		{err == nil, "err==nil"},
		{f.Rank == 5, "rank must be 5"},
		{f.Uniform == false, "uniform should be false"},
		{f.Begin == "\nthis is\nthe begin\nblock", "begin is wrong"},
		{len(f.Rules) == 2, "number of rules"},
		{f.Rules[0][0] == "1", "condition of rule#1 is wrong"},
		{f.Rules[0][1] == "rule number one", "rule#1 is wrong"},
		{f.Rules[1][0] == "2", "condition of rule#2 is wrong"},
		{f.Rules[1][1] == "\n\trule number two\n\thas two lines", "rule#2 is wrong"},
	}
	for _, tc := range tests {
		if tc.cond == false {
			t.Fatalf("%s", tc.msg)
		}
	}
}

const file1 = `
# Strip comment
BEGIN let's start
UNIFORM
x>0:this is the first rule
`

const file2 = `
RANK 5
BEGIN
this is
the begin
block

1:	rule number one

2:
	rule number two
	has two lines
`
