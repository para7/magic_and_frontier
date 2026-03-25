package main

import (
	"fmt"
	"io"
	"os"

	"maf_command_editor/app/cli"
)

func main() {
	cfg := cli.LoadConfig()
	args := os.Args[1:]
	// if len(args) == 0 {
	// 	os.Exit(runEditor(nil, cfg))
	// }

	switch args[0] {
	case "editor":
		os.Exit(cli.Editor(cfg))
	case "validate":
		os.Exit(cli.Validate(cfg))
	case "export":
		os.Exit(cli.Export(cfg))
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		printUsage(os.Stderr)
		os.Exit(2)
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: mcg <command>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "commands:")
	fmt.Fprintln(w, "  editor     start web editor server")
	fmt.Fprintln(w, "  validate   validate savedata and export settings")
	fmt.Fprintln(w, "  export     validate and export datapack")
}
