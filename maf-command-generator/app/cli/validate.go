package cli

import (
	"fmt"
	master "maf_command_editor/app/domain/master"
	"maf_command_editor/app/files"
	"os"
)

func Validate(cfg files.MafConfig) int {
	db := master.NewDBMaster(cfg)

	errs := db.ValidateAll()
	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "error: %v\n", e)
	}
	if len(errs) > 0 {
		return 1
	}
	return 0

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
