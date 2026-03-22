package skills

import (
	"time"

	"tools2/app/internal/domain/common"
)

func ValidateSave(input SaveInput, now time.Time) common.SaveResult[SkillEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	id := common.RequirePrefixedSequenceID(errs, "id", input.ID, "skill_")
	if errs.Any() {
		return common.SaveValidationError[SkillEntry](errs, "Validation failed. Fix the highlighted fields.")
	}
	entry := SkillEntry{
		ID:          id,
		Name:        common.OptionalText(input.Name),
		Description: common.OptionalText(input.Description),
		Script:      common.NormalizeText(input.Script),
		UpdatedAt:   now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}
