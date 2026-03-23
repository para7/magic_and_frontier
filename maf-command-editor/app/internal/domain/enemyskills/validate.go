package enemyskills

import (
	"time"

	"tools2/app/internal/domain/common"
)

func ValidateSave(input SaveInput, now time.Time) common.SaveResult[EnemySkillEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	id := common.RequireNonEmptyID(errs, "id", input.ID)
	if errs.Any() {
		return common.SaveValidationError[EnemySkillEntry](errs, "Validation failed. Fix the highlighted fields.")
	}
	entry := EnemySkillEntry{
		ID:          id,
		Name:        common.OptionalText(input.Name),
		Description: common.OptionalText(input.Description),
		Script:      common.NormalizeText(input.Script),
		UpdatedAt:   now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}
