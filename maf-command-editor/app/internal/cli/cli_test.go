package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"tools2/app/internal/config"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/export"
)

func TestRunValidateSuccess(t *testing.T) {
	cfg := writeFixtureConfig(t, true)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"validate"}, &stdout, &stderr, cfg)
	if code != 0 {
		t.Fatalf("code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "savedata validation ok:") {
		t.Fatalf("stdout = %s", stdout.String())
	}
}

func TestRunValidateFailure(t *testing.T) {
	cfg := writeFixtureConfig(t, false)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"validate"}, &stdout, &stderr, cfg)
	if code != 1 {
		t.Fatalf("code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "Referenced skill does not exist.") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func TestRunExportFailure(t *testing.T) {
	cfg := writeFixtureConfig(t, false)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"export"}, &stdout, &stderr, cfg)
	if code != 1 {
		t.Fatalf("code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "Savedata validation failed.") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func writeFixtureConfig(t *testing.T, valid bool) config.Config {
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
	writeJSON(t, settingsPath, settings)
	writeJSON(t, filepath.Join(root, "item.json"), items.ItemState{
		Items: []items.ItemEntry{{
			ID:      "items_1",
			ItemID:  "minecraft:apple",
			SkillID: "skill_1",
		}},
	})
	skillEntries := []skills.SkillEntry{{
		ID:          "skill_1",
		Name:        "Slash",
		Description: "Basic slash",
		Script:      "say slash",
	}}
	if !valid {
		skillEntries = []skills.SkillEntry{}
	}
	writeJSON(t, filepath.Join(root, "skill.json"), common.EntryState[skills.SkillEntry]{
		Entries: skillEntries,
	})

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
		MinecraftLootTableRoot: writeMinecraftLootTableRoot(t, root),
	}
}

func writeJSON[T any](t *testing.T, path string, value T) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeMinecraftLootTableRoot(t *testing.T, root string) string {
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
