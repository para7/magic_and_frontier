package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"maf-command-editor/app/internal/domain/entity/grimoire"
	"maf-command-editor/app/internal/domain/entity/items"
	"maf-command-editor/app/internal/domain/entity/treasures"
	"maf-command-editor/app/internal/domain/mcsource"
)

func TestGenerateItemOutputsUsesConfiguredLootDir(t *testing.T) {
	settings := ExportSettings{
		OutputRoot: t.TempDir(),
		Namespace:  "maf",
		Paths: ExportPaths{
			ItemFunctionDir: "data/maf/function/generated/item",
			ItemLootDir:     "data/maf/loot_table/generated/item",
		},
	}

	_, err := generateItemOutputs(settings, []items.ItemEntry{{
		ID:     "items_1",
		ItemID: "minecraft:apple",
	}})
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(settings.OutputRoot, settings.Paths.ItemFunctionDir, "items_1.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "loot give @s loot maf:generated/item/items_1") {
		t.Fatalf("item function should reference generated loot table path: %s", string(data))
	}
}

func TestGenerateGrimoireDebugFunctionsCreatesPerEntryFile(t *testing.T) {
	settings := ExportSettings{
		OutputRoot: t.TempDir(),
		Namespace:  "maf",
	}
	entries := []grimoire.GrimoireEntry{
		{
			ID:          "grimoire_1",
			CastID:      1,
			CastTime:    20,
			MPCost:      5,
			Title:       "Firebolt",
			Description: "Basic sample projectile spell.",
		},
	}

	count, err := generateGrimoireDebugFunctions(settings, entries)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("generated count = %d, want 1", count)
	}

	path := filepath.Join(settings.OutputRoot, "data", "maf", "function", "generated", "debug", "grimoire", "grimoire_1.mcfunction")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.TrimSpace(string(data))
	if !strings.HasPrefix(text, "give @s minecraft:written_book[") {
		t.Fatalf("debug file should use direct give command: %s", text)
	}
	if !strings.Contains(text, `custom_data={maf:{grimoire_id:"grimoire_1",spell:{castid:1,cost:5,cast:20,title:"Firebolt",description:"Basic sample projectile spell."}}}`) {
		t.Fatalf("debug file should include spell custom_data: %s", text)
	}
}

