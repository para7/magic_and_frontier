package export_convert

import (
	"fmt"
	"strings"

	model "maf_command_editor/app/domain/model"
	enemyModel "maf_command_editor/app/domain/model/enemy"
	itemModel "maf_command_editor/app/domain/model/item"
)

func ToEnemyFunctionLines(entry enemyModel.Enemy, lootID string, itemsByID map[string]itemModel.Item) []string {
	return []string{
		fmt.Sprintf("# enemyId=%s mobType=%s", entry.ID, entry.MobType),
		fmt.Sprintf("# dropMode=%s", entry.DropMode),
		fmt.Sprintf("summon %s ~ ~ ~ %s", entry.MobType, enemySummonNBT(lootID, entry, itemsByID)),
	}
}

func enemySummonNBT(lootID string, entry enemyModel.Enemy, itemsByID map[string]itemModel.Item) string {
	parts := []string{
		fmt.Sprintf("Health:%sf", formatFloat(entry.HP)),
		fmt.Sprintf("DeathLootTable:%s", JsonString(lootID)),
	}
	if entry.Name != "" {
		parts = append(parts, fmt.Sprintf("CustomName:{text:%s}", JsonString(entry.Name)))
	}
	if tags := enemyTags(entry); len(tags) > 0 {
		parts = append(parts, fmt.Sprintf("Tags:[%s]", strings.Join(tags, ",")))
	}
	if attrs := enemyAttributes(entry); len(attrs) > 0 {
		parts = append(parts, fmt.Sprintf("Attributes:[%s]", strings.Join(attrs, ",")))
	}
	if handItems, handDrops := equipmentArray(itemsByID, entry.Equipment.Mainhand, entry.Equipment.Offhand); handItems != "" {
		parts = append(parts, "HandItems:["+handItems+"]", "HandDropChances:["+handDrops+"]")
	}
	if armorItems, armorDrops := equipmentArray(itemsByID, entry.Equipment.Feet, entry.Equipment.Legs, entry.Equipment.Chest, entry.Equipment.Head); armorItems != "" {
		parts = append(parts, "ArmorItems:["+armorItems+"]", "ArmorDropChances:["+armorDrops+"]")
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func enemyTags(entry enemyModel.Enemy) []string {
	tags := []string{
		JsonString("maf_enemy"),
		JsonString("maf_enemy_" + entry.ID),
		JsonString("maf_vh_checked"),
	}
	if len(entry.EnemySkillIDs) > 0 {
		tags = append(tags, JsonString("EnemySkill"))
	}
	for _, skillID := range entry.EnemySkillIDs {
		tags = append(tags, JsonString(skillID), JsonString("maf_enemy_skill_"+skillID))
	}
	return tags
}

func enemyAttributes(entry enemyModel.Enemy) []string {
	attrs := []string{
		fmt.Sprintf("{Name:generic.max_health,Base:%s}", formatFloat(entry.HP)),
	}
	if entry.Attack != nil {
		attrs = append(attrs, fmt.Sprintf("{Name:generic.attack_damage,Base:%s}", formatFloat(*entry.Attack)))
	}
	if entry.Defense != nil {
		attrs = append(attrs, fmt.Sprintf("{Name:generic.armor,Base:%s}", formatFloat(*entry.Defense)))
	}
	if entry.MoveSpeed != nil {
		attrs = append(attrs, fmt.Sprintf("{Name:generic.movement_speed,Base:%s}", formatFloat(*entry.MoveSpeed)))
	}
	return attrs
}

func equipmentArray(itemsByID map[string]itemModel.Item, slots ...*model.EquipmentSlot) (string, string) {
	itemsOut := make([]string, 0, len(slots))
	dropsOut := make([]string, 0, len(slots))
	for _, slot := range slots {
		if slot == nil {
			itemsOut = append(itemsOut, "{}")
			dropsOut = append(dropsOut, "0.085F")
			continue
		}
		itemsOut = append(itemsOut, fmt.Sprintf("{id:%s,Count:%db}", JsonString(resolveEquipmentItemID(slot, itemsByID)), slot.Count))
		dropChance := 0.085
		if slot.DropChance != nil {
			dropChance = *slot.DropChance
		}
		dropsOut = append(dropsOut, formatFloat(dropChance)+"F")
	}
	return strings.Join(itemsOut, ","), strings.Join(dropsOut, ",")
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
