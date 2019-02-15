package main

import (
	"io"
	"os"
	"testing"

	"github.com/ktye/iv/apl"
	aio "github.com/ktye/iv/apl/io"
	"github.com/ktye/iv/cmd"
)

func TestLui(t *testing.T) {
	if err := os.Chdir("../apl/testdata"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.AplTest(newApl); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir("../../iv/testdata"); err != nil {
		t.Fatal(err)
	}
	newapl := func(f func() *apl.Apl) func(io.ReadCloser) *apl.Apl {
		return func(in io.ReadCloser) *apl.Apl {
			aio.Stdin = in
			return f()
		}
	}
	if err := cmd.IvTest(newapl(newApl)); err != nil {
		t.Fatal(err)
	}
}
