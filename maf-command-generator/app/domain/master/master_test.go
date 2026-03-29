package master

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	"maf_command_editor/app/files"
)

type entriesFile[T any] struct {
	Entries []T `json:"entries"`
}

func writeState[T any](t *testing.T, path string, entries []T) {
	t.Helper()
	data, err := json.MarshalIndent(entriesFile[T]{Entries: entries}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

// 全エンティティのパスを一時ディレクトリに向けた設定を返す。
// grimoires には任意のグリモアデータを渡せる。
func newTestConfig(t *testing.T, grimoires []grimoireModel.Grimoire) files.MafConfig {
	t.Helper()
	dir := t.TempDir()
	p := func(name string) string { return filepath.Join(dir, name) }

	writeState(t, p("grimoire.json"), grimoires)
	writeState[struct{}](t, p("item.json"), nil)
	writeState[struct{}](t, p("skill.json"), nil)
	writeState[struct{}](t, p("enemy_skill.json"), nil)
	writeState[struct{}](t, p("enemy.json"), nil)
	writeState[struct{}](t, p("spawn_table.json"), nil)
	writeState[struct{}](t, p("treasure.json"), nil)
	writeState[struct{}](t, p("loottables.json"), nil)

	cfg := files.LoadConfig()
	cfg.GrimoireStatePath = p("grimoire.json")
	cfg.ItemStatePath = p("item.json")
	cfg.SkillStatePath = p("skill.json")
	cfg.EnemySkillStatePath = p("enemy_skill.json")
	cfg.EnemyStatePath = p("enemy.json")
	cfg.SpawnTableStatePath = p("spawn_table.json")
	cfg.TreasureStatePath = p("treasure.json")
	cfg.LootTablesStatePath = p("loottables.json")
	return cfg
}

func TestDBMasterImplExportMethods(t *testing.T) {
	grimoires := []grimoireModel.Grimoire{
		{ID: "g1", CastID: 1, Script: "say 1", Title: "one"},
		{ID: "g2", CastID: 2, Script: "say 2", Title: "two"},
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
		{ID: "g1", CastID: 1, Script: "say 1", Title: "one"},
	}
	db := NewDBMaster(newTestConfig(t, grimoires))

	list := db.ListGrimoires()
	list[0].Title = "mutated"

	again := db.ListGrimoires()
	if again[0].Title != "one" {
		t.Fatalf("internal state must not be affected by caller mutation, got %q", again[0].Title)
	}
}
