package master

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	bowModel "maf_command_editor/app/domain/model/bow"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	passiveModel "maf_command_editor/app/domain/model/passive"
	"maf_command_editor/app/files"
)

type entriesFile[T any] struct {
	Entries []T `json:"entries"`
}

// writeState creates a directory at dirPath and writes entries to dirPath/entity.json.
func writeState[T any](t *testing.T, dirPath string, entries []T) {
	t.Helper()
	data, err := json.MarshalIndent(entriesFile[T]{Entries: entries}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dirPath, "entity.json"), append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

// 全エンティティのパスを一時ディレクトリに向けた設定を返す。
// grimoires には任意のグリモアデータを渡せる。
func newTestConfig(t *testing.T, grimoires []grimoireModel.Grimoire) files.MafConfig {
	t.Helper()
	dir := t.TempDir()
	p := func(name string) string { return filepath.Join(dir, name) }

	writeState(t, p("grimoire"), grimoires)
	writeState[struct{}](t, p("item"), nil)
	writeState[struct{}](t, p("passive"), nil)
	writeState[struct{}](t, p("bow"), nil)
	writeState[struct{}](t, p("enemy_skill"), nil)
	writeState[struct{}](t, p("enemy"), nil)
	writeState[struct{}](t, p("spawn_table"), nil)
	writeState[struct{}](t, p("treasure"), nil)
	writeState[struct{}](t, p("loottables"), nil)

	cfg := files.LoadConfig()
	cfg.GrimoireStatePath = p("grimoire")
	cfg.ItemStatePath = p("item")
	cfg.PassiveStatePath = p("passive")
	cfg.BowStatePath = p("bow")
	cfg.EnemySkillStatePath = p("enemy_skill")
	cfg.EnemyStatePath = p("enemy")
	cfg.SpawnTableStatePath = p("spawn_table")
	cfg.TreasureStatePath = p("treasure")
	cfg.LootTablesStatePath = p("loottables")
	return cfg
}

func TestDBMasterImplExportMethods(t *testing.T) {
	grimoires := []grimoireModel.Grimoire{
		{ID: "g1", Script: []string{"say 1"}, Title: "one"},
		{ID: "g2", Script: []string{"say 2"}, Title: "two"},
	}
	db := NewDBMaster(newTestConfig(t, grimoires))

	list := db.ListGrimoires()
	if len(list) != 2 {
		t.Fatalf("ListGrimoires length = %d, want 2", len(list))
	}
	if list[0].ID != "g1" || list[1].ID != "g2" {
		t.Fatalf("unexpected list order/content: %#v", list)
	}

	found, ok := db.GetGrimoireByID("g2")
	if !ok || found.ID != "g2" {
		t.Fatalf("GetGrimoireByID(g2) = (%#v, %v), want found", found, ok)
	}
	_, ok = db.GetGrimoireByID("missing")
	if ok {
		t.Fatalf("GetGrimoireByID(missing) should be not found")
	}
}

func TestDBMasterImplListGrimoiresReturnsCopy(t *testing.T) {
	grimoires := []grimoireModel.Grimoire{
		{ID: "g1", Script: []string{"say 1"}, Title: "one"},
	}
	db := NewDBMaster(newTestConfig(t, grimoires))

	list := db.ListGrimoires()
	list[0].Title = "mutated"

	again := db.ListGrimoires()
	if again[0].Title != "one" {
		t.Fatalf("internal state must not be affected by caller mutation, got %q", again[0].Title)
	}
}

func TestDBMasterValidateAllRejectsBowEffectIDCollision(t *testing.T) {
	cfg := newTestConfig(t, nil)
	writeState(t, cfg.PassiveStatePath, []passiveModel.Passive{
		{ID: "bow_test_full", Condition: "always", Slots: []int{1}, Script: []string{"say passive"}},
	})
	writeState(t, cfg.BowStatePath, []bowModel.BowPassive{
		{ID: "test_full"},
	})

	db := NewDBMaster(cfg)
	errs := db.ValidateAll()
	for _, recordErrs := range errs {
		for _, err := range recordErrs {
			if err.Entity == "bow" && err.ID == "test_full" && err.Field == "id" {
				return
			}
		}
	}
	t.Fatalf("expected bow/passive effect id collision error, got %#v", errs)
}
