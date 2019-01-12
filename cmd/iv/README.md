# iv stream processor

## Rewrite
Iv is currently being rewritten and simplified.

## Intro
Iv is a command-line program similar to awk, that operates on a text stream of n-dimensional numeric data using APL syntax.

It reads a possibly large array on stdin and applies a lambda function repeatedly on a subspace of a specified rank.

The input array can be considered to have shape `[* ... * A B C]`, when the program is called with the
a rank parameter of 3.
The number and size of leading dimensions (indicated by stars) is unknown to the program during execution.

Consider the case, when the input is a very large table and `iv` is called with `rank 1`.
Then the complete data has the shape [* COLS], but the lambda is executed on each row of shape [COLS].

The known axes, are not given on the command line, they are part of the input data.
Only the rank argument at which the data should be cut into pieces must be told.

## Termination level
If the input data consists of multiple tables, we could call `iv` with `rank 2` and operate on each table individually.
Alternatively, we can also use `rank 1` and operate on vectors.
In this case it would be good to know, if a table is complete within the stream. E.g. to print a summary.

This is reported as the *termination level* ⍺ to the lambda function,
which is normally 0 and increased by 1 for each higher dimension that is completing when the lambda function is called.

In the multi-table example above, the lambda function receives ⍺←0 for each line, and ⍺←1 at the end of each table.

```
	$ cat data | iv -r2 -b"BEGIN BLOCK" -eEND f:⍺g⍵ ⋄ g:B⍵
	f, g, h are placeholders.
	This translates into APL:
	
	BEGIN BLOCK	    ⍝ Any statement (optional)
	iv←{f:⍺g⍵ ⋄ g:h⍵}  ⍝ ⍵ is the rank 2 sub-array, ⍺ the termination level
	IvC←iv→r 2          ⍝ IvC is a channel.
	⍝ iv→r is the function r in package iv. 
	⍝ It is called with rank 2 given on the command line as -r2.
	⍝ Each take on the channel returns a list (A;E;) with the array and termination level.
TODO update	{(¯1↑⍵) iv 1↑⍵}¨IvC
	END BLOCK           ⍝ optinal end statement, given as -e argument
	
```
TODO: explain termination levels

## Input data
The input data is a stream of numbers in text form, split at whitespace.
The last axis has blanks or tabs as separators (multiples are ignored) and higher dimsions are split by newline.


**Note**: When working on Windows with msys2, arguments such as '+/' may be automagically
translated by a hidden compatibility layer into a path.
To prevent this export `MSYS2_ARG_CONV_EXCL="*"`.
