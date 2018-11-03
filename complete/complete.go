package complete

import (
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/funcs"
	"github.com/ktye/iv/apl/operators"
)

var apldoc *apl.Apl

// TODO: reuse code from split for Line as

// Line completes the current line.
// The completion mechanism is misused in two ways:
//	1. It does not complete, but replace the word started
//	2. It provides help text without having to leave the command
//	   that is currently being assembled.
//
// 	- By pressing TAB after a blank, a list of all symbols appears.
// 	- Pressing TAB replaces the characters just before it. E.g.
//		rho<TAB> is replaced with ⍴
//   	It ignores any characters that are not [a-z] before the word, such
//   	that no spaces are neccessary
//		[2 3]rho<TAB> should work.
//	- Pressing TAB after a symbol shows help text for this symbol, e.g.:
//		⍴<TAB>
//
func Line(line string) []string {

	// Print the list of all commands if there is nothing to complete.
	if line == "" {
		return list()
	}

	// If the last rune of the word is a known symbol, we show the help text.
	// Iterate to find the last rune.
	r := ""
	for _, c := range line {
		r = string(c)
	}
	if help := apldoc.GetDoc(r); help != "" {
		return strings.Split(help, "\n")
	}

	// We look for the [a-z] tail part of it but keep the prefix.
	suffix := ""
	idx := 0
	for i := len(line) - 1; i >= 0; i-- {
		if line[i] >= 'a' && line[i] <= 'z' {
			suffix = string(line[i]) + suffix
			idx = i
		} else {
			break
		}
	}

	// If there is a single match, replace by the symbol itself.
	// If there are multiple, show the name as well.
	prefix := line[:idx]
	type match struct{ short, long string }
	var matches []match
	if len(suffix) > 0 {
		for _, e := range Tab {
			if strings.HasPrefix(e.Name, suffix) {
				matches = append(matches, match{prefix + e.Symbol, prefix + e.Name + e.Symbol})
			}
		}
	}
	if len(matches) == 1 {
		return []string{matches[0].short}
	} else if len(matches) > 1 {
		l := make([]string, len(matches))
		for i, v := range matches {
			l[i] = v.long
		}
		return l
	}

	// If we cannot complete, show the list.
	return list()
}

// LinerWords is used by peterh/liner as a WordCompleter.
func LinerWords(line string, pos int) (head string, completions []string, tail string) {
	// We ignore the tail part and assume the cursor is at the end of the line.
	head, completions = split(line)
	return head, completions, ""
}

func split(line string) (prefix string, matches []string) {
	if line == "" {
		return "", list()
	}

	// If the last rune of the word is a known symbol, we show the help text.
	// Iterate to find the last rune.
	r := ""
	for _, c := range line {
		r = string(c)
	}
	if help := apldoc.GetDoc(r); help != "" {
		return line, strings.Split(help, "\n")
	}

	// We look for the [a-z] tail part of it but keep the prefix.
	suffix := ""
	idx := len(line)
	for i := len(line) - 1; i >= 0; i-- {
		if line[i] >= 'a' && line[i] <= 'z' {
			suffix = string(line[i]) + suffix
			idx = i
		} else {
			break
		}
	}

	prefix = line[:idx]
	if len(suffix) == 0 {
		matches = list()
		return prefix, matches
	}

	type match struct{ short, long string }
	var matchlist []match
	if len(suffix) > 0 {
		for _, e := range Tab {
			if strings.HasPrefix(e.Name, suffix) {
				matchlist = append(matchlist, match{e.Symbol, e.Name + e.Symbol})
			}
		}
	}
	if len(matchlist) == 0 {
		return prefix, list()
	} else if len(matchlist) == 1 {
		return prefix, []string{matchlist[0].short}
	} else {
		for _, v := range matchlist {
			matches = append(matches, v.long)
		}
		return prefix, matches
	}
}

func list() []string {
	l := make([]string, len(Tab)+1)
	l[0] = "Symbols"
	for i, e := range Tab {
		l[1+i] = e.Name + " " + e.Symbol
	}
	return l
}

func init() {
	apldoc = apl.New(nil)
	operators.Register(apldoc)
	funcs.Register(apldoc)
}
