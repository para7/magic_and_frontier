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
			{ID: testUUID("141"), ItemID: "minecraft:apple", Count: 1},
		}},
		GrimoireState: grimoire.GrimoireState{Entries: []grimoire.GrimoireEntry{
			{ID: testUUID("142"), CastID: 1, Script: "say cast", Title: "Spell"},
		}},
		SkillState: common.EntryState[skills.SkillEntry]{Entries: []skills.SkillEntry{{
			ID:     testUUID("143"),
			Name:   "Slash",
			Script: "say slash",
			ItemID: testUUID("999"),
		}}},
		EnemySkillState: common.EntryState[enemyskills.EnemySkillEntry]{Entries: []enemyskills.EnemySkillEntry{{
			ID:     testUUID("144"),
			Name:   "Roar",
			Script: "say roar",
		}}},
		TreasureState: common.EntryState[treasures.TreasureEntry]{Entries: []treasures.TreasureEntry{{
			ID:   testUUID("145"),
			Name: "Chest",
			LootPools: []treasures.DropRef{{
				Kind:   "item",
				RefID:  testUUID("141"),
				Weight: 1,
			}},
		}}},
		EnemyState: common.EntryState[enemies.EnemyEntry]{Entries: []enemies.EnemyEntry{{
			ID:          testUUID("146"),
			Name:        "Zombie",
			HP:          20,
			DropTableID: testUUID("404"),
			SpawnRule: enemies.SpawnRule{
				Origin:   enemies.Vec3{X: 0, Y: 64, Z: 0},
				Distance: enemies.Distance{Min: 0, Max: 16},
			},
		}}},
	}, "", fixedNow())

	if report.OK {
		t.Fatalf("expected validation failure")
	}
	if !strings.Contains(report.String(), "skill["+testUUID("143")+"].itemId: Referenced item does not exist.") {
		t.Fatalf("report = %s", report.String())
	}
	if !strings.Contains(report.String(), "enemy["+testUUID("146")+"].dropTableId: Referenced treasure does not exist.") {
		t.Fatalf("report = %s", report.String())
	}
}

func TestServiceExportDatapackRejectsInvalidSavedata(t *testing.T) {
	cfg := testConfig(t)
	writeJSONFile(t, cfg.ItemStatePath, items.ItemState{Items: []items.ItemEntry{
		{ID: testUUID("141"), ItemID: "minecraft:apple", Count: 1},
	}})
	writeJSONFile(t, cfg.SkillStatePath, map[string]any{
		"entries": []map[string]any{{
			"id":     testUUID("143"),
			"name":   "Slash",
			"script": "say slash",
			"itemId": testUUID("999"),
		}},
	})

	svc := NewService(cfg, Dependencies{Now: fixedNow})
	result := svc.ExportDatapack()
	if result.OK {
		t.Fatalf("expected export failure")
	}
	if result.Code != "VALIDATION_FAILED" {
		t.Fatalf("code = %q", result.Code)
	}
	if !strings.Contains(result.Details, "Referenced item does not exist.") {
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

func testUUID(suffix string) string {
	return "00000000-0000-4000-8000-000000000" + suffix
}
