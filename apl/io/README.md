# Package io provides input and output streams

Linking it into APL leads to an unsafe system.

Io overloads several *primitive functions*:
```
	< filename        returns a Channel reading from a file
	< 0               returns a Channel reading from stdin
	!`ls              execute program return a channel
	!(`ls`-l)         same with arguments
	`cat!A            same reading input from A (String method) or channel (pipe)
	`file<channel     write to file (TODO)
	`dst<<`src        copy idiom (TODO)
	`log<!`prog       redirection (TODO)
```

## Filesystem operations

A *filename* is a string that starts with a slash.
A *directory* is a filename that ends with a slash.

Filesystems are mounted in the current session with the *mount* function or `/m` command.
A later mounted filesystem may shadow a previous one.
Examples:
```
	/m . /                           ⍝ mount the current working directory to root
	/m "c:/very deep directory" `/w  ⍝ mount a windows directory under /w
	/m `/path/a `/a                  ⍝ mount /path/a to /a
	/m                               ⍝ list mtab
	io→umount `/a                    ⍝ unmout /a
	<`/                              ⍝ list the root directory, similar to unix ls
	/cd                              ⍝ show current os directory
	/cd `dir                         ⍝ change current os directory
	/e<`/file                        ⍝ open the file content in the editor (requires pkg u)
```
