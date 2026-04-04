package export_convert

import (
	"fmt"

	model "maf_command_editor/app/domain/model"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func BuildDropLootPool(
	drops []model.DropRef,
	itemsByID map[string]itemModel.Item,
	grimoiresByID map[string]grimoireModel.Grimoire,
	passivesByID map[string]passiveModel.Passive,
	context string,
) (map[string]any, error) {
	entries := make([]any, 0, len(drops))
	for _, drop := range drops {
		switch drop.Kind {
		case "minecraft_item":
			entries = append(entries, map[string]any{
				"type":   "minecraft:item",
				"name":   drop.RefID,
				"weight": ToWeight(drop.Weight),
				"functions": []any{
					map[string]any{"function": "minecraft:set_count", "count": ToCountValue(drop.CountMin, drop.CountMax)},
				},
			})
		case "item":
			item, ok := itemsByID[drop.RefID]
			if !ok {
				return nil, fmt.Errorf("%s: referenced item not found (%s)", context, drop.RefID)
			}
			entry := toItemLootEntry(item, drop.CountMin, drop.CountMax)
			entry["weight"] = ToWeight(drop.Weight)
			entries = append(entries, entry)
		case "grimoire":
			entry, ok := grimoiresByID[drop.RefID]
			if !ok {
				return nil, fmt.Errorf("%s: referenced grimoire not found (%s)", context, drop.RefID)
			}
			out := toSpellLootEntry(entry, drop.CountMin, drop.CountMax)
			out["weight"] = ToWeight(drop.Weight)
			entries = append(entries, out)
		case "passive":
			if drop.Slot == nil {
				return nil, fmt.Errorf("%s: passive drop requires slot (%s)", context, drop.RefID)
			}
			slot := *drop.Slot
			entry, ok := passivesByID[drop.RefID]
			if !ok {
				return nil, fmt.Errorf("%s: referenced passive not found (%s)", context, drop.RefID)
			}
			if !passiveHasSlot(entry, slot) {
				return nil, fmt.Errorf("%s: passive(%s) does not support slot %d", context, drop.RefID, slot)
			}
			out := toPassiveLootEntry(entry, slot, drop.CountMin, drop.CountMax)
			out["weight"] = ToWeight(drop.Weight)
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

func MergeLootTablePools(base map[string]any, pool map[string]any, tablePath string) (map[string]any, error) {
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

func toItemLootEntry(entry itemModel.Item, min, max *float64) map[string]any {
	functions := []any{
		map[string]any{"function": "minecraft:set_count", "add": false, "count": ToCountValue(min, max)},
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

func toSpellLootEntry(entry grimoireModel.Grimoire, min, max *float64) map[string]any {
	functions := []any{
		map[string]any{"function": "minecraft:set_count", "add": false, "count": ToCountValue(min, max)},
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

func toPassiveLootEntry(entry passiveModel.Passive, slot int, min, max *float64) map[string]any {
	functions := []any{
		map[string]any{"function": "minecraft:set_count", "add": false, "count": ToCountValue(min, max)},
		map[string]any{
			"function":   "minecraft:set_components",
			"components": passiveLootComponents(entry, slot),
		},
		map[string]any{"function": "minecraft:set_custom_data", "tag": passiveSpellCustomData(entry, slot)},
	}
	return map[string]any{
		"type":      "minecraft:item",
		"name":      "minecraft:book",
		"functions": functions,
	}
}

func passiveHasSlot(entry passiveModel.Passive, slot int) bool {
	for _, v := range entry.Slots {
		if v == slot {
			return true
		}
	}
	return false
}
