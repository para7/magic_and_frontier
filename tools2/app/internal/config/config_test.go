package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultStatePath_ResolvesSiblingToolsSavedata(t *testing.T) {
	root := t.TempDir()
	tools2Root := filepath.Join(root, "tools2")
	toolsRoot := filepath.Join(root, "tools")
	if err := os.MkdirAll(filepath.Join(tools2Root), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(toolsRoot, "savedata"), 0o755); err != nil {
		t.Fatal(err)
	}
	statePath := filepath.Join(toolsRoot, "savedata", "form-state.json")
	if err := os.WriteFile(statePath, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Chdir(tools2Root)

	got := defaultStatePath("form-state.json")
	want := filepath.Clean(filepath.Join("..", "tools", "savedata", "form-state.json"))
	if got != want {
		t.Fatalf("defaultStatePath = %q, want %q", got, want)
	}
}

func TestDefaultStatePath_PrefersExistingLocalSavedata(t *testing.T) {
	root := t.TempDir()
	tools2Root := filepath.Join(root, "tools2")
	localSavedata := filepath.Join(tools2Root, "savedata")
	siblingSavedata := filepath.Join(root, "tools", "savedata")
	if err := os.MkdirAll(localSavedata, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(siblingSavedata, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(localSavedata, "form-state.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(siblingSavedata, "form-state.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Chdir(tools2Root)

	got := defaultStatePath("form-state.json")
	want := filepath.Clean(filepath.Join(".", "savedata", "form-state.json"))
	if got != want {
		t.Fatalf("defaultStatePath = %q, want %q", got, want)
	}
}

func TestDefaultExportSettingsPath_ResolvesSiblingToolsServerConfig(t *testing.T) {
	root := t.TempDir()
	tools2Root := filepath.Join(root, "tools2")
	toolsConfigDir := filepath.Join(root, "tools", "server", "config")
	if err := os.MkdirAll(toolsConfigDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(tools2Root, 0o755); err != nil {
		t.Fatal(err)
	}
	settingsPath := filepath.Join(toolsConfigDir, "export-settings.json")
	if err := os.WriteFile(settingsPath, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Chdir(tools2Root)

	got := defaultExportSettingsPath()
	want := filepath.Clean(filepath.Join("..", "tools", "server", "config", "export-settings.json"))
	if got != want {
		t.Fatalf("defaultExportSettingsPath = %q, want %q", got, want)
	}
}
