package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	enemyModel "maf_command_editor/app/domain/model/enemy"
	enemyskillModel "maf_command_editor/app/domain/model/enemyskill"
	config "maf_command_editor/app/files"
)

func TestExportDatapackDoesNotWriteLegacyEnemyCompatibilityWrappers(t *testing.T) {
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	settings := map[string]any{
		"outputRoot": filepath.Join(root, "out"),
		"exportPaths": map[string]any{
			"grimoireEffect":     "generated/grimoire/effect",
			"grimoireSelectFile": "generated/grimoire/selectexec.mcfunction",
			"grimoireDebug":      "generated/grimoire/give",
			"enemy":              "generated/enemy/spawn",
			"enemySkill":         "generated/enemy/skill",
			"enemyLoot":          "generated/enemy/loot",
		},
	}
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := config.LoadConfig()
	cfg.ExportSettingsPath = settingsPath
	cfg.MinecraftLootTableRoot = filepath.Join(root, "minecraft", "loot_table")

	master := exportMasterStub{
		enemySkills: []enemyskillModel.EnemySkill{
			{ID: "near_poison", Script: "say poison"},
		},
		enemies: []enemyModel.Enemy{
			{ID: "poison_zombie", MobType: "minecraft:zombie", HP: 20, DropMode: "replace"},
		},
	}

	if err := ExportDatapack(master, cfg); err != nil {
		t.Fatal(err)
	}

	newSkillMain := filepath.Join(root, "out", "data", "maf", "function", "generated", "enemy", "skill", "main.mcfunction")
	if _, err := os.Stat(newSkillMain); err != nil {
		t.Fatalf("missing enemy skill main: %v", err)
	}
	newEnemy := filepath.Join(root, "out", "data", "maf", "function", "generated", "enemy", "spawn", "poison_zombie.mcfunction")
	if _, err := os.Stat(newEnemy); err != nil {
		t.Fatalf("missing enemy spawn function: %v", err)
	}

	legacySkillMain := filepath.Join(root, "out", "data", "maf", "function", "generated", "enemy_skill", "main.mcfunction")
	if _, err := os.Stat(legacySkillMain); !os.IsNotExist(err) {
		t.Fatalf("legacy enemy skill main should not be created: %v", err)
	}

	legacyEnemy := filepath.Join(root, "out", "data", "maf", "function", "generated", "enemy", "poison_zombie.mcfunction")
	if _, err := os.Stat(legacyEnemy); !os.IsNotExist(err) {
		t.Fatalf("legacy enemy wrapper should not be created: %v", err)
	}
}
