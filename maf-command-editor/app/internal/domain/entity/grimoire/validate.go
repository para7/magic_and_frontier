package grimoire

import (
	"time"

	"maf-command-editor/app/internal/domain/common"
)

func ValidateSave(input SaveInput, now time.Time) common.SaveResult[GrimoireEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	id := common.RequireNonEmptyID(errs, "id", input.ID)
	if errs.Any() {
		return common.SaveValidationError[GrimoireEntry](errs, "Validation failed. Fix the highlighted fields.")
	}
	entry := GrimoireEntry{
		ID:          id,
		CastID:      input.CastID,
		CastTime:    input.CastTime,
		MPCost:      input.MPCost,
		Script:      common.NormalizeText(input.Script),
		Title:       common.NormalizeText(input.Title),
		Description: common.OptionalText(input.Description),
		UpdatedAt:   now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}
