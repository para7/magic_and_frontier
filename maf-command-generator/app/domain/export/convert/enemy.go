package export_convert

import (
	"fmt"
	"reflect"
	"strings"

	model "maf_command_editor/app/domain/model"
	enemyModel "maf_command_editor/app/domain/model/enemy"
	itemModel "maf_command_editor/app/domain/model/item"
)

func ToEnemyFunctionLines(entry enemyModel.Enemy, lootID string, itemsByID map[string]itemModel.Item) []string {
	return []string{
		fmt.Sprintf("# enemyId=%s mobType=%s", entry.ID, entry.MobType),
		fmt.Sprintf("# dropMode=%s", entry.Maf.DropMode),
		fmt.Sprintf("summon %s ~ ~ ~ %s", entry.MobType, enemySummonNBT(lootID, entry, itemsByID)),
	}
}

func enemySummonNBT(lootID string, entry enemyModel.Enemy, itemsByID map[string]itemModel.Item) string {
	nbt := deepCopyMap(entry.Minecraft)
	nbt["DeathLootTable"] = lootID
	nbt["Tags"] = mergeTags(nbt["Tags"], enemyTags(entry.ID, entry.Maf.EnemySkillIDs))

	if merged := mergeEquipment(nbt["equipment"], entry.Maf.Equipment, itemsByID); len(merged) > 0 {
		nbt["equipment"] = merged
	}
	if merged := mergeDropChances(nbt["drop_chances"], entry.Maf.Equipment); len(merged) > 0 {
		nbt["drop_chances"] = merged
	}
	if merged := mergePassengers(nbt["Passengers"], passengersNBT(entry.Passengers, itemsByID)); len(merged) > 0 {
		nbt["Passengers"] = merged
	}

	return MapToSNBT(nbt)
}

func passengersNBT(passengers []enemyModel.Passenger, itemsByID map[string]itemModel.Item) []any {
	if len(passengers) == 0 {
		return nil
	}
	parts := make([]any, 0, len(passengers))
	for _, p := range passengers {
		parts = append(parts, passengerNBT(p, itemsByID))
	}
	return parts
}

func passengerNBT(p enemyModel.Passenger, itemsByID map[string]itemModel.Item) map[string]any {
	nbt := deepCopyMap(p.Minecraft)
	nbt["id"] = p.MobType
	if tags := mergeTags(nbt["Tags"], passengerTags(p.Maf)); len(tags) > 0 {
		nbt["Tags"] = tags
	} else {
		delete(nbt, "Tags")
	}

	if p.Maf.Equipment != nil {
		if merged := mergeEquipment(nbt["equipment"], *p.Maf.Equipment, itemsByID); len(merged) > 0 {
			nbt["equipment"] = merged
		}
		if merged := mergeDropChances(nbt["drop_chances"], *p.Maf.Equipment); len(merged) > 0 {
			nbt["drop_chances"] = merged
		}
	}
	if merged := mergePassengers(nbt["Passengers"], passengersNBT(p.Passengers, itemsByID)); len(merged) > 0 {
		nbt["Passengers"] = merged
	}

	return nbt
}

func enemyTags(enemyID string, enemySkillIDs []string) []string {
	tags := []string{
		"maf_enemy",
		"maf_enemy_" + enemyID,
		"maf_vh_checked",
	}
	if len(enemySkillIDs) > 0 {
		tags = append(tags, "EnemySkill")
	}
	for _, skillID := range enemySkillIDs {
		id := strings.TrimSpace(skillID)
		if id == "" {
			continue
		}
		tags = append(tags, id, "maf_enemy_skill_"+id)
	}
	return tags
}

func passengerTags(maf enemyModel.PassengerMaf) []string {
	tags := make([]string, 0, len(maf.Tags)+len(maf.EnemySkillIDs)*2+1)
	for _, tag := range maf.Tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		tags = append(tags, tag)
	}
	if len(maf.EnemySkillIDs) > 0 {
		tags = append(tags, "EnemySkill")
	}
	for _, skillID := range maf.EnemySkillIDs {
		id := strings.TrimSpace(skillID)
		if id == "" {
			continue
		}
		tags = append(tags, id, "maf_enemy_skill_"+id)
	}
	return tags
}

