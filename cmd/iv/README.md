# iv - APL stream processor

Program *iv* is similar to `cmd/apl` but reads data from stdin instead of program text.

The program is given as command line arguments.

```
Usage
	cat data | iv COMMANDS
```

Monadic `<0` is defined to return a Channel that read lines of input from stdin.
Otherwise only the standard packages are included.

Monadic `<"FILE"` reads lines from a file.
To simplify reading a library, there is the l command.
These two lines are equivalent:
```
	cat data | iv '⍎¨<`FILE ⋄ COMMANDS'
	cat data | iv '/l`FILE COMMANDS'
```