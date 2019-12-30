package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/zabio3/transcriber/transcribe"
)

// CLI represents CLI interface.
type CLI struct {
	OutStream, ErrStream io.Writer
}

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK = iota + 1
	ExitCodeParseFlagsError
	ExitCodeArgsError
	ExitCodeInternalError
)

var (
	mode     string
	filePath string
)

// Run ...
func (cli *CLI) Run(args []string) int {
	flags := flag.NewFlagSet("transcriber", flag.ContinueOnError)

	flags.StringVar(&mode, "m", "file", "Use audio recognition mode (\"file\" | \"stream\") ")
	flags.StringVar(&filePath, "f", "", "Path to audio file when \"file\" mode")
	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprint(cli.ErrStream, err)
		return ExitCodeParseFlagsError
	}

	switch mode {
	// audio file recognition
	case "file":
		if filePath == "" {
			fmt.Fprint(cli.ErrStream, fmt.Errorf("empty filepath"))
			return ExitCodeArgsError
		}
		signal, err := transcribe.GetSampleRate(filePath)
		if err != nil {
			fmt.Fprint(cli.ErrStream, err)
			return ExitCodeInternalError
		}

		fmt.Fprint(cli.OutStream, signal.Rate, signal.Channels)
	// voice stream recognition
	// case "stream":
	default:
		fmt.Fprint(cli.ErrStream, fmt.Errorf("unknown mode: %s", mode))
		return ExitCodeArgsError
	}

	return ExitCodeOK
}
