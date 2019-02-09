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

## streaming data
The interpreter has some built-in methods for parsing textual data into values, mostly with dyadic ⍎.

But also the rank operator `⍤` is extended to parse sub-arrays from a channel containing strings, e.g. `<0` (stdin).

This example reads lines of text from stdin and interpretes them as 2 dimensional arrays: `<⍤2<0`

As this is a common pattern for the use case of iv, it is stored as a lambda function in r: `r←{<⍤⍵<0}`.

To apply a function on each 2d subarray of the input stream, we can call iv with:
```
	cat data | iv f¨r 2
```

Another use case is to reshape the input stream if it does not contain a structure.
This can be done with the extension of ⍴ for stream.
It is already stored in the variable s: `s←{⍵⍴<⍤0<0}`.

This example formats each 3 2 subarray of the input stream to a json string:
```
	cat data | iv '`json ⍕¨s 2 3'
```


