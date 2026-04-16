package export_convert

import (
	"strings"
	"testing"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func TestResolveMafLootPoolsExpandsMafEntries(t *testing.T) {
	trueValue := true
	pools := []any{
		map[string]any{
			"rolls": 1.0,
			"entries": []any{
				map[string]any{
					"type":   "maf:item",
					"name":   "item_1",
					"weight": 7.0,
					"count":  2.0,
				},
				map[string]any{
					"type":   "maf:grimoire",
					"name":   "grimoire_1",
					"weight": 3.0,
					"count": map[string]any{
						"min": 1.0,
						"max": 2.0,
					},
				},
				map[string]any{
					"type":   "maf:passive",
					"name":   "passive_1",
					"slot":   1.0,
					"weight": 5.0,
				},
				map[string]any{
					"type":   "minecraft:item",
					"name":   "minecraft:apple",
					"weight": 1.0,
				},
			},
		},
	}

	resolved, err := ResolveMafLootPools(
		pools,
		map[string]itemModel.Item{
			"item_1": {
				ID: "item_1",
				Minecraft: itemModel.MinecraftItem{
					ItemID: "minecraft:stone",
				},
			},
		},
		map[string]grimoireModel.Grimoire{
			"grimoire_1": {
				ID:          "grimoire_1",
				CastTime:    20,
				CoolTime:    40,
				MPCost:      10,
				Script:      []string{"say test"},
				Title:       "Test",
				Description: "desc",
			},
		},
		map[string]passiveModel.Passive{
			"passive_1": {
				ID:               "passive_1",
				Name:             "Passive",
				Condition:        "always",
				Slots:            []int{1},
				Script:           []string{"say test"},
				GenerateGrimoire: &trueValue,
			},
		},
		nil,
		"loot_table(minecraft:chests/example)",
	)
	if err != nil {
		t.Fatalf("ResolveMafLootPools returned error: %v", err)
	}

	pool, ok := resolved[0].(map[string]any)
	if !ok {
		t.Fatalf("resolved pool must be object, got %T", resolved[0])
	}
	entries, ok := pool["entries"].([]any)
	if !ok {
		t.Fatalf("resolved entries must be array, got %T", pool["entries"])
	}
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	itemEntry := expectLootEntryObject(t, entries[0])
	if itemEntry["type"] != "minecraft:item" {
		t.Fatalf("expected minecraft:item, got %v", itemEntry["type"])
	}
	if itemEntry["name"] != "minecraft:stone" {
		t.Fatalf("expected minecraft:stone, got %v", itemEntry["name"])
	}
	if itemEntry["weight"] != 7.0 {
		t.Fatalf("expected weight 7.0, got %v", itemEntry["weight"])
	}
	if itemCount := firstSetCountValue(t, itemEntry); itemCount != 2.0 {
		t.Fatalf("expected count 2.0, got %v", itemCount)
	}

	grimoireEntry := expectLootEntryObject(t, entries[1])
	if grimoireEntry["type"] != "minecraft:item" || grimoireEntry["name"] != "minecraft:book" {
		t.Fatalf("unexpected grimoire entry: %#v", grimoireEntry)
	}
	grimoireCount := firstSetCountValue(t, grimoireEntry)
	countRange, ok := grimoireCount.(map[string]any)
	if !ok {
		t.Fatalf("expected range count object, got %T", grimoireCount)
	}
	if countRange["min"] != 1.0 || countRange["max"] != 2.0 {
		t.Fatalf("unexpected count range: %#v", countRange)
	}

	passiveEntry := expectLootEntryObject(t, entries[2])
	if passiveEntry["type"] != "minecraft:item" || passiveEntry["name"] != "minecraft:book" {
		t.Fatalf("unexpected passive entry: %#v", passiveEntry)
	}
	if passiveEntry["weight"] != 5.0 {
		t.Fatalf("expected weight 5.0, got %v", passiveEntry["weight"])
	}

	vanillaEntry := expectLootEntryObject(t, entries[3])
	if vanillaEntry["type"] != "minecraft:item" || vanillaEntry["name"] != "minecraft:apple" {
		t.Fatalf("unexpected vanilla entry: %#v", vanillaEntry)
	}
}

func TestResolveMafLootPoolsRejectsUnsupportedMafType(t *testing.T) {
	pools := []any{
		map[string]any{
			"entries": []any{
				map[string]any{
					"type": "maf:unknown",
					"name": "foo",
				},
			},
		},
	}

	_, err := ResolveMafLootPools(pools, nil, nil, nil, nil, "loot_table(minecraft:chests/example)")
	if err == nil {
		t.Fatal("expected error for unsupported maf type")
	}
	if !strings.Contains(err.Error(), "unsupported maf entry type") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveMafLootPoolsRejectsInvertedCountRange(t *testing.T) {
	pools := []any{
		map[string]any{
			"entries": []any{
				map[string]any{
					"type": "maf:item",
					"name": "item_1",
					"count": map[string]any{
						"min": 3.0,
						"max": 1.0,
					},
				},
			},
		},
	}

	_, err := ResolveMafLootPools(
		pools,
		map[string]itemModel.Item{
			"item_1": {
				ID: "item_1",
				Minecraft: itemModel.MinecraftItem{
					ItemID: "minecraft:stone",
				},
			},
		},
		nil,
		nil,
		nil,
		"loot_table(minecraft:chests/example)",
	)
	if err == nil {
		t.Fatal("expected error for inverted count range")
	}
	if !strings.Contains(err.Error(), "count.min must be less than or equal to count.max") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func expectLootEntryObject(t *testing.T, raw any) map[string]any {
	t.Helper()
	entry, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("entry must be object, got %T", raw)
	}
	return entry
}

func firstSetCountValue(t *testing.T, entry map[string]any) any {
	t.Helper()
	functions, ok := entry["functions"].([]any)
	if !ok || len(functions) == 0 {
		t.Fatalf("functions must be non-empty array")
	}
	firstFn, ok := functions[0].(map[string]any)
	if !ok {
		t.Fatalf("first function must be object, got %T", functions[0])
	}
	return firstFn["count"]
}
