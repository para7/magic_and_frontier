package skills

import (
	"time"

	"tools2/app/internal/domain/common"
)

var validSkillTypes = map[string]struct{}{
	"sword": {},
	"bow":   {},
	"axe":   {},
}

func ValidateSave(input SaveInput, now time.Time) common.SaveResult[SkillEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	id := common.RequireNonEmptyID(errs, "id", input.ID)
	skillType := common.OptionalText(input.SkillType)
	if skillType == "" {
		skillType = "sword"
	}
	if _, ok := validSkillTypes[skillType]; !ok {
		errs.Add("skilltype", "Must be one of sword, bow, axe.")
	}
	if errs.Any() {
		return common.SaveValidationError[SkillEntry](errs, "Validation failed. Fix the highlighted fields.")
	}
	entry := SkillEntry{
		ID:          id,
		Name:        common.OptionalText(input.Name),
		SkillType:   skillType,
		Description: common.OptionalText(input.Description),
		Script:      common.NormalizeText(input.Script),
		UpdatedAt:   now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}
