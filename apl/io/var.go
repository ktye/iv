package io

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/ktye/iv/apl"
)

// Varfs exports variables as a file system.
type varfs struct {
	*apl.Apl
}

func (v varfs) FileSystem(root string) (FileSystem, error) {
	if v.Apl == nil {
		return nil, fmt.Errorf("varfs: apl is not connected")
	}
	if root != "/" {
		return nil, fmt.Errorf("varfs can only be registerd with root file: var:///")
	}
	return v, nil
}

func (v varfs) String() string {
	return "var:///"
}

func (v varfs) Open(name, mpt string) (io.ReadCloser, error) {
	list := func() (io.ReadCloser, error) {
		pkg := strings.TrimSuffix(name, "/")
		l, err := v.Apl.Vars(pkg)
		if err != nil {
			return nil, err
		}
		if pkg != "" {
			pkg += "→"
		}
		var buf bytes.Buffer
		tw := tabwriter.NewWriter(&buf, 1, 0, 1, ' ', 0)
		for i := range l {
			varname := pkg + l[i]
			info := ""
			if strings.HasSuffix(varname, "/") == false {
				x := v.Apl.Lookup(varname)
				if x == nil {
					info = "?"
				} else {
					info = reflect.TypeOf(x).String()
					if ar, ok := x.(apl.Array); ok {
						info += fmt.Sprintf(" %v", ar.Shape())
					}
					if s, ok := x.(apl.String); ok {
						info += fmt.Sprintf(" %d", len(s))
					}
				}
			}
			fmt.Fprintf(tw, "%s\t%s\n", mpt+varname, info)
		}
		tw.Flush()
		return ioutil.NopCloser(&buf), nil
	}
	if name == "" || strings.HasSuffix(name, "/") {
		return list()
	}

	// Print a variable.
	x := v.Apl.Lookup(name)
	if x != nil {
		return ioutil.NopCloser(strings.NewReader(x.String(v.Apl))), nil
	}

	return nil, &os.PathError{
		Op:   "open",
		Path: name,
		Err:  fmt.Errorf("variable does not exist"),
	}
}
