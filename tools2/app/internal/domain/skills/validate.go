package skills

import (
	"time"

	"tools2/app/internal/domain/common"
)

func ValidateSave(input SaveInput, itemIDs map[string]struct{}, now time.Time) common.SaveResult[SkillEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	itemID := common.NormalizeText(input.ItemID)
	if itemID != "" {
		if _, ok := itemIDs[itemID]; !ok {
			errs.Add("itemId", "Referenced item does not exist.")
		}
	}
	if errs.Any() {
		return common.SaveValidationError[SkillEntry](errs, "Validation failed. Fix the highlighted fields.")
	}
	entry := SkillEntry{
		ID:        common.NormalizeText(input.ID),
		Name:      common.NormalizeText(input.Name),
		Script:    common.NormalizeText(input.Script),
		ItemID:    itemID,
		UpdatedAt: now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}
