package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
)

func TestBuildTreasureArtifacts(t *testing.T) {
	workspace := t.TempDir()
	sourceRoot := filepath.Join(workspace, "savedata", "loot_table")
	minecraftRoot := filepath.Join(workspace, "minecraft", "1.21.11", "loot_table")

	writeTestJSON(t, filepath.Join(sourceRoot, "minecraft", "chests", "abandoned_mineshaft.json"), map[string]any{
		"type": "minecraft:chest",
		"pools": []any{
			map[string]any{
				"rolls": 1.0,
				"entries": []any{
					map[string]any{"type": "maf:item", "name": "item_1", "weight": 5.0},
				},
			},
		},
	})
	writeTestJSON(t, filepath.Join(sourceRoot, "maf", "chests", "custom.json"), map[string]any{
		"type": "minecraft:chest",
		"pools": []any{
			map[string]any{
				"rolls": 1.0,
				"entries": []any{
					map[string]any{"type": "maf:grimoire", "name": "grimoire_1", "weight": 2.0},
				},
			},
		},
	})
	writeTestJSON(t, filepath.Join(sourceRoot, "maf", "spawn_table_patterns.json"), map[string]any{
		"overworld": []any{},
	})

	writeTestJSON(t, filepath.Join(minecraftRoot, "chests", "abandoned_mineshaft.json"), map[string]any{
		"type": "minecraft:chest",
		"pools": []any{
			map[string]any{
				"rolls":   1.0,
				"entries": []any{map[string]any{"type": "minecraft:item", "name": "minecraft:apple"}},
			},
		},
	})

	master := exportMasterStub{
		items: []itemModel.Item{
			{ID: "item_1", Minecraft: itemModel.MinecraftItem{ItemID: "minecraft:stone"}},
		},
		grimoires: []grimoireModel.Grimoire{
			{ID: "grimoire_1", CastTime: 20, CoolTime: 40, MPCost: 10, Script: []string{"say test"}, Title: "Test"},
		},
	}

	artifacts, err := BuildTreasureArtifacts(master, sourceRoot, minecraftRoot)
	if err != nil {
		t.Fatalf("BuildTreasureArtifacts returned error: %v", err)
	}
	if len(artifacts) != 2 {
		t.Fatalf("expected 2 artifacts, got %d", len(artifacts))
	}

	artifactsByKey := map[string]TreasureArtifact{}
	for _, entry := range artifacts {
		artifactsByKey[entry.Namespace+":"+entry.RelPath] = entry
	}

	minecraftArtifact, ok := artifactsByKey["minecraft:chests/abandoned_mineshaft"]
	if !ok {
		t.Fatalf("minecraft artifact not found: %#v", artifacts)
	}
	minecraftPools, ok := minecraftArtifact.LootTable["pools"].([]any)
	if !ok {
		t.Fatalf("minecraft pools must be array, got %T", minecraftArtifact.LootTable["pools"])
	}
	if len(minecraftPools) != 2 {
		t.Fatalf("expected base+custom pools, got %d", len(minecraftPools))
	}
	customPool, ok := minecraftPools[1].(map[string]any)
	if !ok {
		t.Fatalf("custom pool must be object, got %T", minecraftPools[1])
	}
	customEntries, ok := customPool["entries"].([]any)
	if !ok || len(customEntries) != 1 {
		t.Fatalf("custom pool entries must have one element")
	}
	customEntry, ok := customEntries[0].(map[string]any)
	if !ok {
		t.Fatalf("custom entry must be object, got %T", customEntries[0])
	}
	if customEntry["type"] != "minecraft:item" || customEntry["name"] != "minecraft:stone" {
		t.Fatalf("unexpected custom entry: %#v", customEntry)
	}

	mafArtifact, ok := artifactsByKey["maf:chests/custom"]
	if !ok {
		t.Fatalf("maf artifact not found: %#v", artifacts)
	}
	mafPools, ok := mafArtifact.LootTable["pools"].([]any)
	if !ok || len(mafPools) != 1 {
		t.Fatalf("maf pools must have one element")
	}
	mafPool, ok := mafPools[0].(map[string]any)
	if !ok {
		t.Fatalf("maf pool must be object, got %T", mafPools[0])
	}
	mafEntries, ok := mafPool["entries"].([]any)
	if !ok || len(mafEntries) != 1 {
		t.Fatalf("maf entries must have one element")
	}
	mafEntry, ok := mafEntries[0].(map[string]any)
	if !ok {
		t.Fatalf("maf entry must be object, got %T", mafEntries[0])
	}
	if mafEntry["type"] != "minecraft:item" || mafEntry["name"] != "minecraft:book" {
		t.Fatalf("unexpected maf converted entry: %#v", mafEntry)
	}
}

func TestBuildTreasureArtifactsMissingSourceRoot(t *testing.T) {
	artifacts, err := BuildTreasureArtifacts(exportMasterStub{}, filepath.Join(t.TempDir(), "missing"), t.TempDir())
	if err != nil {
		t.Fatalf("BuildTreasureArtifacts returned error: %v", err)
	}
	if len(artifacts) != 0 {
		t.Fatalf("expected 0 artifacts, got %d", len(artifacts))
	}
}

func TestWriteTreasureArtifacts(t *testing.T) {
	outputRoot := t.TempDir()
	artifacts := []TreasureArtifact{
		{
			Namespace: "minecraft",
			RelPath:   "chests/abandoned_mineshaft",
			LootTable: map[string]any{"type": "minecraft:chest", "pools": []any{}},
		},
		{
			Namespace: "maf",
			RelPath:   "chests/custom",
			LootTable: map[string]any{"type": "minecraft:chest", "pools": []any{}},
		},
	}

	if err := WriteTreasureArtifacts(outputRoot, artifacts); err != nil {
		t.Fatalf("WriteTreasureArtifacts returned error: %v", err)
	}

	assertJSONFileExists(t, filepath.Join(outputRoot, "data", "minecraft", "loot_table", "chests", "abandoned_mineshaft.json"))
	assertJSONFileExists(t, filepath.Join(outputRoot, "data", "maf", "loot_table", "chests", "custom.json"))
}

func writeTestJSON(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write json: %v", err)
	}
}

func assertJSONFileExists(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read output json %s: %v", path, err)
	}
	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("invalid json %s: %v", path, err)
	}
}
