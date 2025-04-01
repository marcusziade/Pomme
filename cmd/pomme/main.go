package main

import (
	"fmt"
	"os"

	"github.com/marcusziade/pomme/cmd/pomme/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
