package master

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"maf-command-editor/app/internal/domain/common"
	"maf-command-editor/app/internal/domain/entity/enemies"
	"maf-command-editor/app/internal/domain/entity/enemyskills"
	"maf-command-editor/app/internal/domain/entity/grimoire"
	"maf-command-editor/app/internal/domain/entity/items"
	"maf-command-editor/app/internal/domain/entity/loottables"
	"maf-command-editor/app/internal/domain/entity/skills"
	"maf-command-editor/app/internal/domain/entity/spawntables"
	"maf-command-editor/app/internal/domain/entity/treasures"
)

func TestItemSaveSortsByID(t *testing.T) {
	deps := testDependencies(t)
	master, err := NewJSONMaster(deps)
	if err != nil {
		t.Fatalf("NewJSONMaster() error = %v", err)
	}

	entity := master.Items()
	for _, id := range []string{"items_2", "items_1"} {
		res := entity.Validate(items.SaveInput{ID: id, ItemID: "minecraft:stone"}, master)
		if !res.OK {
			t.Fatalf("Validate(%s) failed: %+v", id, res)
		}
		if err := entity.Create(*res.Entry, master); err != nil {
			t.Fatalf("Create(%s) error = %v", id, err)
		}
	}
	if err := entity.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	state, err := deps.ItemRepo.LoadState()
	if err != nil {
		t.Fatalf("LoadState() error = %v", err)
	}
	if len(state.Entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(state.Entries))
	}
	if got := state.Entries[0].ID; got != "items_1" {
		t.Fatalf("first id = %q, want items_1", got)
	}
	if got := state.Entries[1].ID; got != "items_2" {
		t.Fatalf("second id = %q, want items_2", got)
	}
}

func TestInvalidSavedataDoesNotBlockStartup(t *testing.T) {
	deps := testDependencies(t)
	if err := deps.GrimoireRepo.SaveState(common.EntryState[grimoire.GrimoireEntry]{Entries: []grimoire.GrimoireEntry{{
		ID:       "free-id",
		CastID:   0,
		CastTime: 1,
		MPCost:   0,
		Script:   "say x",
		Title:    "bad",
	}}}); err != nil {
		t.Fatalf("SaveState() error = %v", err)
	}

	master, err := NewJSONMaster(deps)
	if err != nil {
		t.Fatalf("NewJSONMaster() error = %v", err)
	}
	report := master.ValidateSavedAll()
	if report.OK {
		t.Fatalf("ValidateSavedAll() OK = true, want false")
	}
	contains := false
	for _, issue := range report.Issues {
		if issue.Entity == "grimoire" && issue.ID == "free-id" && strings.Contains(issue.Message, "Must satisfy gte 1.") {
			contains = true
			break
		}
	}
	if !contains {
		t.Fatalf("expected grimoire free-id castid issue in report, got: %s", report.String())
	}
}

func testDependencies(t *testing.T) Dependencies {
	t.Helper()
	dir := t.TempDir()

	itemRepo := common.StateRepository[items.ItemEntry]{Path: filepath.Join(dir, "savedata", "item.json")}
	grimoireRepo := common.StateRepository[grimoire.GrimoireEntry]{Path: filepath.Join(dir, "savedata", "grimoire.json")}
	skillRepo := common.StateRepository[skills.SkillEntry]{Path: filepath.Join(dir, "savedata", "skill.json")}
	enemySkillRepo := common.StateRepository[enemyskills.EnemySkillEntry]{Path: filepath.Join(dir, "savedata", "enemy_skill.json")}
	enemyRepo := common.StateRepository[enemies.EnemyEntry]{Path: filepath.Join(dir, "savedata", "enemy.json")}
	spawnRepo := common.StateRepository[spawntables.SpawnTableEntry]{Path: filepath.Join(dir, "savedata", "spawn_table.json")}
	treasureRepo := common.StateRepository[treasures.TreasureEntry]{Path: filepath.Join(dir, "savedata", "treasure.json")}
	loottableRepo := common.StateRepository[loottables.LootTableEntry]{Path: filepath.Join(dir, "savedata", "loottables.json")}

	if err := itemRepo.SaveState(common.EntryState[items.ItemEntry]{Entries: []items.ItemEntry{}}); err != nil {
		t.Fatalf("seed item: %v", err)
	}
	if err := grimoireRepo.SaveState(common.EntryState[grimoire.GrimoireEntry]{Entries: []grimoire.GrimoireEntry{}}); err != nil {
		t.Fatalf("seed grimoire: %v", err)
	}
	if err := skillRepo.SaveState(common.EntryState[skills.SkillEntry]{Entries: []skills.SkillEntry{}}); err != nil {
		t.Fatalf("seed skill: %v", err)
	}
	if err := enemySkillRepo.SaveState(common.EntryState[enemyskills.EnemySkillEntry]{Entries: []enemyskills.EnemySkillEntry{}}); err != nil {
		t.Fatalf("seed enemy skill: %v", err)
	}
	if err := enemyRepo.SaveState(common.EntryState[enemies.EnemyEntry]{Entries: []enemies.EnemyEntry{}}); err != nil {
		t.Fatalf("seed enemy: %v", err)
	}
	if err := spawnRepo.SaveState(common.EntryState[spawntables.SpawnTableEntry]{Entries: []spawntables.SpawnTableEntry{}}); err != nil {
		t.Fatalf("seed spawn table: %v", err)
	}
	if err := treasureRepo.SaveState(common.EntryState[treasures.TreasureEntry]{Entries: []treasures.TreasureEntry{}}); err != nil {
		t.Fatalf("seed treasure: %v", err)
	}
	if err := loottableRepo.SaveState(common.EntryState[loottables.LootTableEntry]{Entries: []loottables.LootTableEntry{}}); err != nil {
		t.Fatalf("seed loottable: %v", err)
	}

	return Dependencies{
		ItemRepo:               itemRepo,
		GrimoireRepo:           grimoireRepo,
		SkillRepo:              skillRepo,
		EnemySkillRepo:         enemySkillRepo,
		EnemyRepo:              enemyRepo,
		SpawnTableRepo:         spawnRepo,
		TreasureRepo:           treasureRepo,
		LootTableRepo:          loottableRepo,
		MinecraftLootTableRoot: filepath.Join(dir, "minecraft", "loot_table"),
		Now:                    func() time.Time { return time.Date(2026, 3, 23, 12, 0, 0, 0, time.UTC) },
	}
}
