# APL\iv - core package

This is the core package of the interpreter.

## Usage:

go get github.com/ktye/iv/apl/...

```go
	// import github.com/ktye/iv/apl

	a := apl.New(os.Stdout)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)

	err := a.ParseAndEval("⍳3")
```

The core package and all additional packages in it's subdirectories require only the Go standard library.
Extra packages with external dependencies can be found in `iv/aplextra`.

## Packages
- [a](a/) access to the go runtime
- [big](big/) big numbers as an alternative
- [io](io/) filesystem access
- [rpc](rpc/) remote procedure calls and ipc communication
- [strings](strings/) wrapper of go strings library
- [xgo](xgo/) generic interface to go types
