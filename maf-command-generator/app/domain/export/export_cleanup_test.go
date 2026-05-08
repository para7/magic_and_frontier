package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	activeModel "maf_command_editor/app/domain/model/active"
	config "maf_command_editor/app/files"
)

func TestExportDatapackCleansConfiguredPaths(t *testing.T) {
	root := t.TempDir()
	outputRoot := filepath.Join(root, "out")

	settings := defaultFixtureExportSettings(outputRoot)
	settings.CleanPaths = []string{
		filepath.Join(outputRoot, "data", "maf", "function", "generated"),
		filepath.Join(outputRoot, "data", "maf", "loot_table", "generated"),
	}
	settingsPath := filepath.Join(root, "export_settings.json")
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	staleFunctionPath := filepath.Join(outputRoot, "data", "maf", "function", "generated", "active", "effect", "fulminant_true_01.mcfunction")
	staleLootPath := filepath.Join(outputRoot, "data", "maf", "loot_table", "generated", "enemy", "loot", "stale_enemy.json")
	if err := writeFunctionFile(staleFunctionPath, "say stale"); err != nil {
		t.Fatal(err)
	}
	if err := writeJSON(staleLootPath, map[string]any{"type": "minecraft:generic"}); err != nil {
		t.Fatal(err)
	}

	cfg := config.LoadConfig()
	cfg.ExportSettingsPath = settingsPath
	cfg.LootTableSourceRoot = filepath.Join(root, "loot_table")
	cfg.MinecraftLootTableRoot = filepath.Join(root, "minecraft", "loot_table")

	master := exportMasterStub{
		actives: []activeModel.Active{
			{
				ID:          "fulminant01",
				Description: "雷鳴で敵を攻撃する",
				MPCost:      20,
				CastTime:    60,
				CoolTime:    20,
				Script:      []string{"say current"},
			},
		},
	}

	if err := ExportDatapack(master, cfg); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(staleFunctionPath); !os.IsNotExist(err) {
		t.Fatalf("stale function should be removed, err=%v", err)
	}
	if _, err := os.Stat(staleLootPath); !os.IsNotExist(err) {
		t.Fatalf("stale loot should be removed, err=%v", err)
	}

	currentPath := filepath.Join(outputRoot, "data", "maf", "function", "generated", "active", "effect", "fulminant01.mcfunction")
	if _, err := os.Stat(currentPath); err != nil {
		t.Fatalf("current active should be generated: %v", err)
	}
}

func TestExportDatapackFailsWhenCleanPathOutsideOutputRoot(t *testing.T) {
	root := t.TempDir()
	outputRoot := filepath.Join(root, "out")

	settings := defaultFixtureExportSettings(outputRoot)
	settings.CleanPaths = []string{
		filepath.Join(root, "other", "outside"),
	}
	settingsPath := filepath.Join(root, "export_settings.json")
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := config.LoadConfig()
	cfg.ExportSettingsPath = settingsPath
	cfg.LootTableSourceRoot = filepath.Join(root, "loot_table")
	cfg.MinecraftLootTableRoot = filepath.Join(root, "minecraft", "loot_table")

	err = ExportDatapack(exportMasterStub{}, cfg)
	if err == nil {
		t.Fatal("expected error when clean path is outside outputRoot")
	}
	if !strings.Contains(err.Error(), "must be inside outputRoot") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExportDatapackFailsWhenCleanPathContainsUnsupportedExtension(t *testing.T) {
	root := t.TempDir()
	outputRoot := filepath.Join(root, "out")
	cleanPath := filepath.Join(outputRoot, "data", "maf", "function", "generated")

	if err := os.MkdirAll(cleanPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cleanPath, "unsafe.txt"), []byte("unsafe"), 0o644); err != nil {
		t.Fatal(err)
	}

	settings := defaultFixtureExportSettings(outputRoot)
	settings.CleanPaths = []string{cleanPath}
	settingsPath := filepath.Join(root, "export_settings.json")
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := config.LoadConfig()
	cfg.ExportSettingsPath = settingsPath
	cfg.LootTableSourceRoot = filepath.Join(root, "loot_table")
	cfg.MinecraftLootTableRoot = filepath.Join(root, "minecraft", "loot_table")

	err = ExportDatapack(exportMasterStub{}, cfg)
	if err == nil {
		t.Fatal("expected error when clean path has unsupported extension")
	}
	if !strings.Contains(err.Error(), "unsupported extension") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExportDatapackSkipsCleanupWhenCleanPathsEmpty(t *testing.T) {
	root := t.TempDir()
	outputRoot := filepath.Join(root, "out")

	settings := defaultFixtureExportSettings(outputRoot)
	settingsPath := filepath.Join(root, "export_settings.json")
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	staleFunctionPath := filepath.Join(outputRoot, "data", "maf", "function", "generated", "active", "effect", "fulminant_true_01.mcfunction")
	if err := writeFunctionFile(staleFunctionPath, "say stale"); err != nil {
		t.Fatal(err)
	}

	cfg := config.LoadConfig()
	cfg.ExportSettingsPath = settingsPath
	cfg.LootTableSourceRoot = filepath.Join(root, "loot_table")
	cfg.MinecraftLootTableRoot = filepath.Join(root, "minecraft", "loot_table")

	master := exportMasterStub{
		actives: []activeModel.Active{
			{
				ID:          "fulminant01",
				Description: "雷鳴で敵を攻撃する",
				MPCost:      20,
				CastTime:    60,
				CoolTime:    20,
				Script:      []string{"say current"},
			},
		},
	}

	if err := ExportDatapack(master, cfg); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(staleFunctionPath); err != nil {
		t.Fatalf("stale function should remain when cleanPaths is empty: %v", err)
	}
}
