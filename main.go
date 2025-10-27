package main

import (
	"os"

	"github.com/piojanu/talksh/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
