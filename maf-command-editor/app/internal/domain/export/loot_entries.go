package export

import (
	"fmt"
	"math"
	"strings"

	"maf-command-editor/app/internal/domain/entity/grimoire"
	"maf-command-editor/app/internal/domain/entity/items"
)

func normalizeFunctionBody(script string) string {
	if strings.HasSuffix(script, "\n") {
		return script
	}
	return script + "\n"
}

func toSpellLootTable(entry grimoire.GrimoireEntry) map[string]any {
	return map[string]any{
		"type": "minecraft:generic",
		"pools": []any{
			map[string]any{
				"rolls": 1,
				"entries": []any{
					toSpellLootEntry(entry, ptrFloat(1), ptrFloat(1)),
				},
			},
		},
	}
}

func toItemLootTable(entry items.ItemEntry) map[string]any {
	value := float64(1)
	return map[string]any{
		"type": "minecraft:generic",
		"pools": []any{
			map[string]any{
				"rolls": 1,
				"entries": []any{
					toItemLootEntry(entry, &value, &value),
				},
			},
		},
	}
}

func toItemLootEntry(entry items.ItemEntry, min, max *float64) map[string]any {
	return map[string]any{
		"type": "minecraft:item",
		"name": entry.ItemID,
		"functions": []any{
			map[string]any{"function": "minecraft:set_count", "add": false, "count": toCountValue(min, max)},
			map[string]any{"function": "minecraft:set_custom_data", "tag": itemCustomData(entry)},
		},
	}
}

func toSpellLootEntry(entry grimoire.GrimoireEntry, min, max *float64) map[string]any {
	functions := []any{
		map[string]any{"function": "minecraft:set_count", "add": false, "count": toCountValue(min, max)},
		map[string]any{"function": "minecraft:set_name", "name": map[string]any{"text": entry.Title}, "target": "item_name"},
	}
	if lore := toLoreComponents(entry.Description); len(lore) > 0 {
		functions = append(functions, map[string]any{"function": "minecraft:set_lore", "mode": "append", "lore": lore})
	}
	functions = append(functions, map[string]any{"function": "minecraft:set_custom_data", "tag": spellCustomData(entry)})
	return map[string]any{
		"type":      "minecraft:item",
		"name":      "minecraft:written_book",
		"functions": functions,
	}
}

func itemCustomData(entry items.ItemEntry) string {
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

func spellCustomData(entry grimoire.GrimoireEntry) string {
	return fmt.Sprintf("{maf:{grimoire_id:%s,spell:{castid:%d,cost:%d,cast:%d,title:%s,description:%s}}}", jsonString(entry.ID), entry.CastID, entry.MPCost, entry.CastTime, jsonString(entry.Title), jsonString(entry.Description))
}

func grimoireDebugGiveCommand(entry grimoire.GrimoireEntry) string {
	parts := []string{
		fmt.Sprintf("item_name=%s", singleQuotedJSON(map[string]any{"text": entry.Title})),
	}
	if loreLines := linesToLoreValues(entry.Description); len(loreLines) > 0 {
		loreParts := make([]string, 0, len(loreLines))
		for _, line := range loreLines {
			loreParts = append(loreParts, singleQuotedJSON(map[string]any{"text": line}))
		}
		parts = append(parts, "lore=["+strings.Join(loreParts, ",")+"]")
	}
	parts = append(parts, "custom_data="+spellCustomData(entry))
	return fmt.Sprintf("give @s minecraft:written_book[%s] 1", strings.Join(parts, ","))
}

func toLoreComponents(value string) []any {
	lines := linesToLoreValues(value)
	out := make([]any, 0, len(lines))
	for _, line := range lines {
		out = append(out, map[string]any{"text": line})
	}
	return out
}

func linesToLoreValues(value string) []string {
	raw := strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
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