func mergeTags(existing any, generated []string) []any {
	tags := make([]any, 0)
	seen := map[string]bool{}

	for _, raw := range anyList(existing) {
		tag, ok := raw.(string)
		if !ok {
			continue
		}
		tag = strings.TrimSpace(tag)
		if tag == "" || seen[tag] {
			continue
		}
		tags = append(tags, tag)
		seen[tag] = true
	}
	for _, raw := range generated {
		tag := strings.TrimSpace(raw)
		if tag == "" || seen[tag] {
			continue
		}
		tags = append(tags, tag)
		seen[tag] = true
	}
	return tags
}

func mergeEquipment(existing any, equipment model.Equipment, itemsByID map[string]itemModel.Item) map[string]any {
	merged := deepCopyMapFromAny(existing)
	slots := []struct {
		name string
		slot *model.EquipmentSlot
	}{
		{name: "mainhand", slot: equipment.Mainhand},
		{name: "offhand", slot: equipment.Offhand},
		{name: "head", slot: equipment.Head},
		{name: "chest", slot: equipment.Chest},
		{name: "legs", slot: equipment.Legs},
		{name: "feet", slot: equipment.Feet},
	}
	for _, entry := range slots {
		if entry.slot == nil {
			continue
		}
		merged[entry.name] = map[string]any{
			"id":    resolveEquipmentItemID(entry.slot, itemsByID),
			"count": entry.slot.Count,
		}
	}
	return merged
}

func mergeDropChances(existing any, equipment model.Equipment) map[string]any {
	merged := deepCopyMapFromAny(existing)
	slots := []struct {
		name string
		slot *model.EquipmentSlot
	}{
		{name: "mainhand", slot: equipment.Mainhand},
		{name: "offhand", slot: equipment.Offhand},
		{name: "head", slot: equipment.Head},
		{name: "chest", slot: equipment.Chest},
		{name: "legs", slot: equipment.Legs},
		{name: "feet", slot: equipment.Feet},
	}
	for _, entry := range slots {
		if entry.slot == nil {
			continue
		}
		dropChance := 0.085
		if entry.slot.DropChance != nil {
			dropChance = *entry.slot.DropChance
		}
		merged[entry.name] = dropChance
	}
	return merged
}

func mergePassengers(existing any, generated []any) []any {
	merged := anyList(existing)
	if len(generated) > 0 {
		merged = append(merged, generated...)
	}
	if len(merged) == 0 {
		return nil
	}
	return merged
}

func resolveEquipmentItemID(slot *model.EquipmentSlot, itemsByID map[string]itemModel.Item) string {
	if slot == nil {
		return ""
	}
	if slot.Kind == "item" {
		if entry, ok := itemsByID[slot.RefID]; ok && entry.ItemID != "" {
			return entry.ItemID
		}
	}
	return slot.RefID
}

func deepCopyMap(src map[string]any) map[string]any {
	if src == nil {
		return map[string]any{}
	}
	dst := make(map[string]any, len(src))
	for key, value := range src {
		dst[key] = deepCopyValue(value)
	}
	return dst
}

func deepCopyMapFromAny(value any) map[string]any {
	switch v := value.(type) {
	case nil:
		return map[string]any{}
	case map[string]any:
		return deepCopyMap(v)
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() || rv.Kind() != reflect.Map || rv.Type().Key().Kind() != reflect.String {
		return map[string]any{}
	}
	dst := make(map[string]any, rv.Len())
	iter := rv.MapRange()
	for iter.Next() {
		dst[iter.Key().String()] = deepCopyValue(iter.Value().Interface())
	}
	return dst
}

func anyList(value any) []any {
	switch v := value.(type) {
	case nil:
		return nil
	case []any:
		out := make([]any, 0, len(v))
		for _, item := range v {
			out = append(out, deepCopyValue(item))
		}
		return out
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() || (rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array) {
		return nil
	}
	out := make([]any, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		out = append(out, deepCopyValue(rv.Index(i).Interface()))
	}
	return out
}

func deepCopyValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		return deepCopyMap(v)
	case []any:
		out := make([]any, 0, len(v))
		for _, item := range v {
			out = append(out, deepCopyValue(item))
		}
		return out
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String {
		out := make(map[string]any, rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			out[iter.Key().String()] = deepCopyValue(iter.Value().Interface())
		}
		return out
	}
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		out := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			out = append(out, deepCopyValue(rv.Index(i).Interface()))
		}
		return out
	}
	return value
}
