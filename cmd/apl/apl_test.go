package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
)

func TestApl(t *testing.T) {
	d, err := os.Open("testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file, ".apl") == false {
			continue
		}

		if err := testfile(file); err != nil {
			t.Fatal(err)
		}
	}
}

func testfile(file string) error {
	var buf bytes.Buffer
	a := apl.New(&buf)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)

	f, err := os.Open(filepath.Join("testdata", file))
	if err != nil {
		return err
	}
	defer f.Close()

	err = run(a, f, file)
	if err != nil {
		return compareError(err, file)
	}
	return compareOut(buf.Bytes(), file)
}

func compareError(goterr error, file string) error {
	name := file[:len(file)-4] + ".err"
	want, err := ioutil.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		return fmt.Errorf("%s: failed but (%s)\n%s", file, err, goterr)
	}
	if got := goterr.Error() + "\n"; got != string(want) {
		return fmt.Errorf("%s: expected:\n%s, got:\n%s", file, want, got)
	}
	return nil
}

func compareOut(got []byte, file string) error {
	name := file[:len(file)-4] + ".out"
	want, err := ioutil.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		return fmt.Errorf("%s: should have failed, but: %s", file, err)
	}

	at := -1
	max := len(want)
	if len(got) < max {
		max = len(got)
	}
	if len(want) != len(got) {
		at = max
	}

	line := 1
	for i := 0; i < max; i++ {
		if want[i] == '\n' {
			line++
		}
		if got[i] != want[i] {
			at = i
			break
		}
	}

	if at < 0 {
		return nil
	}
	return fmt.Errorf("testdata/%s:%d differs (byte %d). Got:\n%s\nWant:\n%s", name, line, at+1, string(got), string(want))
}
