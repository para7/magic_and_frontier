package application

import (
	"path/filepath"
	"strings"
	"testing"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/treasures"
)

func TestValidateBundleDetectsBrokenReferences(t *testing.T) {
	report := ValidateBundle(StateBundle{
		ItemState: common.EntryState[items.ItemEntry]{Entries: []items.ItemEntry{
			{ID: "items_1", ItemID: "minecraft:apple", SkillID: "skill_9"},
		}},
		GrimoireState: common.EntryState[grimoire.GrimoireEntry]{Entries: []grimoire.GrimoireEntry{
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
