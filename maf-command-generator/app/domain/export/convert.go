package export

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	model "maf_command_editor/app/domain/model"
	enemyModel "maf_command_editor/app/domain/model/enemy"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
)

// 独自形式のグリモアのデータをマインクラフトの本に変換する
func grimoireToBook(grimoire grimoireModel.Grimoire) string {
	//example:
	// minecraft:book[minecraft:item_name={text:"動作確認: スウィープ10"},minecraft:lore=[{text:"右クリックで詠唱を開始"},{text:"effect=13 cast=10 cost=0"}],minecraft:consumable={consume_seconds:99999,animation:"bow",has_consume_particles:false},minecraft:custom_data={maf:{grimoire_id:"dev_book2",spell:{castid:13,cost:10,cast:10,title:"スウィープ",description:"近くのアイテムをかき集める"}}}]

	itemName := fmt.Sprintf(`{text:"%s%d"}`, grimoire.Title, grimoire.CastTime)
	lore := fmt.Sprintf(`[{text:"右クリックで詠唱を開始"},{text:"effect=%d cast=%d cost=%d"}]`,
		grimoire.CastID, grimoire.CastTime, grimoire.MPCost)
	consumable := `{consume_seconds:99999,animation:"bow",has_consume_particles:false}`
	customData := fmt.Sprintf(`{maf:{grimoire_id:"%s",spell:{castid:%d,cost:%d,cast:%d,title:"%s",description:"%s"}}}`,
		grimoire.ID, grimoire.CastID, grimoire.MPCost, grimoire.CastTime, grimoire.Title, grimoire.Description)

	return fmt.Sprintf(
		`minecraft:book[minecraft:item_name=%s,minecraft:lore=%s,minecraft:consumable=%s,minecraft:custom_data=%s]`,
		itemName, lore, consumable, customData,
	)
}

func toItemLootEntry(entry itemModel.Item, min, max *float64) map[string]any {
	functions := []any{
		map[string]any{"function": "minecraft:set_count", "add": false, "count": toCountValue(min, max)},
		map[string]any{"function": "minecraft:set_custom_data", "tag": itemCustomData(entry)},
	}
	if components := itemComponentsForLoot(entry); len(components) > 0 {
		functions = append(functions, map[string]any{
			"function":   "minecraft:set_components",
			"components": components,
		})
	}
	if enchMap := itemEnchantmentsForLoot(entry); len(enchMap) > 0 {
		functions = append(functions, map[string]any{
			"function":     "minecraft:set_enchantments",
			"enchantments": enchMap,
			"add":          false,
		})
	}
	return map[string]any{
		"type":      "minecraft:item",
		"name":      entry.ItemID,
		"functions": functions,
	}
}

func itemComponentsForLoot(entry itemModel.Item) map[string]any {
	components := map[string]any{}

	if name := strings.TrimSpace(entry.CustomName); name != "" {
		components["minecraft:custom_name"] = map[string]any{"text": name}
	}

	var loreLines []string
	for _, line := range strings.Split(entry.Lore, "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line != "" {
			loreLines = append(loreLines, line)
		}
	}
	if len(loreLines) > 0 {
		lore := make([]any, 0, len(loreLines))
		for _, line := range loreLines {
			lore = append(lore, map[string]any{"text": line})
		}
		components["minecraft:lore"] = lore
	}

	if entry.Unbreakable {
		components["minecraft:unbreakable"] = map[string]any{}
	}

	return components
}

func itemEnchantmentsForLoot(entry itemModel.Item) map[string]any {
	ench := strings.TrimSpace(entry.Enchantments)
	if ench == "" {
		return nil
	}
	enchMap := map[string]any{}
	for _, line := range strings.Split(strings.ReplaceAll(ench, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			if level, err := strconv.Atoi(fields[1]); err == nil {
				enchMap[fields[0]] = level
			}
		}
	}
	return enchMap
}

func toSpellLootEntry(entry grimoireModel.Grimoire, min, max *float64) map[string]any {
	functions := []any{
		map[string]any{"function": "minecraft:set_count", "add": false, "count": toCountValue(min, max)},
		map[string]any{
			"function": "minecraft:set_components",
			"components": map[string]any{
				"minecraft:item_name": map[string]any{"text": fmt.Sprintf("%s%d", entry.Title, entry.CastTime)},
				"minecraft:lore": []any{
					map[string]any{"text": "右クリックで詠唱を開始"},
					map[string]any{"text": fmt.Sprintf("effect=%d cast=%d cost=%d", entry.CastID, entry.CastTime, entry.MPCost)},
				},
				"minecraft:consumable": map[string]any{
					"consume_seconds":       99999.0,
					"animation":             "bow",
					"has_consume_particles": false,
				},
			},
		},
		map[string]any{"function": "minecraft:set_custom_data", "tag": spellCustomData(entry)},
	}
	return map[string]any{
		"type":      "minecraft:item",
		"name":      "minecraft:book",
		"functions": functions,
	}
}

