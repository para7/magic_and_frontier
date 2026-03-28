package master

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	"maf_command_editor/app/files"
)

type grimoireEntriesFile struct {
	Entries []grimoireModel.Grimoire `json:"entries"`
}

func writeGrimoireState(t *testing.T, path string, entries []grimoireModel.Grimoire) {
	t.Helper()
	data, err := json.MarshalIndent(grimoireEntriesFile{Entries: entries}, "", "  ")
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

func TestDBMasterImplExportMethods(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "grimoire.json")
	writeGrimoireState(t, statePath, []grimoireModel.Grimoire{
		{ID: "g1", CastID: 1, Script: "say 1", Title: "one"},
		{ID: "g2", CastID: 2, Script: "say 2", Title: "two"},
	})

	cfg := files.LoadConfig()
	cfg.GrimoireStatePath = statePath
	db := NewDBMaster(cfg)

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
	statePath := filepath.Join(t.TempDir(), "grimoire.json")
	writeGrimoireState(t, statePath, []grimoireModel.Grimoire{
		{ID: "g1", CastID: 1, Script: "say 1", Title: "one"},
	})

	cfg := files.LoadConfig()
	cfg.GrimoireStatePath = statePath
	db := NewDBMaster(cfg)

	list := db.ListGrimoires()
	list[0].Title = "mutated"

	again := db.ListGrimoires()
	if again[0].Title != "one" {
		t.Fatalf("internal state must not be affected by caller mutation, got %q", again[0].Title)
	}
}
