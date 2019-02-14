package main

import (
	"os"
	"testing"

	"github.com/ktye/iv/cmd"
)

func init() {
	if err := os.Chdir("testdata"); err != nil {
		panic(err)
	}
}

func TestIv(t *testing.T) {
	if err := cmd.IvTest(newApl); err != nil {
		t.Fatal(err)
	}
}
