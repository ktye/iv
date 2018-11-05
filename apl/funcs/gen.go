// +build ignore

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var prefix = regexp.MustCompile(`apl_test.go:[0-9]*:`)
var tabs = regexp.MustCompile(`^\t*`)

// This program is run by go generate.
// It runs go test -v which includes calls to t.Log that show all APL in- and output.
// The output of go test is filtered and written to Tests.md.
func main() {

	cmd := exec.Command("go", "test", "-v")
	testout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	w, err := os.Create("Tests.md")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	fmt.Fprintf(w, "# Test results\n\n```apl\n")

	scn := bufio.NewScanner(testout)
	for scn.Scan() {
		s := scn.Text()
		if strings.HasPrefix(s, "===") {
			continue
		}
		if strings.HasPrefix(s, "---") {
			continue
		}
		if strings.HasPrefix(s, "ok") {
			continue
		}
		s = prefix.ReplaceAllString(s, "")
		s = tabs.ReplaceAllString(s, "")
		fmt.Fprintf(w, "%s\n", s)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "```\n")
}
