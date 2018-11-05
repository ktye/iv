# extra apl packages

This directory contains extra packages for apl.
They may define apl.Values, add or overload primitive functions, add operators,
add functions as workspace variables, or simply pollute the workspace.

To include it to the interpreter, register it:
```go
	// This example includes the image package.
	include "github.com/ktye/iv/apl"
	include "github.com/ktye/iv/apl/funcs"
	include "github.com/ktye/iv/apl/operators"
	include "github.com/ktye/iv/aplextra/image"
	a := apl.New(os.Stdout)
	funcs.Register(a)
	operators.Register(a)
	image.Register(a)
```