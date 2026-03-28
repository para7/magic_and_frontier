package main

import (
	"fmt"
	"io"
	"os"

	cli "maf_command_editor/app/cli"
	"maf_command_editor/app/domain/master"
	config "maf_command_editor/app/files"
)

func main() {
	cfg := config.LoadConfig()
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage(os.Stderr)
		os.Exit(2)
	}

	switch args[0] {
	case "editor":
		dmas := master.NewDBMaster(cfg)
		os.Exit(cli.Editor(dmas))
	case "validate":
		dmas := master.NewDBMaster(cfg)
		os.Exit(cli.Validate(dmas))
	case "export":
		dmas := master.NewDBMaster(cfg)
		os.Exit(cli.Export(dmas, cfg))
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
