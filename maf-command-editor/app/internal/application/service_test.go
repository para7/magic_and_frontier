package application

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"tools2/app/internal/config"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/export"
)

func TestValidateBundleDetectsBrokenReferences(t *testing.T) {
	report := ValidateBundle(StateBundle{
		ItemState: items.ItemState{Items: []items.ItemEntry{
			{ID: "items_1", ItemID: "minecraft:apple", SkillID: "skill_9"},
		}},
		GrimoireState: grimoire.GrimoireState{Entries: []grimoire.GrimoireEntry{
			{ID: "grimoire_1", CastID: 1, CastTime: 10, MPCost: 5, Script: "say cast", Title: "Spell"},
			{ID: "grimoire_2", CastID: 1, CastTime: 10, MPCost: 5, Script: "say cast", Title: "Spell 2"},
		}},
		SkillState: common.EntryState[skills.SkillEntry]{Entries: []skills.SkillEntry{{
			ID:     "skill_1",
			Script: "say slash",
		}}},
		EnemySkillState: common.EntryState[enemyskills.EnemySkillEntry]{Entries: []enemyskills.EnemySkillEntry{{
			ID:     "enemyskill_1",
			Script: "say roar",
		}}},
		TreasureState: common.EntryState[treasures.TreasureEntry]{Entries: []treasures.TreasureEntry{
			{ID: "treasure_1", TablePath: "minecraft:chests/simple_dungeon", LootPools: []treasures.DropRef{{Kind: "item", RefID: "items_1", Weight: 1}}},
			{ID: "treasure_2", TablePath: "minecraft:chests/simple_dungeon", LootPools: []treasures.DropRef{{Kind: "grimoire", RefID: "grimoire_1", Weight: 1}}},
		}},
		LootTableState: common.EntryState[loottables.LootTableEntry]{Entries: []loottables.LootTableEntry{
			{ID: "loottable_1", LootPools: []treasures.DropRef{{Kind: "item", RefID: "items_1", Weight: 1}}},
		}},
		EnemyState: common.EntryState[enemies.EnemyEntry]{Entries: []enemies.EnemyEntry{{
			ID:            "enemy_1",
			MobType:       "minecraft:zombie",
			Name:          "Zombie",
			HP:            20,
			EnemySkillIDs: []string{"enemyskill_404"},
			DropMode:      "replace",
			Drops:         []enemies.DropRef{{Kind: "item", RefID: "items_404", Weight: 1}},
		}}},
	}, "", filepath.Join(repoRoot(t), "minecraft", "1.21.11", "loot_table"), fixedNow())

	if report.OK {
		t.Fatalf("expected validation failure")
	}
	if !strings.Contains(report.String(), "item[items_1].skillId: Referenced skill does not exist.") {
		t.Fatalf("report = %s", report.String())
	}
	if !strings.Contains(report.String(), "grimoire[grimoire_2].castid: Cast ID is already used by grimoire_1.") {
		t.Fatalf("report = %s", report.String())
	}
	if !strings.Contains(report.String(), "treasure[treasure_2].tablePath: Loot table path is already used by treasure_1.") {
		t.Fatalf("report = %s", report.String())
	}
}

func TestServiceExportDatapackRejectsInvalidSavedata(t *testing.T) {
	cfg := testConfig(t)
	writeJSONFile(t, cfg.ItemStatePath, items.ItemState{Items: []items.ItemEntry{
		{ID: "items_1", ItemID: "minecraft:apple", SkillID: "skill_999"},
	}})
	svc := NewService(cfg, Dependencies{Now: fixedNow})
	result := svc.ExportDatapack()
	if result.OK {
		t.Fatalf("expected export failure")
	}
	if result.Code != "VALIDATION_FAILED" {
		t.Fatalf("code = %q", result.Code)
	}
	if !strings.Contains(result.Details, "Referenced skill does not exist.") {
		t.Fatalf("details = %s", result.Details)
	}
}

func TestValidateCheckedInSavedata(t *testing.T) {
	cfg := repoSavedataConfig(t)
	svc := NewService(cfg, Dependencies{Now: fixedNow})

	report, err := svc.ValidateAll()
	if err != nil {
		t.Fatalf("ValidateAll() error = %v", err)
	}
	if !report.OK {
		t.Fatalf("checked-in savedata validation failed:\n%s", report.String())
	}
}

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