func buildDropLootPool(drops []model.DropRef, itemsByID map[string]itemModel.Item, grimoiresByID map[string]grimoireModel.Grimoire, context string) (map[string]any, error) {
	entries := make([]any, 0, len(drops))
	for _, drop := range drops {
		switch drop.Kind {
		case "minecraft_item":
			entries = append(entries, map[string]any{
				"type":   "minecraft:item",
				"name":   drop.RefID,
				"weight": toWeight(drop.Weight),
				"functions": []any{
					map[string]any{"function": "minecraft:set_count", "count": toCountValue(drop.CountMin, drop.CountMax)},
				},
			})
		case "item":
			item, ok := itemsByID[drop.RefID]
			if !ok {
				return nil, fmt.Errorf("%s: referenced item not found (%s)", context, drop.RefID)
			}
			entry := toItemLootEntry(item, drop.CountMin, drop.CountMax)
			entry["weight"] = toWeight(drop.Weight)
			entries = append(entries, entry)
		case "grimoire":
			entry, ok := grimoiresByID[drop.RefID]
			if !ok {
				return nil, fmt.Errorf("%s: referenced grimoire not found (%s)", context, drop.RefID)
			}
			out := toSpellLootEntry(entry, drop.CountMin, drop.CountMax)
			out["weight"] = toWeight(drop.Weight)
			entries = append(entries, out)
		default:
			return nil, fmt.Errorf("%s: unsupported drop kind (%s)", context, drop.Kind)
		}
	}
	return map[string]any{
		"rolls":   1,
		"entries": entries,
	}, nil
}

func mergeLootTablePools(base map[string]any, pool map[string]any, tablePath string) (map[string]any, error) {
	if base == nil {
		base = map[string]any{}
	}
	if rawPools, ok := base["pools"]; ok && rawPools != nil {
		pools, ok := rawPools.([]any)
		if !ok {
			return nil, fmt.Errorf("enemy(%s): base loot table pools must be an array", tablePath)
		}
		base["pools"] = append(pools, pool)
		return base, nil
	}
	base["pools"] = []any{pool}
	return base, nil
}

func toEnemyFunctionLines(entry enemyModel.Enemy, lootID string, itemsByID map[string]itemModel.Item) []string {
	return []string{
		fmt.Sprintf("# enemyId=%s mobType=%s", entry.ID, entry.MobType),
		fmt.Sprintf("# dropMode=%s", entry.DropMode),
		fmt.Sprintf("summon %s ~ ~ ~ %s", entry.MobType, enemySummonNBT(lootID, entry, itemsByID)),
	}
}

func enemySummonNBT(lootID string, entry enemyModel.Enemy, itemsByID map[string]itemModel.Item) string {
	parts := []string{
		fmt.Sprintf("Health:%sf", formatFloat(entry.HP)),
		fmt.Sprintf("DeathLootTable:%s", jsonString(lootID)),
	}
	if entry.Name != "" {
		parts = append(parts, fmt.Sprintf("CustomName:{text:%s}", jsonString(entry.Name)))
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
		jsonString("maf_enemy"),
		jsonString("maf_enemy_" + entry.ID),
		jsonString("maf_vh_checked"),
	}
	if len(entry.EnemySkillIDs) > 0 {
		tags = append(tags, jsonString("EnemySkill"))
	}
	for _, skillID := range entry.EnemySkillIDs {
		tags = append(tags, jsonString(skillID), jsonString("maf_enemy_skill_"+skillID))
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
		itemsOut = append(itemsOut, fmt.Sprintf("{id:%s,Count:%db}", jsonString(resolveEquipmentItemID(slot, itemsByID)), slot.Count))
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

func itemCustomData(entry itemModel.Item) string {
	parts := []string{
		fmt.Sprintf("item_id:%s", jsonString(entry.ItemID)),
		fmt.Sprintf("source_id:%s", jsonString(entry.ID)),
		fmt.Sprintf("nbt_snapshot:%s", jsonString(entry.NBT)),
	}
	if entry.SkillID != "" {
		parts = append(parts, "maf_skill:1b", fmt.Sprintf("maf_skill_id:%s", jsonString(entry.SkillID)))
	}
	return "{maf:{" + strings.Join(parts, ",") + "}}"
}

func spellCustomData(entry grimoireModel.Grimoire) string {
	return fmt.Sprintf(
		"{maf:{grimoire_id:%s,spell:{castid:%d,cost:%d,cast:%d,title:%s,description:%s}}}",
		jsonString(entry.ID),
		entry.CastID,
		entry.MPCost,
		entry.CastTime,
		jsonString(entry.Title),
		jsonString(entry.Description),
	)
}

func toCountValue(min, max *float64) any {
	minValue := 1.0
	maxValue := 1.0
	if min != nil {
		minValue = *min
	}
	if max != nil {
		maxValue = *max
	}
	if minValue == maxValue {
		return minValue
	}
	return map[string]any{
		"type": "minecraft:uniform",
		"min":  minValue,
		"max":  maxValue,
	}
}

func toWeight(weight float64) int {
	if !isFinite(weight) || weight <= 0 {
		return 1
	}
	return int(math.Floor(weight))
}

func isFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func jsonString(value string) string {
	return string(mustJSON(value))
}

func mustJSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return data
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}
