package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type rawExportSettings struct {
	OutputRoot       *string         `json:"outputRoot"`
	Namespace        *string         `json:"namespace"`
	TemplatePackPath *string         `json:"templatePackPath"`
	Paths            *rawExportPaths `json:"paths"`
}

type rawExportPaths struct {
	ItemFunctionDir       *string `json:"itemFunctionDir"`
	ItemLootDir           *string `json:"itemLootDir"`
	SpellFunctionDir      *string `json:"spellFunctionDir"`
	SpellLootDir          *string `json:"spellLootDir"`
	SkillFunctionDir      *string `json:"skillFunctionDir"`
	EnemySkillFunctionDir *string `json:"enemySkillFunctionDir"`
	EnemyFunctionDir      *string `json:"enemyFunctionDir"`
	EnemyLootDir          *string `json:"enemyLootDir"`
	TreasureLootDir       *string `json:"treasureLootDir"`
	LoottableLootDir      *string `json:"loottableLootDir"`
	DebugFunctionDir      *string `json:"debugFunctionDir"`
	MinecraftTagDir       *string `json:"minecraftTagDir"`
}

func loadExportSettings(settingsPath string) (ExportSettings, error) {
	var parsed rawExportSettings
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return ExportSettings{}, fmt.Errorf("failed to read export settings at %s: %w", settingsPath, err)
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return ExportSettings{}, fmt.Errorf("failed to read export settings at %s: %w", settingsPath, err)
	}

	outputRoot, err := requireString(parsed.OutputRoot, "outputRoot")
	if err != nil {
		return ExportSettings{}, err
	}
	namespace, err := requireString(parsed.Namespace, "namespace")
	if err != nil {
		return ExportSettings{}, err
	}
	if !namespacePattern(namespace) {
		return ExportSettings{}, fmt.Errorf("namespace must match [a-z0-9_.-]+")
	}
	templatePackPath, err := requireString(parsed.TemplatePackPath, "templatePackPath")
	if err != nil {
		return ExportSettings{}, err
	}

	if parsed.Paths == nil {
		return ExportSettings{}, fmt.Errorf("paths must be an object")
	}
	pathValues := parsed.Paths

	baseDir := filepath.Dir(settingsPath)
	itemFunctionDir, err := requirePathString(pathValues.ItemFunctionDir, "itemFunctionDir")
	if err != nil {
		return ExportSettings{}, err
	}
	itemLootDir, err := requirePathString(pathValues.ItemLootDir, "itemLootDir")
	if err != nil {
		return ExportSettings{}, err
	}
	spellFunctionDir, err := requirePathString(pathValues.SpellFunctionDir, "spellFunctionDir")
	if err != nil {
		return ExportSettings{}, err
	}
	spellLootDir, err := requirePathString(pathValues.SpellLootDir, "spellLootDir")
	if err != nil {
		return ExportSettings{}, err
	}
	skillFunctionDir, err := pathStringOrDefault(pathValues.SkillFunctionDir, "skillFunctionDir", defaultGeneratedPath(namespace, "function", "skill"))
	if err != nil {
		return ExportSettings{}, err
	}
	enemySkillFunctionDir, err := pathStringOrDefault(pathValues.EnemySkillFunctionDir, "enemySkillFunctionDir", defaultGeneratedPath(namespace, "function", "enemy_skill"))
	if err != nil {
		return ExportSettings{}, err
	}
	enemyFunctionDir, err := pathStringOrDefault(pathValues.EnemyFunctionDir, "enemyFunctionDir", defaultGeneratedPath(namespace, "function", "enemy"))
	if err != nil {
		return ExportSettings{}, err
	}
	enemyLootDir, err := pathStringOrDefault(pathValues.EnemyLootDir, "enemyLootDir", defaultGeneratedPath(namespace, "loot_table", "enemy"))
	if err != nil {
		return ExportSettings{}, err
	}
	treasureLootDir, err := pathStringOrDefault(pathValues.TreasureLootDir, "treasureLootDir", defaultGeneratedPath(namespace, "loot_table", "treasure"))
	if err != nil {
		return ExportSettings{}, err
	}
	loottableLootDir, err := pathStringOrDefault(pathValues.LoottableLootDir, "loottableLootDir", defaultGeneratedPath(namespace, "loot_table", "loottable"))
	if err != nil {
		return ExportSettings{}, err
	}
	debugFunctionDir, err := pathStringOrDefault(pathValues.DebugFunctionDir, "debugFunctionDir", filepath.ToSlash(filepath.Join("data", namespace, "function", "debug", "give")))
	if err != nil {
		return ExportSettings{}, err
	}
	minecraftTagDir, err := requirePathString(pathValues.MinecraftTagDir, "minecraftTagDir")
	if err != nil {
		return ExportSettings{}, err
	}

	return ExportSettings{
		OutputRoot:       filepath.Clean(filepath.Join(baseDir, outputRoot)),
		Namespace:        namespace,
		TemplatePackPath: filepath.Clean(filepath.Join(baseDir, templatePackPath)),
		Paths: ExportPaths{
			ItemFunctionDir:       itemFunctionDir,
			ItemLootDir:           itemLootDir,
			SpellFunctionDir:      spellFunctionDir,
			SpellLootDir:          spellLootDir,
			SkillFunctionDir:      skillFunctionDir,
			EnemySkillFunctionDir: enemySkillFunctionDir,
			EnemyFunctionDir:      enemyFunctionDir,
			EnemyLootDir:          enemyLootDir,
			TreasureLootDir:       treasureLootDir,
			LoottableLootDir:      loottableLootDir,
			DebugFunctionDir:      debugFunctionDir,
			MinecraftTagDir:       minecraftTagDir,
		},
	}, nil
}

func defaultGeneratedPath(namespace, registry string, parts ...string) string {
	segments := append([]string{"data", namespace, registry, "generated"}, parts...)
	return filepath.ToSlash(filepath.Join(segments...))
}

func requireString(value *string, key string) (string, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	return *value, nil
}

func requirePathString(value *string, key string) (string, error) {
	return requireString(value, "paths."+key)
}

func pathStringOrDefault(value *string, key, fallback string) (string, error) {
	if value == nil {
		return fallback, nil
	}
	return requireString(value, "paths."+key)
}

func namespacePattern(value string) bool {
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '.' || r == '-' {
			continue
		}
		return false
	}
	return value != ""
}
