package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Port                int
	ItemStatePath       string
	GrimoireStatePath   string
	SkillStatePath      string
	EnemySkillStatePath string
	EnemyStatePath      string
	TreasureStatePath   string
	ExportSettingsPath  string
}

func Load() Config {
	rawPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || rawPort <= 0 {
		rawPort = 8787
	}

	return Config{
		Port:                rawPort,
		ItemStatePath:       envOrDefault("ITEM_STATE_PATH", defaultStatePath("form-state.json")),
		GrimoireStatePath:   envOrDefault("GRIMOIRE_STATE_PATH", defaultStatePath("grimoire-state.json")),
		SkillStatePath:      envOrDefault("SKILL_STATE_PATH", defaultStatePath("skill-state.json")),
		EnemySkillStatePath: envOrDefault("ENEMY_SKILL_STATE_PATH", defaultStatePath("enemy-skill-state.json")),
		EnemyStatePath:      envOrDefault("ENEMY_STATE_PATH", defaultStatePath("enemy-state.json")),
		TreasureStatePath:   envOrDefault("TREASURE_STATE_PATH", defaultStatePath("treasure-state.json")),
		ExportSettingsPath:  envOrDefault("EXPORT_SETTINGS_PATH", defaultExportSettingsPath()),
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func defaultStatePath(fileName string) string {
	candidates := []string{
		filepath.Clean(filepath.Join(".", "savedata", fileName)),
		filepath.Clean(filepath.Join("..", "savedata", fileName)),
	}
	return firstExistingOrDefault(candidates)
}

func defaultExportSettingsPath() string {
	candidates := []string{
		filepath.Clean(filepath.Join(".", "config", "export-settings.json")),
		filepath.Clean(filepath.Join(".", "server", "config", "export-settings.json")),
		filepath.Clean(filepath.Join("..", "tools", "server", "config", "export-settings.json")),
		filepath.Clean(filepath.Join(".", "tools", "server", "config", "export-settings.json")),
	}
	return firstExistingOrDefault(candidates)
}

func firstExistingOrDefault(candidates []string) string {
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return candidates[0]
}
