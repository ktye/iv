package io

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/scan"
)

// Mount mounts the filesystem indicated by R to the mount point L.
// If L is nil, it returns the mtab as a dictionary.
//
// R may contain a protocol suffix, such as zip:// that is matched against
// known file systems.
// If no protocol can be matched, R is considered to be an os path.
//
// The special file "." can be used, which is always the current working directory.
func mount(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if L == nil {
		mtab.Lock()
		defer mtab.Unlock()

		d := apl.Dict{}
		for _, t := range mtab.tab {
			name := apl.String(t.mpt)
			d.K = append(d.K, name)
			if d.M == nil {
				d.M = make(map[apl.Value]apl.Value)
			}
			d.M[name] = apl.String(t.src.String())
		}
		return &d, nil
	}
	var mpt, src string
	if s, ok := L.(apl.String); ok == false {
		return nil, fmt.Errorf("io mount: left argument must be a string %T", L)
	} else {
		mpt = string(s)
	}
	if s, ok := R.(apl.String); ok == false {
		return nil, fmt.Errorf("io mount: right argument must be a string %T", R)
	} else {
		src = string(s)
	}

	// Test if the filesystem matches a registerd protocol.
	for name, f := range protocols {
		pre := name + "://"
		if strings.HasPrefix(src, pre) {
			fsys, err := f.FileSystem(strings.TrimPrefix(src, pre))
			if err != nil {
				return nil, err
			}
			if err := Mount(mpt, fsys); err != nil {
				return nil, err
			}
			return apl.EmptyArray{}, nil
		}
	}

	// Special case, "." remains always relative.
	if src == "." {
		if err := Mount(mpt, fs(".")); err != nil {
			return nil, err
		}
		return apl.EmptyArray{}, nil
	}

	// Mount a directory.
	fi, err := os.Stat(src)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() == false {
		return nil, fmt.Errorf("io mount: src is not a directory: %s", src)
	}
	abs, err := filepath.Abs(src)
	if err != nil {
		return nil, err
	}
	if err := Mount(mpt, fs(abs)); err != nil {
		return nil, err
	}
	return apl.EmptyArray{}, nil
}

// Umount removes the mount point R.
func umount(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	s, ok := R.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("io umount: argument must be a string %T", R)
	}
	Umount(string(s))
	return apl.EmptyArray{}, nil
}

// cd changes the working directory.
// If R is empty, it returns the current directory.
func cd(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	s, ok := R.(apl.String)
	if ok == false {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		} else {
			return apl.String(dir), nil
		}
	}
	if err := os.Chdir(string(s)); err != nil {
		return nil, err
	}
	return apl.EmptyArray{}, nil
}

func mCmd(t []scan.Token) []scan.Token {
	if len(t) == 0 {
		// List mtab.
		return []scan.Token{
			scan.Token{T: scan.Identifier, S: "io→mount"},
			scan.Token{T: scan.Number, S: "0"},
		}
	}
	if len(t) < 2 {
		return t // one argument is an error.
	}

	// Replace . / of the next two tokens to strings.
	for i := 0; i < 2; i++ {
		if t[i].T == scan.Symbol && (t[i].S == "/" || t[i].S == ".") {
			t[i].T = scan.String
		}
	}
	tokens := []scan.Token{t[1], scan.Token{T: scan.Identifier, S: "io→mount"}, t[0]}
	return append(tokens, t[2:]...)
}

func cdCmd(t []scan.Token) []scan.Token {
	cdt := scan.Token{T: scan.Identifier, S: "io→cd"}
	if len(t) == 0 {
		return []scan.Token{cdt, scan.Token{T: scan.Number, S: "0"}}
	}
	return append([]scan.Token{cdt}, t...)
}

type toCommand func([]scan.Token) []scan.Token

func (f toCommand) Rewrite(t []scan.Token) []scan.Token {
	return f(t)
}
