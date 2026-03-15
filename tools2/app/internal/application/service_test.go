package application

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"tools2/app/internal/config"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/treasures"
)

func TestValidateBundleDetectsBrokenReferences(t *testing.T) {
	report := ValidateBundle(StateBundle{
		ItemState: items.ItemState{Items: []items.ItemEntry{
			{ID: "items_1", ItemID: "minecraft:apple", Count: 1, SkillID: "skill_9"},
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
			{ID: "treasure_1", Mode: "custom", TablePath: "maf:loot/test", LootPools: []treasures.DropRef{{Kind: "item", RefID: "items_1", Weight: 1}}},
			{ID: "treasure_2", Mode: "custom", TablePath: "maf:loot/test", LootPools: []treasures.DropRef{{Kind: "grimoire", RefID: "grimoire_1", Weight: 1}}},
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
	}, "", fixedNow())

	if report.OK {
		t.Fatalf("expected validation failure")
	}
	if !strings.Contains(report.String(), "item[items_1].skillId: Referenced skill does not exist.") {
		t.Fatalf("report = %s", report.String())
	}
	if !strings.Contains(report.String(), "grimoire[grimoire_2].castid: Cast ID is already used by grimoire_1.") {
		t.Fatalf("report = %s", report.String())
	}
	if !strings.Contains(report.String(), "treasure[treasure_2].tablePath: Custom loot table path is already used by treasure_1.") {
		t.Fatalf("report = %s", report.String())
	}
}

func TestServiceAllocateGrimoireIdentity(t *testing.T) {
	cfg := testConfig(t)
	svc := NewService(cfg, Dependencies{Now: fixedNow})
	id, castID, err := svc.AllocateGrimoireIdentity()
	if err != nil {
		t.Fatal(err)
	}
	if id != "grimoire_1" || castID != 1 {
		t.Fatalf("got %s %d", id, castID)
	}
	nextID, err := svc.AllocateID("items")
	if err != nil {
		t.Fatal(err)
	}
	if nextID != "items_1" {
		t.Fatalf("nextID = %s", nextID)
	}
}

func TestServiceExportDatapackRejectsInvalidSavedata(t *testing.T) {
	cfg := testConfig(t)
	writeJSONFile(t, cfg.ItemStatePath, items.ItemState{Items: []items.ItemEntry{
		{ID: "items_1", ItemID: "minecraft:apple", Count: 1, SkillID: "skill_999"},
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

func fixedNow() time.Time {
	return time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
}

func testConfig(t *testing.T) config.Config {
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
			"itemFunctionDir":       "data/maf/function/item",
			"itemLootDir":           "data/maf/loot_table/item",
			"spellFunctionDir":      "data/maf/function/grimoire",
			"spellLootDir":          "data/maf/loot_table/grimoire",
			"skillFunctionDir":      "data/maf/function/skill",
			"enemySkillFunctionDir": "data/maf/function/enemy_skill",
			"enemyFunctionDir":      "data/maf/function/enemy/spawn",
			"enemyLootDir":          "data/maf/loot_table/enemy",
			"treasureLootDir":       "data/maf/loot_table/treasure",
			"debugFunctionDir":      "data/maf/function/debug/give",
			"minecraftTagDir":       "data/minecraft/tags/function",
		},
	}
	writeJSONFile(t, settingsPath, settings)

	return config.Config{
		Port:                8787,
		ItemStatePath:       filepath.Join(root, "item-state.json"),
		GrimoireStatePath:   filepath.Join(root, "grimoire-state.json"),
		SkillStatePath:      filepath.Join(root, "skill-state.json"),
		EnemySkillStatePath: filepath.Join(root, "enemy-skill-state.json"),
		EnemyStatePath:      filepath.Join(root, "enemy-state.json"),
		TreasureStatePath:   filepath.Join(root, "treasure-state.json"),
		IDCounterStatePath:  filepath.Join(root, "id-counters.json"),
		ExportSettingsPath:  settingsPath,
	}
}

func writeJSONFile(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}
