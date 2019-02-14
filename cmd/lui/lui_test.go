package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/cmd"
)

func TestLui(t *testing.T) {
	if err := os.Chdir("../apl/testdata"); err != nil {
		t.Fatal(err)
	}
	fmt.Println(os.Getwd())
	if err := cmd.AplTest(newApl); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir("../../iv/testdata"); err != nil {
		t.Fatal(err)
	}
	fmt.Println(os.Getwd())
	newapl := func(f func() *apl.Apl) func(io.ReadCloser) *apl.Apl {
		return func(io.ReadCloser) *apl.Apl {
			return f()
		}
	}
	if err := cmd.IvTest(newapl(newApl)); err != nil {
		t.Fatal(err)
	}
}
