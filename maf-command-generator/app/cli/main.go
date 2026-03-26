package cli

import (
	"fmt"
	master "maf_command_editor/app/domain/master"
	"maf_command_editor/app/files"
)

func Validate(cfg files.MafConfig) int {
	_ = master.NewDBMaster(cfg)

	// entity := grimoire.NewGrimoireEntity(cfg.GrimoireStatePath)
	// if err := entity.Load(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "failed to load grimoire: %v\n", err)
	// 	return 1
	// }
	// fmt.Printf("Loaded %d grimoires\n", len(entity.GetAll()))
	// errs := entity.ValidateAll(nil)
	// for _, e := range errs {
	// 	fmt.Fprintln(os.Stderr, e)
	// }
	// if len(errs) > 0 {
	// 	return 1
	// }
	return 0
}

func Export(cfg files.MafConfig) int {
	fmt.Println("Export")
	return 0
}

func Editor(cfg files.MafConfig) int {
	fmt.Println("Editor")
	return 0
}
