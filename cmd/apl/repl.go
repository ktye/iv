package main

import (
	"fmt"

	"github.com/ktye/iv/complete"
	"github.com/peterh/liner"
)

// Repl runs the interactive interpreter which uses liner for console handling.
// The completion mechanism can be used for APL symbol input or help text.
func (opt *options) repl() error {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetTabCompletionStyle(liner.TabPrints)
	line.SetWordCompleter(complete.LinerWords)

	for {

		if s, err := line.Prompt("        "); err == nil {
			if err := opt.state.ParseAndEval(s); err != nil {
				fmt.Fprintf(opt.stderr, "error: %s\n", err)
			}
			line.AppendHistory(s)
		} else if err == liner.ErrPromptAborted {
			return nil
		} else {
			return err
		}
	}
}
