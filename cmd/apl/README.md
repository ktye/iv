# cmd/apl

Apl is a simple command line program that runs APL\iv.
It includes only the basic packages *numbers*, *primitives* and *operators*.

It is just one example to use the interpreter.
A more advanced program is `cmd/lui`.

## Usage
```
	apl
```
If no input argument is given, the program acts as a simple REPL reading a line at a time.
Multiline statements are errors.
On error it prints a message but continues.

Apl does not contain a readline library. If you need line editing or history, consider `rlwrap`.

```
	apl FILE ...
```
If one or more arguments are given, it reads input from these files.
Multiline statements are allowed.
On error it prints a message to stderr and exits.
If an argument is `-` it reads from stdin, but otherwise behaves like reading from a file.


## Testing
`go test` runs all file in `testdata/*.apl` and compares the results to the corresponding `.out` files.
If a file is known to fail, it's error message is in a `.err` file.

More unit tests are in `apl/primitives/apl_test.go`.

## Example
```
	apl testdata/a.apl
```
