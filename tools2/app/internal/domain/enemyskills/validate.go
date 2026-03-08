package enemyskills

import (
	"time"

	"tools2/app/internal/domain/common"
)

func ValidateSave(input SaveInput, now time.Time) common.SaveResult[EnemySkillEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	cooldown := input.Cooldown

	var trigger *Trigger
	if input.Trigger != "" {
		t := Trigger(common.NormalizeText(input.Trigger))
		trigger = &t
	}
	if errs.Any() {
		return common.SaveValidationError[EnemySkillEntry](errs, "Validation failed. Fix the highlighted fields.")
	}
	entry := EnemySkillEntry{
		ID:        common.NormalizeText(input.ID),
		Name:      common.NormalizeText(input.Name),
		Script:    common.NormalizeText(input.Script),
		Cooldown:  cooldown,
		Trigger:   trigger,
		UpdatedAt: now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}
