package items

import (
	"time"

	"maf-command-editor/app/internal/domain/common"
)

func ValidateSave(input SaveInput, skillIDs map[string]struct{}, now time.Time) common.SaveResult[ItemEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	id := common.RequireNonEmptyID(errs, "id", input.ID)
	skillID := common.OptionalText(input.SkillID)
	if skillID != "" {
		if _, ok := skillIDs[skillID]; !ok {
			errs.Add("skillId", "Referenced skill does not exist.")
		}
	}
	nbt, enchantErr := buildItemNBT(input)
	if enchantErr != "" {
		errs.Add("enchantments", enchantErr)
	}
	if errs.Any() {
		return common.SaveValidationError[ItemEntry](errs, "Validation failed. Fix the highlighted fields.")
	}
	entry := ItemEntry{
		ID:                  id,
		ItemID:              common.NormalizeText(input.ItemID),
		SkillID:             skillID,
		CustomName:          common.OptionalText(input.CustomName),
		Lore:                common.OptionalText(input.Lore),
		Enchantments:        common.OptionalText(input.Enchantments),
		Unbreakable:         input.Unbreakable,
		CustomModelData:     common.OptionalText(input.CustomModelData),
		RepairCost:          common.OptionalText(input.RepairCost),
		HideFlags:           common.OptionalText(input.HideFlags),
		PotionID:            common.OptionalText(input.PotionID),
		CustomPotionColor:   common.OptionalText(input.CustomPotionColor),
		CustomPotionEffects: common.OptionalText(input.CustomPotionEffects),
		AttributeModifiers:  common.OptionalText(input.AttributeModifiers),
		CustomNBT:           common.OptionalText(input.CustomNBT),
		NBT:                 nbt,
		UpdatedAt:           now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}
