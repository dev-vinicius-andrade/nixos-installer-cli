package main

import (
	"fmt"
	"os"

	"github.com/dev-vinicius-andrade/nioscli/commands/nioscli"
)

var version = "dsadasd"

func main() {
	//dotfilesInformation := &context.NiOsiContext{}
	command := nioscli.CreateCommand()
	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
