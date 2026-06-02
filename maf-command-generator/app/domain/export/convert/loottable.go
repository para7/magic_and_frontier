package export_convert

import (
	"fmt"
	"strings"

	model "maf_command_editor/app/domain/model"
	activeModel "maf_command_editor/app/domain/model/active"
	bowModel "maf_command_editor/app/domain/model/bow"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func ResolveMafLootPools(
	pools []any,
	itemsByID map[string]itemModel.Item,
	activesByID map[string]activeModel.Active,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
	context string,
	ver ...int64,
) ([]any, error) {
	resolved := make([]any, 0, len(pools))
	for i, rawPool := range pools {
		pool, ok := rawPool.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%s: pools[%d] must be an object", context, i)
		}

		rawEntries, exists := pool["entries"]
		if !exists || rawEntries == nil {
			resolved = append(resolved, rawPool)
			continue
		}

		entries, ok := rawEntries.([]any)
		if !ok {
			return nil, fmt.Errorf("%s: pools[%d].entries must be an array", context, i)
		}

		resolvedEntries := make([]any, 0, len(entries))
		for j, rawEntry := range entries {
			entryContext := fmt.Sprintf("%s: pools[%d].entries[%d]", context, i, j)
			nextEntry, err := resolveMafLootEntry(rawEntry, itemsByID, activesByID, passivesByID, bowsByID, entryContext, itemVersion(ver...))
			if err != nil {
				return nil, err
			}
			resolvedEntries = append(resolvedEntries, nextEntry)
		}

		nextPool := make(map[string]any, len(pool))
		for key, value := range pool {
			nextPool[key] = value
		}
		nextPool["entries"] = resolvedEntries
		resolved = append(resolved, nextPool)
	}
	return resolved, nil
}

func resolveMafLootEntry(
	rawEntry any,
	itemsByID map[string]itemModel.Item,
	activesByID map[string]activeModel.Active,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
	context string,
	ver int64,
) (any, error) {
	entry, ok := rawEntry.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s must be an object", context)
	}

	rawType, _ := entry["type"].(string)
	entryType := strings.TrimSpace(rawType)
	if !strings.HasPrefix(entryType, "maf:") {
		return rawEntry, nil
	}

	rawName, _ := entry["name"].(string)
	refID := strings.TrimSpace(rawName)
	if refID == "" {
		return nil, fmt.Errorf("%s: maf entry requires non-empty name", context)
	}

	weight, hasWeight := entry["weight"]
	min, max, err := model.ParseLootEntryCount(entry)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", context, err)
	}

	switch entryType {
	case "maf:item":
		item, found := itemsByID[refID]
		if !found {
			return nil, fmt.Errorf("%s: referenced item not found (%s)", context, refID)
		}
		out, err := toItemLootEntry(item, activesByID, passivesByID, bowsByID, min, max, ver)
		if err != nil {
			return nil, err
		}
		if hasWeight {
			out["weight"] = weight
		}
		return out, nil
	case "maf:active":
		active, found := activesByID[refID]
		if !found {
			return nil, fmt.Errorf("%s: referenced active not found (%s)", context, refID)
		}
		out := toSpellLootEntry(active, min, max)
		if hasWeight {
			out["weight"] = weight
		}
		return out, nil
	case "maf:passive":
		passive, found := passivesByID[refID]
		if !found {
			return nil, fmt.Errorf("%s: referenced passive not found (%s)", context, refID)
		}
		slot, err := model.ParseLootEntrySlot(entry)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", context, err)
		}
		if !passiveHasSlot(passive, slot) {
			return nil, fmt.Errorf("%s: passive(%s) does not support slot %d", context, refID, slot)
		}
		out := toPassiveLootEntry(passive, slot, min, max)
		if hasWeight {
			out["weight"] = weight
		}
		return out, nil
	default:
		return nil, fmt.Errorf("%s: unsupported maf entry type (%s)", context, entryType)
	}
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

func toItemLootEntry(
	entry itemModel.Item,
	activesByID map[string]activeModel.Active,
	passivesByID map[string]passiveModel.Passive,
	bowsByID map[string]bowModel.BowPassive,
	min, max *float64,
	ver int64,
) (map[string]any, error) {
	customData, err := itemCustomData(entry, activesByID, passivesByID, bowsByID, ver)
	if err != nil {
		return nil, err
	}
	functions := []any{
		map[string]any{"function": "minecraft:set_count", "add": false, "count": ToCountValue(min, max)},
		map[string]any{"function": "minecraft:set_custom_data", "tag": customData},
	}
	components, err := itemComponentsForLoot(entry, activesByID, passivesByID, bowsByID)
	if err != nil {
		return nil, err
	}
	if len(components) > 0 {
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
	}, nil
}

func toSpellLootEntry(entry activeModel.Active, min, max *float64) map[string]any {
	book := activeSpellBookModel(entry)
	functions := []any{
		map[string]any{"function": "minecraft:set_count", "add": false, "count": ToCountValue(min, max)},
		map[string]any{
			"function":   "minecraft:set_components",
			"components": book.LootComponents(),
		},
		map[string]any{"function": "minecraft:set_custom_data", "tag": book.customData},
	}
	return map[string]any{
		"type":      "minecraft:item",
		"name":      "minecraft:book",
		"functions": functions,
	}
}

func toPassiveLootEntry(entry passiveModel.Passive, slot int, min, max *float64) map[string]any {
	book := passiveSpellBookModel(entry, slot)
	functions := []any{
		map[string]any{"function": "minecraft:set_count", "add": false, "count": ToCountValue(min, max)},
		map[string]any{
			"function":   "minecraft:set_components",
			"components": book.LootComponents(),
		},
		map[string]any{"function": "minecraft:set_custom_data", "tag": book.customData},
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
