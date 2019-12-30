package cmd

import (
	"flag"
	"fmt"
	"io"
)

// CLI represents CLI interface.
type CLI struct {
	OutStream, ErrStream io.Writer
}

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK = iota + 1
	ExitCodeParseFlagsError
	ExitCodeModeError
)

var (
	modeStr  string
	filePath string
)

// Run ...
func (cli *CLI) Run(args []string) int {
	flags := flag.NewFlagSet("transcriber", flag.ContinueOnError)

	flags.StringVar(&modeStr, "m", "file", "mode for voice recognition")
	flags.StringVar(&filePath, "f", "", "file path of audio file")
	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprint(cli.ErrStream, err)
		return ExitCodeParseFlagsError
	}

	switch modeStr {
	// audio file recognition
	case "file":

	// voice stream recognition
	// case "stream":
	default:
		fmt.Fprint(cli.ErrStream, fmt.Errorf("unknown mode: %s", modeStr))
		return ExitCodeModeError
	}

	return ExitCodeOK
}
