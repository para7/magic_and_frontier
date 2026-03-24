package export

import (
	"fmt"

	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/treasures"
)

func buildDropLootTable(drops []treasures.DropRef, itemsByID map[string]items.ItemEntry, grimoiresByID map[string]grimoire.GrimoireEntry, context string) (map[string]any, error) {
	pool, err := buildDropLootPool(drops, itemsByID, grimoiresByID, context)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"type":  "minecraft:generic",
		"pools": []any{pool},
	}, nil
}

func buildDropLootPool(drops []treasures.DropRef, itemsByID map[string]items.ItemEntry, grimoiresByID map[string]grimoire.GrimoireEntry, context string) (map[string]any, error) {
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
			return nil, fmt.Errorf("treasure(%s): base loot table pools must be an array", tablePath)
		}
		base["pools"] = append(pools, pool)
		return base, nil
	}
	base["pools"] = []any{pool}
	return base, nil
}
