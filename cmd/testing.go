package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ktye/iv/apl"
)

// AplTest is run from go test in cmd/apl and cmd/lui.
func AplTest(newapl func() *apl.Apl) error {
	d, err := os.Open(".")
	if err != nil {
		return err
	}
	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file, ".apl") == false {
			continue
		}
		if err := testAplFile(newapl, file); err != nil {
			return err
		}
	}
	return nil
}

// IvTest is run from go test in cmd/iv and cmd/lui.
func IvTest(newapl func(io.ReadCloser) *apl.Apl) error {
	d, err := os.Open(".")
	if err != nil {
		return err
	}
	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file, ".iv") == false {
			continue
		}
		if err := testIvFile(newapl, file[:len(file)-3]); err != nil {
			return err
		}
	}
	return nil
}

func testAplFile(newapl func() *apl.Apl, file string) error {
	var out bytes.Buffer
	a := newapl()
	a.SetOutput(&out)
	if err := Apl(a, nil, []string{file}); err != nil {
		return compareError(err, file)
	}
	return compareOut(out.Bytes(), file[:len(file)-4])
}
func testIvFile(newapl func(io.ReadCloser) *apl.Apl, file string) error {
	f, err := os.Open(file + ".iv")
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
	if err := Iv(newapl(ioutil.NopCloser(r)), prog[1:], &out); err != nil {
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

func compareError(got error, file string) error {
	name := file[:len(file)-4] + ".err"
	want, err := ioutil.ReadFile(name)
	if err != nil {
		return fmt.Errorf("%s: failed but (%s)\n%s", file, err, got)
	}
	if got := got.Error() + "\n"; got != string(want) {
		return fmt.Errorf("%s: expected:\n%sgot:\n%s", file, want, got)
	}
	return nil
}

func compareOut(got []byte, file string) error {
	name := file + ".out"
	want, err := ioutil.ReadFile(name)
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
	return fmt.Errorf("%s:%d differs (byte %d). Got:\n%s\nWant:\n%s", name, line, at+1, string(got), string(want))
}
