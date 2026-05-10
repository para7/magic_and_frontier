package export_convert

import (
	"strings"
	"testing"

	activeModel "maf_command_editor/app/domain/model/active"
	bowModel "maf_command_editor/app/domain/model/bow"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func TestItemLootHelpersReadMinecraftComponents(t *testing.T) {
	entry := itemModel.Item{
		ID:     "items_1",
		ItemID: "minecraft:stone",
		Maf: itemModel.ItemMaf{
			ActiveID:  "tempest01",
			PassiveID: "regeneration",
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

	activesByID := map[string]activeModel.Active{
		"tempest01": {
			ID:          "tempest01",
			MPCost:      13,
			CastTime:    40,
			CoolTime:    20,
			Title:       "テンペスト",
			Description: "敵1体に雷を落とし周辺に特大ダメージ",
			Target:      "視線",
			Range:       10,
		},
	}
	passivesByID := map[string]passiveModel.Passive{
		"regeneration": {
			ID:          "regeneration",
			Name:        "いつでもリジェネ",
			Lore:        "HPが常時回復する",
			Condition:   "always",
			Slots:       []int{1},
			Description: "",
			Target:      "自分",
			Range:       0,
		},
	}

	customData, err := itemCustomData(entry, activesByID, passivesByID, nil)
	if err != nil {
		t.Fatalf("itemCustomData returned error: %v", err)
	}
	if !strings.Contains(customData, `item_id:"minecraft:stone"`) {
		t.Fatalf("item_id missing from custom data: %s", customData)
	}
	if !strings.Contains(customData, `active_id:"tempest01"`) {
		t.Fatalf("active_id missing from custom data: %s", customData)
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
	if strings.Contains(customData, `spell:{`) {
		t.Fatalf("item custom data should not embed runtime spell metadata: %s", customData)
	}
	if !strings.Contains(customData, `nbt_snapshot:{`) {
		t.Fatalf("nbt snapshot should be derived from components: %s", customData)
	}

	components, err := itemComponentsForLoot(entry, activesByID, passivesByID, nil)
	if err != nil {
		t.Fatalf("itemComponentsForLoot returned error: %v", err)
	}
	if _, ok := components["minecraft:custom_name"]; !ok {
		t.Fatalf("custom_name should be exported: %#v", components)
	}
	if _, ok := components["minecraft:lore"]; !ok {
		t.Fatalf("lore should be exported: %#v", components)
	}
	lore := loreLinesByKey(t, components, "minecraft:lore")
	if len(lore) != 17 {
		t.Fatalf("lore line count = %d, want 17: %#v", len(lore), lore)
	}
	for _, want := range []string{"Sample item", "アクティブスキル", "テンペスト", "target : 視線", "range  : 10", "パッシブスキル", "いつでもリジェネ", "HPが常時回復する", "MP     : 10"} {
		if !containsString(lore, want) {
			t.Fatalf("lore should contain %q: %#v", want, lore)
		}
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
			Name:      "いつでもリジェネ",
			Lore:      "HPが常時回復する",
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
	lore := loreLinesByKey(t, components, "minecraft:lore")
	if !containsString(lore, "パッシブスキル") || !containsString(lore, "いつでもリジェネ") || !containsString(lore, "HPが常時回復する") {
		t.Fatalf("passive-only item should include generated lore: %#v", lore)
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
			ActiveID: "tempest01",
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
	activesByID := map[string]activeModel.Active{
		"tempest01": {
			ID:          "tempest01",
			MPCost:      13,
			CastTime:    40,
			CoolTime:    20,
			Title:       "テンペスト",
			Description: "敵1体に雷を落とし周辺に特大ダメージ",
			Target:      "視線",
			Range:       10,
		},
	}

	command, err := ItemToGiveCommand(entry, activesByID, nil, nil)
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
	if !strings.Contains(command, `minecraft:lore=[{text:"Sample item"},{color:"white",font:"minecraft:uniform",italic:0b,text:""},{color:"light_purple",font:"minecraft:uniform",italic:0b,text:"アクティブスキル"}`) {
		t.Fatalf("generated active lore should be appended after JSON lore: %s", command)
	}
	if !strings.Contains(command, `{bold:1b,color:"white",font:"minecraft:uniform",italic:0b,text:"テンペスト"}`) {
		t.Fatalf("generated active title should be rendered as bold lore: %s", command)
	}
	if !strings.Contains(command, `{color:"white",font:"minecraft:uniform",italic:0b,text:"target : 視線"}`) {
		t.Fatalf("generated target lore should be rendered: %s", command)
	}
}

func TestItemToGiveCommandOverridesMinecraftConsumableForSpellItems(t *testing.T) {
	entry := itemModel.Item{
		ID:     "items_1",
		ItemID: "minecraft:stone",
		Maf: itemModel.ItemMaf{
			ActiveID: "tempest01",
		},
		Minecraft: map[string]any{
			"components": map[string]any{
				"minecraft:consumable": map[string]any{"consume_seconds": 10},
			},
		},
	}
	activesByID := map[string]activeModel.Active{
		"tempest01": {ID: "tempest01", MPCost: 1, CastTime: 1, CoolTime: 1, Title: "Spell"},
	}

	command, err := ItemToGiveCommand(entry, activesByID, nil, nil)
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
			ActiveID: "tempest01",
		},
		Minecraft: map[string]any{
			"components": map[string]any{
				"minecraft:consumable":  map[string]any{"consume_seconds": 10},
				"minecraft:custom_data": map[string]any{"legacy": true},
			},
		},
	}
	activesByID := map[string]activeModel.Active{
		"tempest01": {ID: "tempest01", MPCost: 1, CastTime: 1, CoolTime: 1, Title: "Spell"},
	}

	components, err := itemComponentsForLoot(entry, activesByID, nil, nil)
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
		"test_full": {ID: "test_full", Name: "Test Bow", Lore: "弓スキル", MPCost: 0},
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
	lore := loreLinesByKey(t, components, "minecraft:lore")
	if !containsString(lore, "パッシブスキル") || !containsString(lore, "Test Bow") || !containsString(lore, "弓スキル") || !containsString(lore, "MP     : 0") {
		t.Fatalf("bow item should include generated passive lore: %#v", lore)
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

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestBowItemRejectsHybridActiveMetadata(t *testing.T) {
	entry := itemModel.Item{
		ID:     "bow_hybrid",
		ItemID: "minecraft:bow",
		Maf: itemModel.ItemMaf{
			BowID:    "test_full",
			ActiveID: "tempest01",
		},
	}
	activesByID := map[string]activeModel.Active{
		"tempest01": {ID: "tempest01", MPCost: 1, CastTime: 1, CoolTime: 1, Title: "Spell"},
	}
	bowsByID := map[string]bowModel.BowPassive{
		"test_full": {ID: "test_full"},
	}

	if _, err := ItemToGiveCommand(entry, activesByID, nil, bowsByID); err == nil {
		t.Fatal("expected hybrid bow/active item to be rejected")
	}
}
