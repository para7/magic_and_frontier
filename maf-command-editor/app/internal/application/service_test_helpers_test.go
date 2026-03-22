package application

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"tools2/app/internal/config"
	"tools2/app/internal/export"
)

func fixedNow() time.Time {
	return time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
}

func testConfig(t *testing.T) config.Config {
	t.Helper()
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	templatePath := filepath.Join(root, "pack-template.mcmeta")
	if err := os.WriteFile(templatePath, []byte("{\"pack\":{\"pack_format\":61,\"description\":\"test\"}}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	settings := export.ExportSettings{
		OutputRoot:       "./out",
		Namespace:        "maf",
		TemplatePackPath: "./pack-template.mcmeta",
		Paths: export.ExportPaths{
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
	writeJSONFile(t, settingsPath, settings)

	return config.Config{
		Port:                   8787,
		ItemStatePath:          filepath.Join(root, "item.json"),
		GrimoireStatePath:      filepath.Join(root, "grimoire.json"),
		SkillStatePath:         filepath.Join(root, "skill.json"),
		EnemySkillStatePath:    filepath.Join(root, "enemy_skill.json"),
		EnemyStatePath:         filepath.Join(root, "enemy.json"),
		SpawnTableStatePath:    filepath.Join(root, "spawn_table.json"),
		TreasureStatePath:      filepath.Join(root, "treasure.json"),
		LootTablesStatePath:    filepath.Join(root, "loottables.json"),
		IDCounterStatePath:     filepath.Join(root, "id_counters.json"),
		ExportSettingsPath:     settingsPath,
		MinecraftLootTableRoot: writeTestMinecraftLootTableRoot(t, root),
	}
}

func repoSavedataConfig(t *testing.T) config.Config {
	t.Helper()

	root := repoRoot(t)
	savedataDir := filepath.Join(root, "savedata")
	fixtureDir := t.TempDir()
	settingsPath := filepath.Join(fixtureDir, "export_settings.json")
	writeJSONFile(t, settingsPath, export.ExportSettings{
		OutputRoot:       "./out",
		Namespace:        "maf",
		TemplatePackPath: "./pack-template.mcmeta",
		Paths: export.ExportPaths{
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
	})

	return config.Config{
		Port:                   8787,
		ItemStatePath:          filepath.Join(savedataDir, "item.json"),
		GrimoireStatePath:      filepath.Join(savedataDir, "grimoire.json"),
		SkillStatePath:         filepath.Join(savedataDir, "skill.json"),
		EnemySkillStatePath:    filepath.Join(savedataDir, "enemy_skill.json"),
		EnemyStatePath:         filepath.Join(savedataDir, "enemy.json"),
		SpawnTableStatePath:    filepath.Join(savedataDir, "spawn_table.json"),
		TreasureStatePath:      filepath.Join(savedataDir, "treasure.json"),
		LootTablesStatePath:    filepath.Join(savedataDir, "loottables.json"),
		IDCounterStatePath:     filepath.Join(savedataDir, "id_counters.json"),
		ExportSettingsPath:     settingsPath,
		MinecraftLootTableRoot: filepath.Join(root, "minecraft", "1.21.11", "loot_table"),
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", ".."))
}

func writeJSONFile[T any](t *testing.T, path string, value T) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeTestMinecraftLootTableRoot(t *testing.T, root string) string {
	t.Helper()
	dir := filepath.Join(root, "minecraft", "1.21.11", "loot_table", "chests")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "simple_dungeon.json"), []byte("{\"type\":\"minecraft:generic\",\"pools\":[]}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return filepath.Join(root, "minecraft", "1.21.11", "loot_table")
}
