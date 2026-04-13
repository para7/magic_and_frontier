package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	bowModel "maf_command_editor/app/domain/model/bow"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
	config "maf_command_editor/app/files"
)

func TestBuildItemArtifactsBuildsGiveCommands(t *testing.T) {
	master := exportMasterStub{
		grimoires: []grimoireModel.Grimoire{
			{ID: "tempest01", MPCost: 13, CastTime: 40, CoolTime: 20, Title: "テンペスト", Description: "敵1体に雷を落とし周辺に特大ダメージ"},
		},
		passives: []passiveModel.Passive{
			{ID: "regeneration", Condition: "always", Slots: []int{1}},
		},
		items: []itemModel.Item{
			{
				ID: "items_1",
				Maf: itemModel.ItemMaf{
					GrimoireID: "tempest01",
					PassiveID:  "regeneration",
				},
				Minecraft: itemModel.MinecraftItem{
					ItemID: "minecraft:stone",
					Components: map[string]string{
						"minecraft:custom_name": `'{"text":"Starter Stone"}'`,
						"minecraft:lore":        `['{"text":"Sample item"}']`,
					},
				},
			},
		},
	}

	artifacts, err := BuildItemArtifacts(master)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 1 {
		t.Fatalf("artifacts length = %d, want 1", len(artifacts))
	}
	if artifacts[0].ID != "items_1" {
		t.Fatalf("unexpected artifact id: %#v", artifacts[0])
	}
	if !strings.Contains(artifacts[0].Body, `give @p minecraft:stone[`) {
		t.Fatalf("unexpected give command: %q", artifacts[0].Body)
	}
	if !strings.Contains(artifacts[0].Body, `minecraft:consumable={consume_seconds:99999,animation:"bow",has_consume_particles:false}`) {
		t.Fatalf("spell item should include consumable: %q", artifacts[0].Body)
	}
	if !strings.Contains(artifacts[0].Body, `minecraft:custom_data={maf:{`) {
		t.Fatalf("give command should include custom_data: %q", artifacts[0].Body)
	}
}

func TestBuildItemArtifactsBuildsBowGiveCommand(t *testing.T) {
	master := exportMasterStub{
		bows: []bowModel.BowPassive{
			{ID: "test_full"},
		},
		items: []itemModel.Item{
			{
				ID: "bow_item",
				Maf: itemModel.ItemMaf{
					BowID: "test_full",
				},
				Minecraft: itemModel.MinecraftItem{
					ItemID: "minecraft:bow",
				},
			},
		},
	}

	artifacts, err := BuildItemArtifacts(master)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 1 {
		t.Fatalf("artifacts length = %d, want 1", len(artifacts))
	}
	if !strings.Contains(artifacts[0].Body, `bowId:"test_full"`) {
		t.Fatalf("bowId should be exported: %q", artifacts[0].Body)
	}
	if !strings.Contains(artifacts[0].Body, `passiveId:"bow_test_full"`) {
		t.Fatalf("passive bridge id should be exported: %q", artifacts[0].Body)
	}
}

func TestBuildItemArtifactsBuildsCrossbowGiveCommand(t *testing.T) {
	master := exportMasterStub{
		bows: []bowModel.BowPassive{
			{ID: "test_full"},
		},
		items: []itemModel.Item{
			{
				ID: "crossbow_item",
				Maf: itemModel.ItemMaf{
					BowID: "test_full",
				},
				Minecraft: itemModel.MinecraftItem{
					ItemID: "minecraft:crossbow",
				},
			},
		},
	}

	artifacts, err := BuildItemArtifacts(master)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 1 {
		t.Fatalf("artifacts length = %d, want 1", len(artifacts))
	}
	if !strings.Contains(artifacts[0].Body, `give @p minecraft:crossbow[`) {
		t.Fatalf("crossbow item should export as crossbow: %q", artifacts[0].Body)
	}
	if !strings.Contains(artifacts[0].Body, `bowId:"test_full"`) {
		t.Fatalf("bowId should be exported: %q", artifacts[0].Body)
	}
	if !strings.Contains(artifacts[0].Body, `passiveId:"bow_test_full"`) {
		t.Fatalf("passive bridge id should be exported: %q", artifacts[0].Body)
	}
}

func TestWriteItemArtifactsWritesFiles(t *testing.T) {
	root := t.TempDir()
	artifacts := []ItemGiveFunction{
		{ID: "items_1", Body: "give @p minecraft:stone 1"},
	}

	if err := WriteItemArtifacts(root, artifacts); err != nil {
		t.Fatal(err)
	}

	body, err := os.ReadFile(filepath.Join(root, "items_1.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "give @p minecraft:stone 1\n" {
		t.Fatalf("unexpected item give body: %q", string(body))
	}
}

func TestExportDatapackWritesItemArtifacts(t *testing.T) {
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	settings := map[string]any{
		"outputRoot": filepath.Join(root, "out"),
		"exportPaths": map[string]any{
			"grimoireEffect": "generated/grimoire/effect",
			"grimoireDebug":  "generated/grimoire/give",
			"itemGive":       "generated/item/give",
			"enemy":          "generated/enemy/spawn",
			"enemySkill":     "generated/enemy/skill",
			"enemyLoot":      "generated/enemy/loot",
		},
	}
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := config.LoadConfig()
	cfg.ExportSettingsPath = settingsPath
	cfg.MinecraftLootTableRoot = filepath.Join(root, "minecraft", "loot_table")

	master := exportMasterStub{
		items: []itemModel.Item{
			{
				ID: "items_1",
				Minecraft: itemModel.MinecraftItem{
					ItemID: "minecraft:stone",
					Components: map[string]string{
						"minecraft:custom_name": `'{"text":"Starter Stone"}'`,
					},
				},
			},
		},
	}

	if err := ExportDatapack(master, cfg); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(root, "out", "data", "maf", "function", "generated", "item", "give", "items_1.mcfunction")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `give @p minecraft:stone[`) {
		t.Fatalf("unexpected exported item give body: %q", string(body))
	}
	if !strings.Contains(string(body), `minecraft:custom_data={maf:{`) {
		t.Fatalf("item give file should contain custom_data: %q", string(body))
	}
}

func TestExportDatapackUsesDefaultItemGivePathWhenUnset(t *testing.T) {
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	settings := map[string]any{
		"outputRoot": filepath.Join(root, "out"),
		"exportPaths": map[string]any{
			"grimoireEffect": "generated/grimoire/effect",
			"grimoireDebug":  "generated/grimoire/give",
			"enemy":          "generated/enemy/spawn",
			"enemySkill":     "generated/enemy/skill",
			"enemyLoot":      "generated/enemy/loot",
		},
	}
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := config.LoadConfig()
	cfg.ExportSettingsPath = settingsPath
	cfg.MinecraftLootTableRoot = filepath.Join(root, "minecraft", "loot_table")

	master := exportMasterStub{
		items: []itemModel.Item{
			{
				ID:        "items_1",
				Minecraft: itemModel.MinecraftItem{ItemID: "minecraft:stone"},
			},
		},
	}

	if err := ExportDatapack(master, cfg); err != nil {
		t.Fatal(err)
	}

	defaultPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "item", "give", "items_1.mcfunction")
	if _, err := os.Stat(defaultPath); err != nil {
		t.Fatalf("default item give file should exist: %v", err)
	}
}
