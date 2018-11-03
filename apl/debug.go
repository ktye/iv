package apl

import (
	"fmt"
	"strings"

	"github.com/ktye/iv/apl/scan"
)

// Set ddebug to true to print a call stack of the parser.
const ddebug = false

func enter(s string, t scan.Token) {
	if ddebug {
		fmt.Printf("%s%s> %s\n", indent(), s, t)
		level++
	}
}

func leave(s string) {
	if ddebug {
		level--
		fmt.Printf("%s%s<\n", indent(), s)
	}
}

func indent() string {
	return strings.Repeat("    ", level)
}

var level int
