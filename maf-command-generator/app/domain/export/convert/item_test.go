package export_convert

import (
	"strings"
	"testing"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func TestItemLootHelpersReadMinecraftComponents(t *testing.T) {
	entry := itemModel.Item{
		ID: "items_1",
		Maf: itemModel.ItemMaf{
			GrimoireID: "tempest01",
			PassiveID:  "regeneration",
		},
		Minecraft: itemModel.MinecraftItem{
			ItemID: "minecraft:stone",
			Components: map[string]string{
				"minecraft:custom_name":  `'{"text":"Starter Stone"}'`,
				"minecraft:lore":         `['{"text":"Sample item"}']`,
				"minecraft:unbreakable":  `{}`,
				"minecraft:enchantments": `{"minecraft:sharpness":5}`,
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

	customData, err := itemCustomData(entry, grimoiresByID, passivesByID)
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
	if !strings.Contains(customData, `nbt_snapshot:"{`) {
		t.Fatalf("nbt snapshot should be derived from components: %s", customData)
	}

	components, err := itemComponentsForLoot(entry, grimoiresByID, passivesByID)
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
	if _, ok := components["minecraft:consumable"]; !ok {
		t.Fatalf("consumable should be added for spell items: %#v", components)
	}

	enchantments := itemEnchantmentsForLoot(entry)
	if enchantments["minecraft:sharpness"] != float64(5) {
		t.Fatalf("unexpected enchantments: %#v", enchantments)
	}
}

func TestPassiveOnlyItemDoesNotBecomeRightClickSpell(t *testing.T) {
	entry := itemModel.Item{
		ID: "items_passive_only",
		Maf: itemModel.ItemMaf{
			PassiveID: "regeneration",
		},
		Minecraft: itemModel.MinecraftItem{
			ItemID: "minecraft:stone",
			Components: map[string]string{
				"minecraft:custom_name": `'{"text":"Passive Only"}'`,
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

	customData, err := itemCustomData(entry, nil, passivesByID)
	if err != nil {
		t.Fatalf("itemCustomData returned error: %v", err)
	}
	if !strings.Contains(customData, `hasPassive:1b`) || !strings.Contains(customData, `passiveId:"regeneration"`) {
		t.Fatalf("passive metadata should be embedded: %s", customData)
	}
	if strings.Contains(customData, `spell:{`) {
		t.Fatalf("passive-only item should not embed spell metadata: %s", customData)
	}

	components, err := itemComponentsForLoot(entry, nil, passivesByID)
	if err != nil {
		t.Fatalf("itemComponentsForLoot returned error: %v", err)
	}
	if _, ok := components["minecraft:consumable"]; ok {
		t.Fatalf("passive-only item should not be consumable: %#v", components)
	}
}
