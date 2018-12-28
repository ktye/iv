package shell

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ktye/iv/apl"
)

// ls prints file in the current directory
type ls struct{}

func (l ls) String(a *apl.Apl) string {
	return "ls"
}

func (_ ls) Call(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	if l != nil {
		return nil, fmt.Errorf("ls cannot be called dyadically")
	}

	dir, ok := r.(apl.String)
	if ok == false {
		dir = "."
	}

	if dir == "*" {
		return walk()
	}

	f, err := os.Open(string(dir))
	if err != nil {
		return nil, err
	}

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	ar := apl.MixedArray{
		Dims:   []int{len(names), 1},
		Values: make([]apl.Value, len(names)),
	}
	for i, s := range names {
		ar.Values[i] = apl.String(s)
	}
	return ar, nil
}

func walk() (apl.Value, error) {
	var names []string

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := path + info.Name()
		if info.IsDir() {
			name += "/"
		}
		names = append(names, name)
		return nil
	}

	if err := filepath.Walk(".", walker); err != nil {
		return nil, err
	}

	ar := apl.MixedArray{
		Dims:   []int{len(names), 1},
		Values: make([]apl.Value, len(names)),
	}
	for i, s := range names {
		ar.Values[i] = apl.String(s)
	}
	return ar, nil
}
