package cli

import (
	"fmt"
	cv "maf_command_editor/app/domain/custom_validator"
	master "maf_command_editor/app/domain/master"
	"maf_command_editor/app/files"
	"os"
)

func Validate(cfg files.MafConfig) int {
	db := master.NewDBMaster(cfg)

	allErrs := db.ValidateAll()
	total := 0
	for i, recordErrs := range allErrs {
		for _, _e := range recordErrs {
			e := _e
			e.Entity = e.Entity + fmt.Sprintf("[%d]", i+1)
			fmt.Fprintf(os.Stderr, " %s\n", cv.FormatValidationError(e))
			total++
		}
	}

	if total > 0 {
		fmt.Fprintf(os.Stderr, "\nvalidation failed: %d error(s)\n", total)
		return 1
	}
	fmt.Print("validation passed\n")
	return 0
}
