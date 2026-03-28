package files

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ExportSettings struct {
	OutputRoot  string      `json:"outputRoot"`
	ExportPaths ExportPaths `json:"exportPaths"`
}

type ExportPaths struct {
	SpellEffect string `json:"spellEffect"`
}

func LoadExportSettings(path string) (ExportSettings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ExportSettings{}, err
	}
	var s ExportSettings
	return s, json.Unmarshal(data, &s)
}

type MafConfig struct {
	Port                   int
	ItemStatePath          string
	GrimoireStatePath      string
	SkillStatePath         string
	EnemySkillStatePath    string
	EnemyStatePath         string
	SpawnTableStatePath    string
	TreasureStatePath      string
	LootTablesStatePath    string
	IDCounterStatePath     string
	ExportSettingsPath     string
	MinecraftLootTableRoot string
}

func LoadConfig() MafConfig {
	return MafConfig{
		Port:                   3000,
		ItemStatePath:          filepath.Clean(filepath.Join("savedata", "item.json")),
		GrimoireStatePath:      filepath.Clean(filepath.Join("savedata", "grimoire.json")),
		SkillStatePath:         filepath.Clean(filepath.Join("savedata", "skill.json")),
		EnemySkillStatePath:    filepath.Clean(filepath.Join("savedata", "enemy_skill.json")),
		EnemyStatePath:         filepath.Clean(filepath.Join("savedata", "enemy.json")),
		SpawnTableStatePath:    filepath.Clean(filepath.Join("savedata", "spawn_table.json")),
		TreasureStatePath:      filepath.Clean(filepath.Join("savedata", "treasure.json")),
		LootTablesStatePath:    filepath.Clean(filepath.Join("savedata", "loottables.json")),
		IDCounterStatePath:     filepath.Clean(filepath.Join("savedata", "id_counters.json")),
		ExportSettingsPath:     filepath.Clean(filepath.Join("config", "export_settings.json")),
		MinecraftLootTableRoot: filepath.Clean(filepath.Join("minecraft", "1.21.11", "loot_table")),
	}
}
