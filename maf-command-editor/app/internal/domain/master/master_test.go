package master

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/store"
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

	state, err := deps.ItemRepo.LoadItemState()
	if err != nil {
		t.Fatalf("LoadItemState() error = %v", err)
	}
	if len(state.Items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(state.Items))
	}
	if got := state.Items[0].ID; got != "items_1" {
		t.Fatalf("first id = %q, want items_1", got)
	}
	if got := state.Items[1].ID; got != "items_2" {
		t.Fatalf("second id = %q, want items_2", got)
	}
}

func TestInvalidSavedataDoesNotBlockStartup(t *testing.T) {
	deps := testDependencies(t)
	if err := deps.GrimoireRepo.SaveGrimoireState(grimoire.GrimoireState{Entries: []grimoire.GrimoireEntry{{
		ID:       "free-id",
		CastID:   0,
		CastTime: 1,
		MPCost:   0,
		Script:   "say x",
		Title:    "bad",
	}}}); err != nil {
		t.Fatalf("SaveGrimoireState() error = %v", err)
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

	itemRepo := store.NewItemStateRepository(filepath.Join(dir, "savedata", "item.json"))
	grimoireRepo := store.NewGrimoireStateRepository(filepath.Join(dir, "savedata", "grimoire.json"))
	skillRepo := store.NewEntryStateRepository[skills.SkillEntry](filepath.Join(dir, "savedata", "skill.json"))
	enemySkillRepo := store.NewEntryStateRepository[enemyskills.EnemySkillEntry](filepath.Join(dir, "savedata", "enemy_skill.json"))
	enemyRepo := store.NewEntryStateRepository[enemies.EnemyEntry](filepath.Join(dir, "savedata", "enemy.json"))
	spawnRepo := store.NewEntryStateRepository[spawntables.SpawnTableEntry](filepath.Join(dir, "savedata", "spawn_table.json"))
	treasureRepo := store.NewEntryStateRepository[treasures.TreasureEntry](filepath.Join(dir, "savedata", "treasure.json"))
	loottableRepo := store.NewEntryStateRepository[loottables.LootTableEntry](filepath.Join(dir, "savedata", "loottables.json"))

	if err := itemRepo.SaveItemState(items.ItemState{Items: []items.ItemEntry{}}); err != nil {
		t.Fatalf("seed item: %v", err)
	}
	if err := grimoireRepo.SaveGrimoireState(grimoire.GrimoireState{Entries: []grimoire.GrimoireEntry{}}); err != nil {
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
