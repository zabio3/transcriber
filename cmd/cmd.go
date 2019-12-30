package cmd

import "io"

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK = iota + 1
	ExitCodeParseFlagsError
)

// CLI represents CLI interface.
type CLI struct {
	ErrStream io.Writer
}

// Run ...
func (cli *CLI) Run(args []string) int {
	return ExitCodeOK
}
