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
	if !strings.Contains(customData, `nbt_snapshot:"{`) {
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

func TestItemToGiveCommandBuildsSortedComponentsAndCustomData(t *testing.T) {
	entry := itemModel.Item{
		ID: "items_1",
		Maf: itemModel.ItemMaf{
			GrimoireID: "tempest01",
		},
		Minecraft: itemModel.MinecraftItem{
			ItemID: "minecraft:stone",
			Components: map[string]string{
				"minecraft:lore":        `['{"text":"Sample item"}']`,
				"minecraft:custom_name": `'{"text":"Starter Stone"}'`,
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
	if !strings.Contains(command, `minecraft:custom_name={"text":"Starter Stone"}`) {
		t.Fatalf("custom_name should be normalized for give: %s", command)
	}
	if !strings.Contains(command, `minecraft:lore=[{"text":"Sample item"}]`) {
		t.Fatalf("lore should be normalized for give: %s", command)
	}
	if strings.Contains(command, `minecraft:custom_name='{"text":"Starter Stone"}'`) {
		t.Fatalf("quoted custom_name JSON should not remain in give output: %s", command)
	}
}

func TestItemToGiveCommandNormalizesRawJSONTextComponents(t *testing.T) {
	entry := itemModel.Item{
		ID: "items_raw_json",
		Minecraft: itemModel.MinecraftItem{
			ItemID: "minecraft:stone",
			Components: map[string]string{
				"minecraft:item_name": `{"text":"Debug Title","italic":false}`,
				"minecraft:lore":      `[{"text":"Role line"},{"text":"Extra"}]`,
			},
		},
	}

	command, err := ItemToGiveCommand(entry, nil, nil, nil)
	if err != nil {
		t.Fatalf("ItemToGiveCommand returned error: %v", err)
	}
	if !strings.Contains(command, `minecraft:item_name={"italic":false,"text":"Debug Title"}`) &&
		!strings.Contains(command, `minecraft:item_name={"text":"Debug Title","italic":false}`) {
		t.Fatalf("item_name should be preserved as a structured text component: %s", command)
	}
	if !strings.Contains(command, `minecraft:lore=[{"text":"Role line"},{"text":"Extra"}]`) {
		t.Fatalf("lore should remain structured JSON for give: %s", command)
	}
}

func TestItemToGiveCommandPreservesDirectGiveTextComponentSNBT(t *testing.T) {
	entry := itemModel.Item{
		ID: "items_snbt",
		Minecraft: itemModel.MinecraftItem{
			ItemID: "minecraft:stone",
			Components: map[string]string{
				"minecraft:custom_name": `{text:"Starter Stone"}`,
				"minecraft:lore":        `[{text:"Sample item"}]`,
			},
		},
	}

	command, err := ItemToGiveCommand(entry, nil, nil, nil)
	if err != nil {
		t.Fatalf("ItemToGiveCommand returned error: %v", err)
	}
	if !strings.Contains(command, `minecraft:custom_name={text:"Starter Stone"}`) {
		t.Fatalf("direct give custom_name SNBT should be preserved: %s", command)
	}
	if !strings.Contains(command, `minecraft:lore=[{text:"Sample item"}]`) {
		t.Fatalf("direct give lore SNBT should be preserved: %s", command)
	}
}

func TestItemToGiveCommandDoesNotDuplicateConsumable(t *testing.T) {
	entry := itemModel.Item{
		ID: "items_1",
		Maf: itemModel.ItemMaf{
			GrimoireID: "tempest01",
		},
		Minecraft: itemModel.MinecraftItem{
			ItemID: "minecraft:stone",
			Components: map[string]string{
				"minecraft:consumable": `{consume_seconds:10}`,
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
		t.Fatalf("consumable should not be duplicated: %s", command)
	}
	if !strings.Contains(command, `minecraft:consumable={consume_seconds:10}`) {
		t.Fatalf("existing consumable should be preserved: %s", command)
	}
}

func TestItemToGiveCommandPreservesEnchantmentsComponent(t *testing.T) {
	entry := itemModel.Item{
		ID: "items_1",
		Minecraft: itemModel.MinecraftItem{
			ItemID: "minecraft:stone",
			Components: map[string]string{
				"minecraft:enchantments": `{"minecraft:aqua_affinity":1,"minecraft:bane_of_arthropods":9}`,
			},
		},
	}

	command, err := ItemToGiveCommand(entry, nil, nil, nil)
	if err != nil {
		t.Fatalf("ItemToGiveCommand returned error: %v", err)
	}
	if !strings.Contains(command, `minecraft:enchantments={"minecraft:aqua_affinity":1,"minecraft:bane_of_arthropods":9}`) {
		t.Fatalf("enchantments component should be preserved: %s", command)
	}
}

func TestBowItemEmbedsBowAndPassiveIdsWithoutConsumable(t *testing.T) {
	entry := itemModel.Item{
		ID: "bow_item",
		Maf: itemModel.ItemMaf{
			BowID: "test_full",
		},
		Minecraft: itemModel.MinecraftItem{
			ItemID: "minecraft:bow",
			Components: map[string]string{
				"minecraft:custom_name": `'{"text":"Bow Item"}'`,
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
	if strings.Contains(customData, `hasPassive:1b`) || strings.Contains(customData, `passiveSlot:`) || strings.Contains(customData, `passiveCondition:`) {
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
		ID: "crossbow_item",
		Maf: itemModel.ItemMaf{
			BowID: "test_full",
		},
		Minecraft: itemModel.MinecraftItem{
			ItemID: "minecraft:crossbow",
			Components: map[string]string{
				"minecraft:custom_name": `'{"text":"Crossbow Item"}'`,
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
		ID: "bow_hybrid",
		Maf: itemModel.ItemMaf{
			BowID:      "test_full",
			GrimoireID: "tempest01",
		},
		Minecraft: itemModel.MinecraftItem{
			ItemID: "minecraft:bow",
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
