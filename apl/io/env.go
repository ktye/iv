package io

import (
	"fmt"
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

func (e Env) String(a *apl.Apl) string {
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
	return apl.String(os.Getenv(string(s)))
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
