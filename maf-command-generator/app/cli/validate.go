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
}
