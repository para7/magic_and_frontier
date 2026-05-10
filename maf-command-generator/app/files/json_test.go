package files

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type testEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func writeJSON(t *testing.T, path string, entries []testEntry) {
	t.Helper()
	data, err := json.MarshalIndent(entriesFile[testEntry]{Entries: entries}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestJsonStore_Load_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, filepath.Join(dir, "swords.json"), []testEntry{{ID: "s1", Name: "Iron Sword"}, {ID: "s2", Name: "Gold Sword"}})
	writeJSON(t, filepath.Join(dir, "potions.json"), []testEntry{{ID: "p1", Name: "Potion"}})

	s := NewJsonStore[testEntry](dir)
	entries, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("want 3 entries, got %d", len(entries))
	}
}

func TestJsonStore_Load_NestedFiles(t *testing.T) {
	dir := t.TempDir()
	nestedDir := filepath.Join(dir, "level", "one")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeJSON(t, filepath.Join(dir, "root.json"), []testEntry{{ID: "root", Name: "Root"}})
	writeJSON(t, filepath.Join(nestedDir, "nested.json"), []testEntry{{ID: "nested", Name: "Nested"}})

	s := NewJsonStore[testEntry](dir)
	entries, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}
	seen := map[string]bool{}
	for _, entry := range entries {
		seen[entry.ID] = true
	}
	if !seen["root"] || !seen["nested"] {
		t.Fatalf("expected root and nested entries, got %#v", entries)
	}
}

func TestJsonStore_Load_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	s := NewJsonStore[testEntry](dir)
	entries, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("want 0 entries, got %d", len(entries))
	}
}
