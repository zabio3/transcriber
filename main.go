package main

import (
	"os"

	"github.com/zabio3/transcriber/cmd"
)

func main() {
	cli := &cmd.CLI{OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
