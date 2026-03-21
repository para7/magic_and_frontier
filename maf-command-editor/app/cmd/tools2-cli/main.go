package main

import (
	"os"

	"tools2/app/internal/cli"
	"tools2/app/internal/config"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, config.Load()))
}
