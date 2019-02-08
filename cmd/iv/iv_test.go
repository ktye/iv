package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIv(t *testing.T) {
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
		if strings.HasSuffix(file, ".iv") == false {
			continue
		}

		if err := testfile(file); err != nil {
			t.Fatal(err)
		}
	}
}

func testfile(file string) error {
	f, err := os.Open(filepath.Join("testdata", file))
	if err != nil {
		return err
	}
	defer f.Close()

	// First line in each iv test file is the program with a comment.
	r := bufio.NewReader(f)
	prog, err := readline(r)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	stdin = ioutil.NopCloser(r)
	if err := iv(prog[1:], &out); err != nil {
		return err
	}

	return compareOut(out.Bytes(), file)
}

func readline(s io.RuneScanner) (string, error) {
	var b strings.Builder
	for {
		if r, _, err := s.ReadRune(); err != nil {
			return "", err
		} else {
			if r == '\n' {
				return b.String(), nil
			}
			b.WriteRune(r)
		}
	}
}

func compareOut(got []byte, file string) error {
	name := file[:len(file)-3] + ".out"
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
