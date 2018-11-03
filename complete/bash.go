// Package complete provides command line completion for APL symbols
//

package complete

import "fmt"

//go:generate go run zsh_gen.go

// Bash is used from the shell, if it is about to complete arguments for iv.
//
// To enable completion for iv in bash, use:
//	complete -o nospace -C 'iv -complete-bash' iv
// This uses an external program to find completion for any iv argument.
// The external program is iv itself.
//
// See Line for the completion rules.
func Bash(args []string) {

	if len(args) < 2 {
		fmt.Println("complete-bash-error")
		fmt.Println("bash did not set command line arguments")
		return
	}

	// The word to complete is given on the command line.
	word := args[1]

	// Zsh sends also quotes.
	if len(word) > 0 && (word[0] == '\'' || word[0] == '"') {
		word = word[1:]
	}

	matches := Line(word)
	for _, s := range matches {
		fmt.Println(s)
	}
}
