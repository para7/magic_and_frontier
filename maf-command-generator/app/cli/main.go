package cli

import (
	"fmt"
	"os"

	"maf_command_editor/app/domain/model/grimoire"
)

func Validate(cfg MafConfig) int {
	store := grimoire.NewStore(cfg.GrimoireStatePath)
	if err := store.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load grimoire: %v\n", err)
		return 1
	}
	fmt.Printf("Loaded %d grimoires\n", len(store.Entries))
	errs := store.ValidateAll()
	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e)
	}
	if len(errs) > 0 {
		return 1
	}
	return 0
}

func Export(cfg MafConfig) int {
	fmt.Println("Export")
	return 0
}

func Editor(cfg MafConfig) int {
	fmt.Println("Editor")
	return 0
}
