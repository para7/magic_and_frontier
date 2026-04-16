package export_convert

import (
	"strings"
	"testing"

	model "maf_command_editor/app/domain/model"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func TestBuildDropLootPoolAllowsPassiveWithGenerateGrimoireFalse(t *testing.T) {
	falseValue := false
	slot := 1
	_, err := BuildDropLootPool(
		[]model.DropRef{{Kind: "passive", RefID: "passive_1", Slot: &slot, Weight: 1}},
		nil,
		nil,
		map[string]passiveModel.Passive{
			"passive_1": {
				ID:               "passive_1",
				Name:             "Passive",
				Condition:        "always",
				Slots:            []int{1},
				Script:           []string{"say test"},
				GenerateGrimoire: &falseValue,
			},
		},
		nil,
		"enemy(enemy_1)",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestBuildDropLootPoolPassiveStillRequiresSupportedSlot(t *testing.T) {
	trueValue := true
	slot := 2
	_, err := BuildDropLootPool(
		[]model.DropRef{{Kind: "passive", RefID: "passive_1", Slot: &slot, Weight: 1}},
		nil,
		nil,
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
		"enemy(enemy_1)",
	)
	if err == nil {
		t.Fatal("expected slot validation error")
	}
}

func TestResolveMafLootPoolsExpandsMafEntries(t *testing.T) {
	pools := []any{
		map[string]any{
			"rolls": 1.0,
			"entries": []any{
				map[string]any{
					"type":   "maf:item",
					"name":   "item_1",
					"weight": 7.0,
				},
				map[string]any{
					"type":   "maf:grimoire",
					"name":   "grimoire_1",
					"weight": 3.0,
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
				ID:       "grimoire_1",
				CastTime: 20,
				CoolTime: 40,
				MPCost:   10,
				Script:   []string{"say test"},
				Title:    "Test",
			},
		},
		nil,
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
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	itemEntry, ok := entries[0].(map[string]any)
	if !ok {
		t.Fatalf("item entry must be object, got %T", entries[0])
	}
	if itemEntry["type"] != "minecraft:item" {
		t.Fatalf("expected minecraft:item, got %v", itemEntry["type"])
	}
	if itemEntry["name"] != "minecraft:stone" {
		t.Fatalf("expected minecraft:stone, got %v", itemEntry["name"])
	}
	if itemEntry["weight"] != 7.0 {
		t.Fatalf("expected weight 7.0, got %v", itemEntry["weight"])
	}
	itemFunctions, ok := itemEntry["functions"].([]any)
	if !ok || len(itemFunctions) == 0 {
		t.Fatalf("item functions must be non-empty array")
	}
	firstFunction, ok := itemFunctions[0].(map[string]any)
	if !ok {
		t.Fatalf("item first function must be object, got %T", itemFunctions[0])
	}
	if firstFunction["count"] != 1.0 {
		t.Fatalf("expected default count 1.0, got %v", firstFunction["count"])
	}

	grimoireEntry, ok := entries[1].(map[string]any)
	if !ok {
		t.Fatalf("grimoire entry must be object, got %T", entries[1])
	}
	if grimoireEntry["type"] != "minecraft:item" {
		t.Fatalf("expected minecraft:item for grimoire, got %v", grimoireEntry["type"])
	}
	if grimoireEntry["name"] != "minecraft:book" {
		t.Fatalf("expected minecraft:book, got %v", grimoireEntry["name"])
	}
	if grimoireEntry["weight"] != 3.0 {
		t.Fatalf("expected weight 3.0, got %v", grimoireEntry["weight"])
	}
	grimoireFunctions, ok := grimoireEntry["functions"].([]any)
	if !ok || len(grimoireFunctions) == 0 {
		t.Fatalf("grimoire functions must be non-empty array")
	}
	grimoireFirstFunction, ok := grimoireFunctions[0].(map[string]any)
	if !ok {
		t.Fatalf("grimoire first function must be object, got %T", grimoireFunctions[0])
	}
	if grimoireFirstFunction["count"] != 1.0 {
		t.Fatalf("expected default count 1.0, got %v", grimoireFirstFunction["count"])
	}

	vanillaEntry, ok := entries[2].(map[string]any)
	if !ok {
		t.Fatalf("vanilla entry must be object, got %T", entries[2])
	}
	if vanillaEntry["type"] != "minecraft:item" || vanillaEntry["name"] != "minecraft:apple" {
		t.Fatalf("unexpected vanilla entry: %#v", vanillaEntry)
	}
}

func TestResolveMafLootPoolsRejectsUnsupportedMafType(t *testing.T) {
	pools := []any{
		map[string]any{
			"entries": []any{
				map[string]any{
					"type": "maf:passive",
					"name": "passive_1",
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
