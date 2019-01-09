package iv

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestTabularText(t *testing.T) {
	input := `1  2 3 4
	  5 6 7
 
8 9

x`
	expect := "1\n2\n3\n4\n\n5\n6\n7\n\n\n8\n9\n\n\nx\n\n"

	r := strings.NewReader(input)
	b, err := ioutil.ReadAll(tabularText(r))
	if err != nil {
		t.Fatal(err)
	}
	if got := string(b); got != expect {
		t.Fatalf("expected:\n%q\ngot:\n%q\n", expect, got)
	}

}
