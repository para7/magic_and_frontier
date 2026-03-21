package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"tools2/app/internal/config"
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
	settingsPath := filepath.Join(root, "export-settings.json")
	templatePath := filepath.Join(root, "pack-template.mcmeta")
	if err := os.WriteFile(templatePath, []byte("{\"pack\":{\"pack_format\":61,\"description\":\"test\"}}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	settings := map[string]any{
		"outputRoot":       "./out",
		"namespace":        "maf",
		"templatePackPath": "./pack-template.mcmeta",
		"paths": map[string]any{
			"itemFunctionDir":       "data/maf/function/generated/item",
			"itemLootDir":           "data/maf/loot_table/generated/item",
			"spellFunctionDir":      "data/maf/function/generated/grimoire",
			"spellLootDir":          "data/maf/loot_table/generated/grimoire",
			"skillFunctionDir":      "data/maf/function/generated/skill",
			"enemySkillFunctionDir": "data/maf/function/generated/enemy_skill",
			"enemyFunctionDir":      "data/maf/function/generated/enemy",
			"enemyLootDir":          "data/maf/loot_table/generated/enemy",
			"treasureLootDir":       "data/maf/loot_table/generated/treasure",
			"loottableLootDir":      "data/maf/loot_table/generated/loottable",
			"debugFunctionDir":      "data/maf/function/debug/give",
			"minecraftTagDir":       "data/minecraft/tags/function",
		},
	}
	writeJSON(t, settingsPath, settings)
	writeJSON(t, filepath.Join(root, "item.json"), map[string]any{
		"items": []map[string]any{{
			"id":      "items_1",
			"itemId":  "minecraft:apple",
			"count":   1,
			"skillId": "skill_1",
		}},
	})
	skillEntries := []map[string]any{{
		"id":          "skill_1",
		"name":        "Slash",
		"description": "Basic slash",
		"script":      "say slash",
	}}
	if !valid {
		skillEntries = []map[string]any{}
	}
	writeJSON(t, filepath.Join(root, "skill.json"), map[string]any{
		"entries": skillEntries,
	})

	return config.Config{
		Port:                   8787,
		ItemStatePath:          filepath.Join(root, "item.json"),
		GrimoireStatePath:      filepath.Join(root, "grimoire.json"),
		SkillStatePath:         filepath.Join(root, "skill.json"),
		EnemySkillStatePath:    filepath.Join(root, "enemy-skill.json"),
		EnemyStatePath:         filepath.Join(root, "enemy.json"),
		TreasureStatePath:      filepath.Join(root, "treasure.json"),
		LootTablesStatePath:    filepath.Join(root, "loottables.json"),
		IDCounterStatePath:     filepath.Join(root, "id-counters.json"),
		ExportSettingsPath:     settingsPath,
		MinecraftLootTableRoot: writeMinecraftLootTableRoot(t, root),
	}
}

func writeJSON(t *testing.T, path string, value any) {
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
