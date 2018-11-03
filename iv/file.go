package iv

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// File contains the parsed content of an iv program file.
type File struct {
	Rank    int
	Uniform bool
	Begin   string
	Rules   [][2]string
}

// ParseFile parses the program from r.
// The file format is a text file with this format:
// 	# This line is a comment.
//	# A rank may be given which overwrites the
//	# rank given on the command line and the default rank
//  	RANK (int)
//	# Uniform is a bool, which can be enabled by the line:
// 	UNIFORM
//	# An optional begin block, which overwrite the -b option.
// 	BEGIN [begin block]
//	# Following are one or more APL rules: conditional:expression.
//	rule:[rule block]
//	...
// A block may be given on the same line as a single line block,
// or start at the following line and end on an empty line.
func ParseFile(r io.Reader, file string, rank int, uniform bool, begin string) (*File, error) {
	scanner := bufio.NewScanner(r)
	f := File{
		Rank:    rank,
		Uniform: uniform,
		Begin:   begin,
	}
	line := 0
	beginBlock := false
	ruleBlock := false
	for scanner.Scan() {
		t := scanner.Text()
		line++
		if strings.HasPrefix(t, "#") {
			continue
		}
		if strings.HasPrefix(t, "RANK") {
			if n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(t, "RANK"))); err != nil {
				return nil, fmt.Errorf("%s:%d cannot parse rank", file, line)
			} else {
				f.Rank = n
			}
			continue
		}
		if strings.HasPrefix(t, "UNIFORM") {
			f.Uniform = true
			continue
		}
		if strings.HasPrefix(t, "BEGIN") {
			t = strings.TrimSpace(strings.TrimPrefix(t, "BEGIN"))
			if t != "" {
				f.Begin = t
				continue
			}
			beginBlock = true
			continue
		}
		if t == "" {
			beginBlock = false
			ruleBlock = false
			continue
		}
		if beginBlock {
			f.Begin += "\n" + t
			continue
		}
		if ruleBlock {
			r := f.Rules[len(f.Rules)-1]
			r[1] += "\n" + t
			f.Rules[len(f.Rules)-1] = r
			continue
		}
		if idx := strings.Index(t, ":"); idx == -1 {
			return nil, fmt.Errorf("%s:%d: rule has no colon", file, line)
		} else {
			r := [2]string{t[:idx], ""}
			if t = strings.TrimSpace(t[idx+1:]); t != "" {
				r[1] = t
				f.Rules = append(f.Rules, r)
			} else {
				f.Rules = append(f.Rules, r)
				ruleBlock = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &f, nil
}
