package shell

import "github.com/ktye/iv/apl"

// Register adds functions ls, cd and cat global function variables.
//	ls 0
//	ls "dir"
//	ls "*"
//	cd 0
//	cd "path/to/file"
//	cat/"file"
func Register(a *apl.Apl) {
	must(a.Assign("ls", ls{}))
	//must(a.Assign("cd", cd{}))
	//must(a.Assign("cat", cat{}))
	a.RegisterDoc("ls", `Z←ls R  list directory
Z←ls R  R string: Z: array of directory names, shape N 1.
Z←ls R  R "*": Z: array of directory names (recursive), shape N 1.
Z←ls R  R not a string: same as R←"." (current directory)
`)
	a.RegisterDoc("cd", `Z←cd R  change working directory
Z←cd R  R string: Z: new directory name, changes directory.
Z←cd R  R not a string: Z: name of current directory (string)
`)
	a.RegisterDoc("cat", `Z←cat R  return file content
Z←cat R  R string: Z: array of strings, shape (number of lines) 1.
`)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
