package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/zabio3/transcriber/transcribe"
	"github.com/zabio3/transcriber/transcribe/gcp"
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
	lang     string
	filePath string
)

// Run ...
func (cli *CLI) Run(args []string) int {
	flags := flag.NewFlagSet("transcriber", flag.ContinueOnError)

	flags.StringVar(&mode, "m", "file", "Use audio recognition mode (\"file\", \"stream\") ")
	flags.StringVar(&filePath, "f", "", "Path to audio file when \"file\" mode")
	flags.StringVar(&lang, "l", "ja-JP", "Use to recognize language (ref https://cloud.google.com/speech-to-text/docs/languages)")
	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprint(cli.ErrStream, err)
		return ExitCodeParseFlagsError
	}

	ctx := context.Background()
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

		b, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Fprint(cli.ErrStream, err)
			return ExitCodeInternalError
		}

		res, err := gcp.RecognizeSpeech(ctx, b, signal, lang)
		if err != nil {
			fmt.Fprint(cli.ErrStream, err)
			return ExitCodeInternalError
		}

		for _, v := range res {
			fmt.Fprint(cli.OutStream, fmt.Sprintf("confidence: %f, content: %s", v.Confidence, v.Content))
		}

		//fmt.Fprint(cli.OutStream, res)
		return ExitCodeOK
	// voice stream recognition
	// case "stream":
	default:
		fmt.Fprint(cli.ErrStream, fmt.Errorf("unknown mode: %s", mode))
		return ExitCodeArgsError
	}
}
