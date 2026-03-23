package items

import (
	"time"

	"tools2/app/internal/domain/common"
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

func Upsert(state ItemState, entry ItemEntry) (ItemState, common.SaveMode) {
	for i := range state.Items {
		if state.Items[i].ID == entry.ID {
			next := append([]ItemEntry{}, state.Items...)
			next[i] = entry
			return ItemState{Items: next}, common.SaveModeUpdated
		}
	}
	next := append(append([]ItemEntry{}, state.Items...), entry)
	return ItemState{Items: next}, common.SaveModeCreated
}

func Delete(state ItemState, id string) (ItemState, bool) {
	next := make([]ItemEntry, 0, len(state.Items))
	found := false
	for _, it := range state.Items {
		if it.ID == id {
			found = true
			continue
		}
		next = append(next, it)
	}
	return ItemState{Items: next}, found
}
