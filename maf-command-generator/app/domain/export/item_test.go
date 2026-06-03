package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	activeModel "maf_command_editor/app/domain/model/active"
	bowModel "maf_command_editor/app/domain/model/bow"
	itemModel "maf_command_editor/app/domain/model/item"
	passiveModel "maf_command_editor/app/domain/model/passive"
	config "maf_command_editor/app/files"
)

func TestItemExportFixtures(t *testing.T) {
	cases := discoverCases(t, filepath.Join("testdata", "item"))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			master := loadFixtureMaster(t, tc.dir)
			artifacts, err := BuildItemArtifacts(master)
			if err != nil {
				t.Fatal(err)
			}

			actualDir := t.TempDir()
			if err := WriteItemArtifacts(filepath.Join(actualDir, "give"), artifacts); err != nil {
				t.Fatal(err)
			}

			assertGoldenDir(t, filepath.Join(tc.dir, "output"), actualDir)
		})
	}
}

func TestBuildItemArtifactsBuildsGiveCommands(t *testing.T) {
	master := exportMasterStub{
		actives: []activeModel.Active{
			{ID: "tempest01", MPCost: 13, CastTime: 40, CoolTime: 20, Title: "テンペスト", Description: "敵1体に雷を落とし周辺に特大ダメージ"},
		},
		passives: []passiveModel.Passive{
			{ID: "regeneration", Condition: "always", Slots: []int{1}},
		},
		items: []itemModel.Item{
			{
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
				ID:     "bow_item",
				ItemID: "minecraft:bow",
				Maf: itemModel.ItemMaf{
					BowID: "test_full",
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
				ID:     "crossbow_item",
				ItemID: "minecraft:crossbow",
				Maf: itemModel.ItemMaf{
					BowID: "test_full",
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

func TestBuildItemUpdateArtifactsBuildsUpdateCommands(t *testing.T) {
	master := exportMasterStub{
		items: []itemModel.Item{
			{
				ID:     "items_1",
				ItemID: "minecraft:stone",
			},
		},
	}

	artifacts, err := BuildItemUpdateArtifacts(master)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 1 {
		t.Fatalf("artifacts length = %d, want 1", len(artifacts))
	}
	if artifacts[0].ID != "items_1" {
		t.Fatalf("unexpected artifact id: %#v", artifacts[0])
	}
	if !strings.Contains(artifacts[0].Body, `$item replace entity @s $(slot) with minecraft:stone[`) {
		t.Fatalf("unexpected update command: %q", artifacts[0].Body)
	}
}

func TestWriteItemUpdateArtifactsWritesTimestampAndFiles(t *testing.T) {
	root := t.TempDir()
	artifacts := []ItemUpdateFunction{
		{ID: "items_1", Body: "$item replace entity @s $(slot) with minecraft:stone 1"},
	}

	if err := WriteItemUpdateArtifacts(root, 12345, artifacts); err != nil {
		t.Fatal(err)
	}

	timestampBody, err := os.ReadFile(filepath.Join(root, "timestamp.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if string(timestampBody) != "scoreboard players set #maf_item_ver tmp 12345\n" {
		t.Fatalf("unexpected timestamp body: %q", string(timestampBody))
	}

	body, err := os.ReadFile(filepath.Join(root, "items_1.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "$item replace entity @s $(slot) with minecraft:stone 1\n" {
		t.Fatalf("unexpected item update body: %q", string(body))
	}
}

func TestExportDatapackWritesItemArtifacts(t *testing.T) {
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	settings := map[string]any{
		"outputRoot": filepath.Join(root, "out"),
		"exportPaths": map[string]any{
			"activeEffect": "generated/active/effect",
			"activeDebug":  "generated/active/give",
			"itemGive":     "generated/item/give",
			"itemUpdate":   "generated/item/update",
			"enemy":        "generated/enemy/spawn",
			"enemySkill":   "generated/enemy/skill",
			"enemyLoot":    "generated/enemy/loot",
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
				ID:     "items_1",
				ItemID: "minecraft:stone",
				Minecraft: map[string]any{
					"components": map[string]any{
						"minecraft:custom_name": map[string]any{"text": "Starter Stone"},
					},
				},
			},
		},
	}

	restoreExportUnixNow(t, 12345)
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

	updatePath := filepath.Join(root, "out", "data", "maf", "function", "generated", "item", "update", "items_1.mcfunction")
	updateBody, err := os.ReadFile(updatePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(updateBody), `$item replace entity @s $(slot) with minecraft:stone[`) {
		t.Fatalf("unexpected exported item update body: %q", string(updateBody))
	}

	timestampPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "item", "update", "timestamp.mcfunction")
	timestampBody, err := os.ReadFile(timestampPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(timestampBody) != "scoreboard players set #maf_item_ver tmp 12345\n" {
		t.Fatalf("unexpected exported timestamp body: %q", string(timestampBody))
	}
}

func TestExportDatapackUsesDefaultItemGivePathWhenUnset(t *testing.T) {
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	settings := map[string]any{
		"outputRoot": filepath.Join(root, "out"),
		"exportPaths": map[string]any{
			"activeEffect": "generated/active/effect",
			"activeDebug":  "generated/active/give",
			"enemy":        "generated/enemy/spawn",
			"enemySkill":   "generated/enemy/skill",
			"enemyLoot":    "generated/enemy/loot",
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
				ID:     "items_1",
				ItemID: "minecraft:stone",
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
