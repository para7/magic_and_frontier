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
	GrimoireEffect     string `json:"grimoireEffect"`
	GrimoireDebug      string `json:"grimoireDebug"`
	PassiveEffect      string `json:"passiveEffect"`
	PassiveGive        string `json:"passiveGive"`
	PassiveApply       string `json:"passiveApply"`
	Enemy              string `json:"enemy"`
	EnemySkill         string `json:"enemySkill"`
	EnemyLoot          string `json:"enemyLoot"`
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
	PassiveStatePath       string
	EnemySkillStatePath    string
	EnemyStatePath         string
	SpawnTableStatePath    string
	TreasureStatePath      string
	LootTablesStatePath    string
	ExportSettingsPath     string
	MinecraftLootTableRoot string
}

func LoadConfig() MafConfig {
	return MafConfig{
		Port:                   3000,
		ItemStatePath:          filepath.Clean(filepath.Join("savedata", "item.json")),
		GrimoireStatePath:      filepath.Clean(filepath.Join("savedata", "grimoire.json")),
		PassiveStatePath:       filepath.Clean(filepath.Join("savedata", "passive.json")),
		EnemySkillStatePath:    filepath.Clean(filepath.Join("savedata", "enemy_skill.json")),
		EnemyStatePath:         filepath.Clean(filepath.Join("savedata", "enemy.json")),
		SpawnTableStatePath:    filepath.Clean(filepath.Join("savedata", "spawn_table.json")),
		TreasureStatePath:      filepath.Clean(filepath.Join("savedata", "treasure.json")),
		LootTablesStatePath:    filepath.Clean(filepath.Join("savedata", "loottables.json")),
		ExportSettingsPath:     filepath.Clean(filepath.Join("config", "export_settings.json")),
		MinecraftLootTableRoot: filepath.Clean(filepath.Join("minecraft", "1.21.11", "loot_table")),
	}
}
