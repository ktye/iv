# extra apl packages

This directory contains extra packages for apl written in go.

In contrast to `../apl/*` which only depends on the go standard library, these packages may contain lots of dependencies.

To be included in the interpreter, each package needs to be compiled in:
```go
	/*
	include "github.com/ktye/iv/apl"
	include "github.com/ktye/iv/apl/numbers"
	include "github.com/ktye/iv/apl/primitives"
	include "github.com/ktye/iv/apl/operators"
	include "github.com/ktye/iv/aplextra/q"
	*/
	
	// This example uses the standard interpreter with the extra q package.
	a := apl.New(os.Stdout)
	numbers.Register(a)
	operators.Register(a)
	q.Register(a)
```

## extra packages
- q: rpc interface from APL to kdb (via `github.com/sv/kdbgo`)
- u: ui elements used by `cmd/lui`

### planned:
- [ ] mat: linear algebra package based on gonum
- [ ] plot: 
- [ ] stats: statistics package for data streams
- [ ] images/animations should be possible with the std lib.
