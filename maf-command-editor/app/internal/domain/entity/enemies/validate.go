package enemies

import (
	"fmt"
	"time"

	"tools2/app/internal/domain/common"
)

func ValidateSave(input SaveInput, enemySkillIDs, itemIDs, grimoireIDs map[string]struct{}, now time.Time) common.SaveResult[EnemyEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	id := common.RequireNonEmptyID(errs, "id", input.ID)

	normalizedEnemySkillIDs := make([]string, 0, len(input.EnemySkillIDs))
	seen := map[string]struct{}{}
	for i, sid := range input.EnemySkillIDs {
		idv := common.NormalizeText(sid)
		if idv == "" {
			continue
		}
		if _, ok := enemySkillIDs[idv]; !ok {
			errs.Add(fmt.Sprintf("enemySkillIds.%d", i), "Referenced enemy skill does not exist.")
			continue
		}
		if _, exists := seen[idv]; exists {
			continue
		}
		seen[idv] = struct{}{}
		normalizedEnemySkillIDs = append(normalizedEnemySkillIDs, idv)
	}

	equipment := Equipment{
		Mainhand: normalizeEquipmentSlot(errs, "equipment.mainhand", input.Equipment.Mainhand, itemIDs),
		Offhand:  normalizeEquipmentSlot(errs, "equipment.offhand", input.Equipment.Offhand, itemIDs),
		Head:     normalizeEquipmentSlot(errs, "equipment.head", input.Equipment.Head, itemIDs),
		Chest:    normalizeEquipmentSlot(errs, "equipment.chest", input.Equipment.Chest, itemIDs),
		Legs:     normalizeEquipmentSlot(errs, "equipment.legs", input.Equipment.Legs, itemIDs),
		Feet:     normalizeEquipmentSlot(errs, "equipment.feet", input.Equipment.Feet, itemIDs),
	}

	drops := make([]DropRef, 0, len(input.Drops))
	for i, d := range input.Drops {
		kind := common.NormalizeText(d.Kind)
		refID := common.NormalizeText(d.RefID)
		if refID != "" {
			switch kind {
			case "item":
				if _, ok := itemIDs[refID]; !ok {
					errs.Add(fmt.Sprintf("drops.%d.refId", i), "Referenced entry does not exist.")
				}
			case "grimoire":
				if _, ok := grimoireIDs[refID]; !ok {
					errs.Add(fmt.Sprintf("drops.%d.refId", i), "Referenced entry does not exist.")
				}
			case "minecraft_item":
				if !common.IsNamespacedResourceID(refID) {
					errs.Add(fmt.Sprintf("drops.%d.refId", i), "Must be a minecraft item id.")
				}
			}
		}
		if d.CountMin != nil && d.CountMax != nil && *d.CountMin > *d.CountMax {
			errs.Add(fmt.Sprintf("drops.%d.countMin", i), "Must be <= countMax.")
		}
		if errs[fmt.Sprintf("drops.%d.kind", i)] != "" || errs[fmt.Sprintf("drops.%d.refId", i)] != "" || errs[fmt.Sprintf("drops.%d.weight", i)] != "" {
			continue
		}
		drops = append(drops, DropRef{
			Kind:     kind,
			RefID:    refID,
			Weight:   d.Weight,
			CountMin: d.CountMin,
			CountMax: d.CountMax,
		})
	}

	if errs.Any() {
		return common.SaveValidationError[EnemyEntry](errs, "Validation failed. Fix the highlighted fields.")
	}

	entry := EnemyEntry{
		ID:            id,
		MobType:       common.NormalizeText(input.MobType),
		Name:          common.OptionalText(input.Name),
		HP:            input.HP,
		Memo:          common.OptionalText(input.Memo),
		Attack:        input.Attack,
		Defense:       input.Defense,
		MoveSpeed:     input.MoveSpeed,
		Equipment:     equipment,
		EnemySkillIDs: normalizedEnemySkillIDs,
		DropMode:      common.NormalizeText(input.DropMode),
		Drops:         drops,
		UpdatedAt:     now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}

func normalizeEquipmentSlot(errs common.FieldErrors, field string, slot *EquipmentSlot, itemIDs map[string]struct{}) *EquipmentSlot {
	if slot == nil {
		return nil
	}
	kind := common.NormalizeText(slot.Kind)
	refID := common.NormalizeText(slot.RefID)
	switch kind {
	case "item":
		if _, ok := itemIDs[refID]; !ok {
			errs.Add(field+".refId", "Referenced entry does not exist.")
		}
	case "minecraft_item":
		if !common.IsNamespacedResourceID(refID) {
			errs.Add(field+".refId", "Must be a minecraft item id.")
		}
	default:
		errs.Add(field+".kind", "Invalid value.")
	}
	if slot.Count < 1 || slot.Count > 64 {
		errs.Add(field+".count", "Must satisfy range 1..64.")
	}
	if slot.DropChance != nil && (*slot.DropChance < 0 || *slot.DropChance > 1) {
		errs.Add(field+".dropChance", "Must satisfy range 0..1.")
	}
	if errs[field+".kind"] != "" || errs[field+".refId"] != "" || errs[field+".count"] != "" || errs[field+".dropChance"] != "" {
		return nil
	}
	return &EquipmentSlot{
		Kind:       kind,
		RefID:      refID,
		Count:      slot.Count,
		DropChance: slot.DropChance,
	}
}
