package help

import (
	"fmt"
	"strings"

	"github.com/ktye/iv/apl"
)

type help struct{}

func (_ help) String(a *apl.Apl) string {
	return "help"
}

func (_ help) Call(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	filter := ""
	if l != nil {
		if s, ok := l.(apl.String); ok == false {
			return usage()
		} else {
			filter = string(s)
		}
	}

	key := ""
	if filter == "" {
		if s, ok := r.(apl.String); ok == true {
			key = string(s)
		}
	}

	if key == "" && filter == "" {
		return overview(a)
	}

	if filter != "" {
		return query(a, filter)
	}
	return apl.String(a.GetDoc(key)), nil
}

func overview(a *apl.Apl) (apl.Value, error) {
	var buf strings.Builder
	keys := a.GetDocKeys()
	for _, q := range keys {
		doc := a.GetDoc(q)
		fmt.Fprintf(&buf, "%s\n", firstline(doc))
	}

	list := func(l []string) {
		n := 0
		for _, s := range l {
			n++
			sep := " "
			if n > 30 {
				sep = "\n"
				n = 0
			}
			fmt.Fprintf(&buf, "%s%s", s, sep)
		}
		fmt.Fprintf(&buf, "\n")
	}

	fmt.Fprintf(&buf, "\nPrimitives:\n")
	list(a.ListAllPrimitives())
	fmt.Fprintf(&buf, "\nOperators:\n")
	list(a.ListAllOperators())

	fmt.Fprintf(&buf, "\nhelp 0  |  help \"key\"   |   \"query\" help 0")

	return apl.String(buf.String()), nil
}

func query(a *apl.Apl, filter string) (apl.Value, error) {
	var buf strings.Builder
	keys := a.GetDocKeys()
	for _, q := range keys {
		doc := a.GetDoc(q)
		if strings.Contains(firstline(doc), filter) {
			fmt.Fprintf(&buf, "%s\n", doc)
		}
	}
	return apl.String(buf.String()), nil
}

func firstline(s string) string {
	if idx := strings.IndexByte(s, '\n'); idx > 0 {
		return s[:idx]
	}
	return s
}

func usage() (apl.Value, error) {
	return nil, fmt.Errorf(`wrong usage of help
	
help 0         - overview
help "key"     - help text for a single key
"query" help 0 - show help entries whos first line matches the query
`)
}