func TestExportDatapackKeepsStaticFilesAndWritesGeneratedTick(t *testing.T) {
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	settings := ExportSettings{
		OutputRoot:       "./out",
		Namespace:        "maf",
		TemplatePackPath: "./pack-template.mcmeta",
		Paths: ExportPaths{
			ItemFunctionDir:       "data/maf/function/generated/item",
			ItemLootDir:           "data/maf/loot_table/generated/item",
			SpellFunctionDir:      "data/maf/function/generated/grimoire",
			SpellLootDir:          "data/maf/loot_table/generated/grimoire",
			SkillFunctionDir:      "data/maf/function/generated/skill",
			EnemySkillFunctionDir: "data/maf/function/generated/enemy_skill",
			EnemyFunctionDir:      "data/maf/function/generated/enemy",
			EnemyLootDir:          "data/maf/loot_table/generated/enemy",
			TreasureLootDir:       "data/maf/loot_table/generated/treasure",
			LoottableLootDir:      "data/maf/loot_table/generated/loottable",
			DebugFunctionDir:      "data/maf/function/debug/give",
			MinecraftTagDir:       "data/minecraft/tags/function",
		},
	}
	writeExportSettingsFile(t, settingsPath, settings)

	outRoot := filepath.Join(root, "out")
	packPath := filepath.Join(outRoot, "pack.mcmeta")
	loadTagPath := filepath.Join(outRoot, "data", "minecraft", "tags", "function", "load.json")
	tickTagPath := filepath.Join(outRoot, "data", "minecraft", "tags", "function", "tick.json")
	lootSentinelPath := filepath.Join(outRoot, "data", "minecraft", "loot_table", "keep.json")
	debugSentinelPath := filepath.Join(outRoot, "data", "maf", "function", "debug", "give", "item", "keep.mcfunction")
	generatedFunctionStalePath := filepath.Join(outRoot, "data", "maf", "function", "generated", "legacy", "stale.mcfunction")
	generatedLootStalePath := filepath.Join(outRoot, "data", "maf", "loot_table", "generated", "legacy", "stale.json")

	if err := os.MkdirAll(filepath.Dir(packPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(loadTagPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(tickTagPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(lootSentinelPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(debugSentinelPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(generatedFunctionStalePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(generatedLootStalePath), 0o755); err != nil {
		t.Fatal(err)
	}

	staticPack := []byte("static-pack\n")
	staticLoadTag := []byte("{\"values\":[\"keep:load\"]}\n")
	staticTickTag := []byte("{\"values\":[\"keep:tick\"]}\n")
	if err := os.WriteFile(packPath, staticPack, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(loadTagPath, staticLoadTag, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tickTagPath, staticTickTag, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lootSentinelPath, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(debugSentinelPath, []byte("# keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(generatedFunctionStalePath, []byte("# stale\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(generatedLootStalePath, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := ExportDatapack(ExportParams{
		ExportSettingsPath: settingsPath,
	})
	if !result.OK {
		t.Fatalf("ExportDatapack() failed: %+v", result)
	}
	if result.Generated == nil || result.Generated.TotalFiles != 4 {
		t.Fatalf("TotalFiles = %+v, want 4 (generated tick + vh replacer tick + enemy skill main + grimoire selectexec)", result.Generated)
	}

	packAfter, err := os.ReadFile(packPath)
	if err != nil {
		t.Fatal(err)
	}
	loadAfter, err := os.ReadFile(loadTagPath)
	if err != nil {
		t.Fatal(err)
	}
	tickAfter, err := os.ReadFile(tickTagPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(packAfter) != string(staticPack) {
		t.Fatalf("pack.mcmeta should not be rewritten: %s", string(packAfter))
	}
	if string(loadAfter) != string(staticLoadTag) {
		t.Fatalf("load tag should not be rewritten: %s", string(loadAfter))
	}
	if string(tickAfter) != string(staticTickTag) {
		t.Fatalf("tick tag should not be rewritten: %s", string(tickAfter))
	}

	if _, err := os.Stat(lootSentinelPath); err != nil {
		t.Fatalf("minecraft loot_table outside generated should be kept: %v", err)
	}
	if _, err := os.Stat(debugSentinelPath); err != nil {
		t.Fatalf("debug path outside generated should be kept: %v", err)
	}
	if _, err := os.Stat(generatedFunctionStalePath); !os.IsNotExist(err) {
		t.Fatalf("generated function stale file should be removed, err=%v", err)
	}
	if _, err := os.Stat(generatedLootStalePath); !os.IsNotExist(err) {
		t.Fatalf("generated loot stale file should be removed, err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(outRoot, "data", "maf", "function", "generated", "tick.mcfunction")); err != nil {
		t.Fatalf("generated tick should exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outRoot, "data", "maf", "function", "generated", "vh", "replacer", "tick.mcfunction")); err != nil {
		t.Fatalf("generated vh replacer tick should exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outRoot, "data", "maf", "function", "generated", "enemy_skill", "main.mcfunction")); err != nil {
		t.Fatalf("generated enemy skill main should exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outRoot, "data", "maf", "function", "generated", "load.mcfunction")); !os.IsNotExist(err) {
		t.Fatalf("generated load should not exist, err=%v", err)
	}
}

func TestGenerateTreasureOutputsRemovesStaleOverrides(t *testing.T) {
	root := t.TempDir()
	settings := ExportSettings{
		OutputRoot: root,
		Namespace:  "maf",
	}
	sourceRoot := filepath.Join(root, "minecraft-source")
	writeBaseLootTable := func(tablePath string) {
		t.Helper()
		path, err := mcsource.FilePathForTable(sourceRoot, tablePath)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := writeJSON(path, map[string]any{"type": "minecraft:generic", "pools": []any{}}); err != nil {
			t.Fatal(err)
		}
	}
	writeBaseLootTable("minecraft:chests/old")
	writeBaseLootTable("minecraft:chests/new")

	firstEntries := []treasures.TreasureEntry{
		{
			ID:        "treasure_1",
			TablePath: "minecraft:chests/old",
			LootPools: []treasures.DropRef{
				{Kind: "minecraft_item", RefID: "minecraft:apple", Weight: 1},
			},
		},
	}
	if _, err := generateTreasureOutputs(settings, sourceRoot, firstEntries, nil, nil); err != nil {
		t.Fatal(err)
	}
	oldPath, err := lootTableOutputPath(settings, "minecraft:chests/old")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(oldPath); err != nil {
		t.Fatalf("first export should write old override: %v", err)
	}

	secondEntries := []treasures.TreasureEntry{
		{
			ID:        "treasure_1",
			TablePath: "minecraft:chests/new",
			LootPools: []treasures.DropRef{
				{Kind: "minecraft_item", RefID: "minecraft:apple", Weight: 1},
			},
		},
	}
	if _, err := generateTreasureOutputs(settings, sourceRoot, secondEntries, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("stale override should be removed, err=%v", err)
	}
	newPath, err := lootTableOutputPath(settings, "minecraft:chests/new")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("new override should exist: %v", err)
	}

	manifestPath := treasureOverlayManifestPath(settings)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("manifest should exist: %v", err)
	}
	var manifest treasureOverlayManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("manifest should be valid json: %v", err)
	}
	expectedRel, err := filepath.Rel(settings.OutputRoot, newPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Paths) != 1 || manifest.Paths[0] != filepath.ToSlash(expectedRel) {
		t.Fatalf("manifest paths = %#v, want [%q]", manifest.Paths, filepath.ToSlash(expectedRel))
	}
}
