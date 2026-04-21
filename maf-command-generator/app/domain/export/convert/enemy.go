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
	if entry.IsBaby != nil && *entry.IsBaby {
		parts = append(parts, "IsBaby:1b")
	}
	if tags := enemyTags(entry); len(tags) > 0 {
		parts = append(parts, fmt.Sprintf("Tags:[%s]", strings.Join(tags, ",")))
	}
	if attrs := enemyattributes(entry); len(attrs) > 0 {
		parts = append(parts, fmt.Sprintf("attributes:[%s]", strings.Join(attrs, ",")))
	}
	if handItems, handDrops := equipmentArray(itemsByID, entry.Equipment.Mainhand, entry.Equipment.Offhand); handItems != "" {
		parts = append(parts, "HandItems:["+handItems+"]", "HandDropChances:["+handDrops+"]")
	}
	if armorItems, armorDrops := equipmentArray(itemsByID, entry.Equipment.Feet, entry.Equipment.Legs, entry.Equipment.Chest, entry.Equipment.Head); armorItems != "" {
		parts = append(parts, "ArmorItems:["+armorItems+"]", "ArmorDropChances:["+armorDrops+"]")
	}
	if pnbt := passengersNBT(entry.Passengers); pnbt != "" {
		parts = append(parts, "Passengers:["+pnbt+"]")
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func passengersNBT(passengers []enemyModel.PassengerEntity) string {
	if len(passengers) == 0 {
		return ""
	}
	parts := make([]string, 0, len(passengers))
	for _, p := range passengers {
		parts = append(parts, passengerNBT(p))
	}
	return strings.Join(parts, ",")
}

func passengerNBT(p enemyModel.PassengerEntity) string {
	parts := []string{fmt.Sprintf("id:%s", JsonString(p.MobType))}
	if p.Name != "" {
		parts = append(parts, fmt.Sprintf("CustomName:%s", JsonString(p.Name)))
	}

	attrs := []string{}
	if p.HP != nil {
		attrs = append(attrs, fmt.Sprintf("{id:\"minecraft:max_health\",base:%s}", formatFloat(*p.HP)))
		parts = append(parts, fmt.Sprintf("Health:%sf", formatFloat(*p.HP)))
	}
	if p.Attack != nil {
		attrs = append(attrs, fmt.Sprintf("{id:\"minecraft:attack_damage\",base:%s}", formatFloat(*p.Attack)))
	}
	if p.Defense != nil {
		attrs = append(attrs, fmt.Sprintf("{id:\"minecraft:armor\",base:%s}", formatFloat(*p.Defense)))
	}
	if p.MoveSpeed != nil {
		attrs = append(attrs, fmt.Sprintf("{id:\"minecraft:movement_speed\",base:%s}", formatFloat(*p.MoveSpeed)))
	}
	if len(attrs) > 0 {
		parts = append(parts, fmt.Sprintf("attributes:[%s]", strings.Join(attrs, ",")))
	}

	tags := make([]string, 0, len(p.Tags)+len(p.EnemySkillIDs)*2+1)
	for _, t := range p.Tags {
		tags = append(tags, JsonString(t))
	}
	if len(p.EnemySkillIDs) > 0 {
		tags = append(tags, JsonString("EnemySkill"))
	}
	for _, skillID := range p.EnemySkillIDs {
		tags = append(tags, JsonString(skillID), JsonString("maf_enemy_skill_"+skillID))
	}
	if len(tags) > 0 {
		parts = append(parts, fmt.Sprintf("Tags:[%s]", strings.Join(tags, ",")))
	}
	if p.IsBaby != nil && *p.IsBaby {
		parts = append(parts, "IsBaby:1b")
	}
	if nested := passengersNBT(p.Passengers); nested != "" {
		parts = append(parts, "Passengers:["+nested+"]")
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

func enemyattributes(entry enemyModel.Enemy) []string {
	attrs := []string{
		fmt.Sprintf("{id:\"minecraft:max_health\",base:%s}", formatFloat(entry.HP)),
	}
	if entry.Attack != nil {
		attrs = append(attrs, fmt.Sprintf("{id:\"minecraft:attack_damage\",base:%s}", formatFloat(*entry.Attack)))
	}
	if entry.Defense != nil {
		attrs = append(attrs, fmt.Sprintf("{id:\"minecraft:armor\",base:%s}", formatFloat(*entry.Defense)))
	}
	if entry.MoveSpeed != nil {
		attrs = append(attrs, fmt.Sprintf("{id:\"minecraft:movement_speed\",base:%s}", formatFloat(*entry.MoveSpeed)))
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
		if entry, ok := itemsByID[slot.RefID]; ok && entry.Minecraft.ItemID != "" {
			return entry.Minecraft.ItemID
		}
	}
	return slot.RefID
}
