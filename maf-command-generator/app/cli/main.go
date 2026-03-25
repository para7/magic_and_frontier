package cli

import (
	"fmt"
	"maf_command_editor/app/domain/master/masterimpl"
	"maf_command_editor/app/domain/model/grimoire"
)

func Validate(cfg MafConfig) int {
	_ = masterimpl.NewDBMaster(grimoire.NewGrimoireEntity(cfg.GrimoireStatePath))

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

func Export(cfg MafConfig) int {
	fmt.Println("Export")
	return 0
}

func Editor(cfg MafConfig) int {
	fmt.Println("Editor")
	return 0
}
