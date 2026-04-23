package export_convert

import (
	"strings"
	"testing"

	bowModel "maf_command_editor/app/domain/model/bow"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func TestItemLootHelpersReadMinecraftComponents(t *testing.T) {
	entry := itemModel.Item{
		ID:     "items_1",
		ItemID: "minecraft:stone",
		Maf: itemModel.ItemMaf{
			GrimoireID: "tempest01",
			PassiveID:  "regeneration",
		},
		Minecraft: map[string]any{
			"components": map[string]any{
				"minecraft:custom_name": map[string]any{"text": "Starter Stone"},
				"minecraft:lore": []any{
					map[string]any{"text": "Sample item"},
				},
				"minecraft:unbreakable":  map[string]any{},
				"minecraft:enchantments": map[string]any{"minecraft:sharpness": 5},
				"minecraft:foo":          map[string]any{"bar": true},
			},
		},
	}

	grimoiresByID := map[string]grimoireModel.Grimoire{
		"tempest01": {
			ID:          "tempest01",
			MPCost:      13,
			CastTime:    40,
			CoolTime:    20,
			Title:       "テンペスト",
			Description: "敵1体に雷を落とし周辺に特大ダメージ",
		},
	}
	passivesByID := map[string]passiveModel.Passive{
		"regeneration": {
			ID:          "regeneration",
			Name:        "いつでもリジェネ",
			Condition:   "always",
			Slots:       []int{1},
			Description: "",
		},
	}

	customData, err := itemCustomData(entry, grimoiresByID, passivesByID, nil)
	if err != nil {
		t.Fatalf("itemCustomData returned error: %v", err)
	}
	if !strings.Contains(customData, `item_id:"minecraft:stone"`) {
		t.Fatalf("item_id missing from custom data: %s", customData)
	}
	if !strings.Contains(customData, `grimoire_id:"tempest01"`) {
		t.Fatalf("grimoire_id missing from custom data: %s", customData)
	}
	if !strings.Contains(customData, `hasPassive:1b`) {
		t.Fatalf("hasPassive tag missing from custom data: %s", customData)
	}
	if !strings.Contains(customData, `passiveId:"regeneration"`) {
		t.Fatalf("passiveId missing from custom data: %s", customData)
	}
	if !strings.Contains(customData, `passiveSlot:1`) {
		t.Fatalf("passiveSlot missing from custom data: %s", customData)
	}
	if !strings.Contains(customData, `spell:{kind:"grimoire",id:"tempest01",cost:13,cast:40,cooltime:20`) {
		t.Fatalf("spell metadata should be derived from grimoire: %s", customData)
	}
	if !strings.Contains(customData, `nbt_snapshot:{`) {
		t.Fatalf("nbt snapshot should be derived from components: %s", customData)
	}

	components, err := itemComponentsForLoot(entry, grimoiresByID, passivesByID, nil)
	if err != nil {
		t.Fatalf("itemComponentsForLoot returned error: %v", err)
	}
	if _, ok := components["minecraft:custom_name"]; !ok {
		t.Fatalf("custom_name should be exported: %#v", components)
	}
	if _, ok := components["minecraft:lore"]; !ok {
		t.Fatalf("lore should be exported: %#v", components)
	}
	if _, ok := components["minecraft:unbreakable"]; !ok {
		t.Fatalf("unbreakable should be exported: %#v", components)
	}
	if _, ok := components["minecraft:foo"]; !ok {
		t.Fatalf("minecraft custom component should pass through: %#v", components)
	}
	if _, ok := components["minecraft:enchantments"]; ok {
		t.Fatalf("enchantments should be handled by set_enchantments: %#v", components)
	}
	if _, ok := components["minecraft:consumable"]; !ok {
		t.Fatalf("consumable should be added for spell items: %#v", components)
	}

	enchantments := itemEnchantmentsForLoot(entry)
	if enchantments["minecraft:sharpness"] != 5 {
		t.Fatalf("unexpected enchantments: %#v", enchantments)
	}
}

func TestPassiveOnlyItemDoesNotBecomeRightClickSpell(t *testing.T) {
	entry := itemModel.Item{
		ID:     "items_passive_only",
		ItemID: "minecraft:stone",
		Maf: itemModel.ItemMaf{
			PassiveID: "regeneration",
		},
		Minecraft: map[string]any{
			"components": map[string]any{
				"minecraft:custom_name": map[string]any{"text": "Passive Only"},
			},
		},
	}
	passivesByID := map[string]passiveModel.Passive{
		"regeneration": {
			ID:        "regeneration",
			Condition: "always",
			Slots:     []int{1},
		},
	}

	customData, err := itemCustomData(entry, nil, passivesByID, nil)
	if err != nil {
		t.Fatalf("itemCustomData returned error: %v", err)
	}
	if !strings.Contains(customData, `hasPassive:1b`) || !strings.Contains(customData, `passiveId:"regeneration"`) {
		t.Fatalf("passive metadata should be embedded: %s", customData)
	}
	if strings.Contains(customData, `spell:{`) {
		t.Fatalf("passive-only item should not embed spell metadata: %s", customData)
	}

	components, err := itemComponentsForLoot(entry, nil, passivesByID, nil)
	if err != nil {
		t.Fatalf("itemComponentsForLoot returned error: %v", err)
	}
	if _, ok := components["minecraft:consumable"]; ok {
		t.Fatalf("passive-only item should not be consumable: %#v", components)
	}
}

func TestItemCustomDataEmbedsMaxMPWhenConfigured(t *testing.T) {
	maxMP := -12
	entry := itemModel.Item{
		ID:     "items_with_maxmp",
		ItemID: "minecraft:stone",
		Maf: itemModel.ItemMaf{
			MaxMP: &maxMP,
		},
	}

	customData, err := itemCustomData(entry, nil, nil, nil)
	if err != nil {
		t.Fatalf("itemCustomData returned error: %v", err)
	}
	if !strings.Contains(customData, `maxmp:-12`) {
		t.Fatalf("maxmp should be embedded into custom data: %s", customData)
	}
}

func TestItemToGiveCommandBuildsSortedComponentsAndCustomData(t *testing.T) {
	entry := itemModel.Item{
		ID:     "items_1",
		ItemID: "minecraft:stone",
		Maf: itemModel.ItemMaf{
			GrimoireID: "tempest01",
		},
		Minecraft: map[string]any{
			"components": map[string]any{
				"minecraft:lore": []any{
					map[string]any{"text": "Sample item"},
				},
				"minecraft:custom_name": map[string]any{"text": "Starter Stone"},
			},
		},
	}
	grimoiresByID := map[string]grimoireModel.Grimoire{
		"tempest01": {
			ID:          "tempest01",
			MPCost:      13,
			CastTime:    40,
			CoolTime:    20,
			Title:       "テンペスト",
			Description: "敵1体に雷を落とし周辺に特大ダメージ",
		},
	}

	command, err := ItemToGiveCommand(entry, grimoiresByID, nil, nil)
	if err != nil {
		t.Fatalf("ItemToGiveCommand returned error: %v", err)
	}
	if !strings.Contains(command, `give @p minecraft:stone[`) {
		t.Fatalf("unexpected give command: %s", command)
	}
	if !strings.Contains(command, `minecraft:consumable={consume_seconds:99999,animation:"bow",has_consume_particles:false}`) {
		t.Fatalf("spell item should include consumable: %s", command)
	}
	if !strings.Contains(command, `minecraft:custom_data={maf:{`) {
		t.Fatalf("custom_data missing from give command: %s", command)
	}
	customNameIndex := strings.Index(command, "minecraft:custom_name=")
	loreIndex := strings.Index(command, "minecraft:lore=")
	if customNameIndex == -1 || loreIndex == -1 || customNameIndex > loreIndex {
		t.Fatalf("components should be sorted by key: %s", command)
	}
	if !strings.Contains(command, `minecraft:custom_name={text:"Starter Stone"}`) {
		t.Fatalf("custom_name should be rendered from JSON object: %s", command)
	}
	if !strings.Contains(command, `minecraft:lore=[{text:"Sample item"}]`) {
		t.Fatalf("lore should be rendered from JSON array: %s", command)
	}
}

func TestItemToGiveCommandOverridesMinecraftConsumableForSpellItems(t *testing.T) {
	entry := itemModel.Item{
		ID:     "items_1",
		ItemID: "minecraft:stone",
		Maf: itemModel.ItemMaf{
			GrimoireID: "tempest01",
		},
		Minecraft: map[string]any{
			"components": map[string]any{
				"minecraft:consumable": map[string]any{"consume_seconds": 10},
			},
		},
	}
	grimoiresByID := map[string]grimoireModel.Grimoire{
		"tempest01": {ID: "tempest01", MPCost: 1, CastTime: 1, CoolTime: 1, Title: "Spell"},
	}

	command, err := ItemToGiveCommand(entry, grimoiresByID, nil, nil)
	if err != nil {
		t.Fatalf("ItemToGiveCommand returned error: %v", err)
	}
	if strings.Count(command, "minecraft:consumable=") != 1 {
		t.Fatalf("consumable should appear exactly once: %s", command)
	}
	if !strings.Contains(command, `minecraft:consumable={consume_seconds:99999,animation:"bow",has_consume_particles:false}`) {
		t.Fatalf("maf consumable should override minecraft value: %s", command)
	}
}

func TestItemComponentsForLootOverridesMinecraftConsumableForSpellItems(t *testing.T) {
	entry := itemModel.Item{
		ID:     "items_1",
		ItemID: "minecraft:stone",
		Maf: itemModel.ItemMaf{
			GrimoireID: "tempest01",
		},
		Minecraft: map[string]any{
			"components": map[string]any{
				"minecraft:consumable":  map[string]any{"consume_seconds": 10},
				"minecraft:custom_data": map[string]any{"legacy": true},
			},
		},
	}
	grimoiresByID := map[string]grimoireModel.Grimoire{
		"tempest01": {ID: "tempest01", MPCost: 1, CastTime: 1, CoolTime: 1, Title: "Spell"},
	}

	components, err := itemComponentsForLoot(entry, grimoiresByID, nil, nil)
	if err != nil {
		t.Fatalf("itemComponentsForLoot returned error: %v", err)
	}
	consumable, ok := components["minecraft:consumable"].(map[string]any)
	if !ok {
		t.Fatalf("consumable should be map: %#v", components["minecraft:consumable"])
	}
	if consumable["consume_seconds"] != 99999.0 {
		t.Fatalf("consume_seconds should be overwritten: %#v", consumable)
	}
	if _, ok := components["minecraft:custom_data"]; ok {
		t.Fatalf("custom_data should be excluded from set_components: %#v", components)
	}
}

func TestItemToGiveCommandPreservesEnchantmentsComponent(t *testing.T) {
	entry := itemModel.Item{
		ID:     "items_1",
		ItemID: "minecraft:stone",
		Minecraft: map[string]any{
			"components": map[string]any{
				"minecraft:enchantments": map[string]any{"minecraft:aqua_affinity": 1, "minecraft:bane_of_arthropods": 9},
			},
		},
	}

	command, err := ItemToGiveCommand(entry, nil, nil, nil)
	if err != nil {
		t.Fatalf("ItemToGiveCommand returned error: %v", err)
	}
	if !strings.Contains(command, `minecraft:enchantments={`) {
		t.Fatalf("enchantments component should be preserved: %s", command)
	}
}

func TestBowItemEmbedsBowAndPassiveIdsWithoutConsumable(t *testing.T) {
	entry := itemModel.Item{
		ID:     "bow_item",
		ItemID: "minecraft:bow",
		Maf: itemModel.ItemMaf{
			BowID: "test_full",
		},
		Minecraft: map[string]any{
			"components": map[string]any{
				"minecraft:custom_name": map[string]any{"text": "Bow Item"},
			},
		},
	}
	bowsByID := map[string]bowModel.BowPassive{
		"test_full": {ID: "test_full"},
	}

	customData, err := itemCustomData(entry, nil, nil, bowsByID)
	if err != nil {
		t.Fatalf("itemCustomData returned error: %v", err)
	}
	if !strings.Contains(customData, `bowId:"test_full"`) {
		t.Fatalf("bowId should be embedded: %s", customData)
	}
	if !strings.Contains(customData, `passiveId:"bow_test_full"`) {
		t.Fatalf("passiveId bridge should be embedded: %s", customData)
	}
	if !strings.Contains(customData, `passiveCondition:"always"`) {
		t.Fatalf("bow item should embed passiveCondition always: %s", customData)
	}
	if strings.Contains(customData, `hasPassive:1b`) || strings.Contains(customData, `passiveSlot:`) {
		t.Fatalf("bow item should not embed passive slot metadata: %s", customData)
	}
	if strings.Contains(customData, `spell:{`) {
		t.Fatalf("bow item should not embed spell metadata: %s", customData)
	}

	components, err := itemComponentsForLoot(entry, nil, nil, bowsByID)
	if err != nil {
		t.Fatalf("itemComponentsForLoot returned error: %v", err)
	}
	if _, ok := components["minecraft:consumable"]; ok {
		t.Fatalf("bow item should not become consumable: %#v", components)
	}
}

func TestCrossbowItemEmbedsBowAndPassiveIdsWithoutConsumable(t *testing.T) {
	entry := itemModel.Item{
		ID:     "crossbow_item",
		ItemID: "minecraft:crossbow",
		Maf: itemModel.ItemMaf{
			BowID: "test_full",
		},
	}
	bowsByID := map[string]bowModel.BowPassive{
		"test_full": {ID: "test_full"},
	}

	customData, err := itemCustomData(entry, nil, nil, bowsByID)
	if err != nil {
		t.Fatalf("itemCustomData returned error: %v", err)
	}
	if !strings.Contains(customData, `bowId:"test_full"`) {
		t.Fatalf("bowId should be embedded: %s", customData)
	}
	if !strings.Contains(customData, `passiveId:"bow_test_full"`) {
		t.Fatalf("passiveId bridge should be embedded: %s", customData)
	}
	if !strings.Contains(customData, `passiveCondition:"always"`) {
		t.Fatalf("crossbow item should embed passiveCondition always: %s", customData)
	}

	components, err := itemComponentsForLoot(entry, nil, nil, bowsByID)
	if err != nil {
		t.Fatalf("itemComponentsForLoot returned error: %v", err)
	}
	if _, ok := components["minecraft:consumable"]; ok {
		t.Fatalf("crossbow item should not become consumable: %#v", components)
	}
}

func TestBowItemRejectsHybridGrimoireMetadata(t *testing.T) {
	entry := itemModel.Item{
		ID:     "bow_hybrid",
		ItemID: "minecraft:bow",
		Maf: itemModel.ItemMaf{
			BowID:      "test_full",
			GrimoireID: "tempest01",
		},
	}
	grimoiresByID := map[string]grimoireModel.Grimoire{
		"tempest01": {ID: "tempest01", MPCost: 1, CastTime: 1, CoolTime: 1, Title: "Spell"},
	}
	bowsByID := map[string]bowModel.BowPassive{
		"test_full": {ID: "test_full"},
	}

	if _, err := ItemToGiveCommand(entry, grimoiresByID, nil, bowsByID); err == nil {
		t.Fatal("expected hybrid bow/grimoire item to be rejected")
	}
}
