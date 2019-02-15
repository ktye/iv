# iv - APL stream processor

Program *iv* is similar to `cmd/apl` but reads data from stdin instead of program text.

The program is given as command line arguments.

```
Usage
	cat data | iv COMMANDS
```

## streaming data
Iv provides two pre-defined functions that read data from stdin:
```
	r ← {<⍤⍵ io→r 0}
	s ← {⍵⍴<⍤0 io→r 0}
```
io→r 0 reads from stdin are returns a channel that provides a line of input on each read.

The rank operator `⍤` is extended to parse sub-arrays from a channel delivering strings.
See `ScanRankArray` in `apl/fmt.go`.
The number and array syntax is not as strict as usual program text to allow interoperatiliby with other programs.

Function `r` is called monadically with a rank argument: `r 0` reads a scalar at a time, `r 1` a line at a time and so on.
Higher order arrays are terminated by multiple newlines, or bracket notation can be used (like json arrays), or matlab style.

Function `s` ignores the structure of incoming data and always reads a scalar at a time, reshaping it according to it's right argument.

## examples
To apply a function on each 2d subarray of the input stream, we can call iv with:
```
	cat data | iv f¨r 2
```


Format each 3 2 subarray of the input stream to a json string:
```
	cat data | iv '`json ⍕¨s 2 3'
```

More examples are given in `testdata`.

## extra libraries
cmd/iv includes only the base packages. Even `io` mentioned above only provides a single function to read from stdin.
To use all packages in this repository, cmd/lui can be called to work like iv by including the argument `-i`.
```
	lui -i COMMANDS
```
