package io

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ktye/iv/apl"
)

// Env returns the environment as an object.
// Assigning to a key changes the environment variable.
func env(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return Env{}, nil
}

type Env struct{}

func (e Env) String(f apl.Format) string {
	v := os.Environ()
	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 1, 0, 1, ' ', 0)
	for _, s := range v {
		s = strings.Replace(s, "=", ":\t", 1)
		fmt.Fprintln(tw, s)
	}
	tw.Flush()
	return b.String()
}
func (e Env) Copy() apl.Value { return e }

func (e Env) Keys() []apl.Value {
	v := os.Environ()
	res := make([]apl.Value, len(v))
	for i, s := range v {
		idx := strings.Index(s, "=")
		if idx > 0 {
			s = s[:idx]
		}
		res[i] = apl.String(s)
	}
	return res
}

func (e Env) At(a *apl.Apl, v apl.Value) apl.Value {
	s, ok := v.(apl.String)
	if ok == false {
		return apl.String("")
	}
	val, ok := os.LookupEnv(string(s))
	if ok {
		return apl.String(val)
	}
	return nil
}

func (e Env) Set(a *apl.Apl, key, val apl.Value) error {
	k, ok := key.(apl.String)
	if ok == false {
		return fmt.Errorf("setenv: key must be a string: %T", key)
	}
	v, ok := val.(apl.String)
	if ok == false {
		return fmt.Errorf("setenv: value must be a string: %T", val)
	}
	return os.Setenv(string(k), string(v))
}

type envfs struct{}

func (e envfs) FileSystem(root string) (FileSystem, error) {
	if root != "/" {
		return nil, fmt.Errorf("envfs can only be registered with root file env:///, not %s", root)
	}
	return envfs{}, nil
}

func (e envfs) String() string { return "env:///" }

func (e envfs) Open(name, mpt string) (io.ReadCloser, error) {
	if name == "" {
		v := os.Environ()
		var buf bytes.Buffer
		tw := tabwriter.NewWriter(&buf, 1, 0, 1, ' ', 0)
		for _, s := range v {
			s = strings.Replace(s, "=", "\t", 1)
			fmt.Fprintln(tw, mpt+s)
		}
		tw.Flush()
		return ioutil.NopCloser(&buf), nil
	}
	v, ok := os.LookupEnv(name)
	if ok {
		return ioutil.NopCloser(strings.NewReader(v)), nil
	}
	return nil, &os.PathError{
		Op:   "open",
		Path: name,
		Err:  fmt.Errorf("environment variable does not exist"),
	}
}

func (e envfs) Write(name string) (io.WriteCloser, error) {
	var b strings.Builder
	return envwriter{name: name, Builder: &b}, nil
}

type envwriter struct {
	name string
	*strings.Builder
}

func (e envwriter) Close() error {
	return os.Setenv(e.name, e.String())
}
