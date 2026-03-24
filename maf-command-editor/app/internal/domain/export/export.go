package export

import (
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
)

type ExportSettings struct {
	OutputRoot       string      `json:"outputRoot"`
	Namespace        string      `json:"namespace"`
	TemplatePackPath string      `json:"templatePackPath"`
	Paths            ExportPaths `json:"paths"`
}

type ExportPaths struct {
	ItemFunctionDir       string `json:"itemFunctionDir"`
	ItemLootDir           string `json:"itemLootDir"`
	SpellFunctionDir      string `json:"spellFunctionDir"`
	SpellLootDir          string `json:"spellLootDir"`
	SkillFunctionDir      string `json:"skillFunctionDir"`
	EnemySkillFunctionDir string `json:"enemySkillFunctionDir"`
	EnemyFunctionDir      string `json:"enemyFunctionDir"`
	EnemyLootDir          string `json:"enemyLootDir"`
	TreasureLootDir       string `json:"treasureLootDir"`
	LoottableLootDir      string `json:"loottableLootDir"`
	DebugFunctionDir      string `json:"debugFunctionDir"`
	MinecraftTagDir       string `json:"minecraftTagDir"`
}

type ExportStats struct {
	ItemFunctions       int `json:"itemFunctions"`
	ItemLootTables      int `json:"itemLootTables"`
	SpellFunctions      int `json:"spellFunctions"`
	SpellLootTables     int `json:"spellLootTables"`
	SkillFunctions      int `json:"skillFunctions"`
	EnemySkillFunctions int `json:"enemySkillFunctions"`
	EnemyFunctions      int `json:"enemyFunctions"`
	EnemyLootTables     int `json:"enemyLootTables"`
	TreasureLootTables  int `json:"treasureLootTables"`
	LoottableLootTables int `json:"loottableLootTables"`
	TotalFiles          int `json:"totalFiles"`
}

type SaveDataResponse struct {
	OK         bool         `json:"ok"`
	Message    string       `json:"message,omitempty"`
	OutputRoot string       `json:"outputRoot,omitempty"`
	Generated  *ExportStats `json:"generated,omitempty"`
	Code       string       `json:"code,omitempty"`
	Details    string       `json:"details,omitempty"`
}

type ExportParams struct {
	ItemState              items.ItemState
	GrimoireState          grimoire.GrimoireState
	Skills                 []skills.SkillEntry
	EnemySkills            []enemyskills.EnemySkillEntry
	Enemies                []enemies.EnemyEntry
	SpawnTables            []spawntables.SpawnTableEntry
	Treasures              []treasures.TreasureEntry
	LootTables             []loottables.LootTableEntry
	ExportSettingsPath     string
	MinecraftLootTableRoot string
}

func ExportDatapack(params ExportParams) SaveDataResponse {
	settings, err := loadExportSettings(params.ExportSettingsPath)
	if err != nil {
		return exportFailure("INVALID_CONFIG", "Invalid export settings.", err)
	}
	if err := writeDatapackScaffold(settings); err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}

	itemStats, err := generateItemOutputs(settings, params.ItemState.Items)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	grimoireStats, err := generateGrimoireOutputs(settings, params.GrimoireState.Entries)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	debugGrimoireFunctions, err := generateGrimoireDebugFunctions(settings, params.GrimoireState.Entries)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	skillStats, err := generateSkillOutputs(settings, params.Skills)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	enemySkillStats, err := generateEnemySkillOutputs(settings, params.EnemySkills)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	enemyStats, err := generateEnemyOutputs(settings, params.Enemies, params.ItemState.Items, params.GrimoireState.Entries)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	spawnTableStats, err := generateSpawnTableOutputs(settings, params.SpawnTables)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	treasureStats, err := generateTreasureOutputs(settings, params.MinecraftLootTableRoot, params.Treasures, params.ItemState.Items, params.GrimoireState.Entries)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	loottableStats, err := generateLootTableOutputs(settings, params.LootTables, params.ItemState.Items, params.GrimoireState.Entries)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	tickDispatcherFiles, err := generateTickDispatcher(settings)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}

	stats := &ExportStats{
		ItemFunctions:       itemStats.ItemFunctions,
		ItemLootTables:      itemStats.ItemLootTables,
		SpellFunctions:      grimoireStats.SpellFunctions,
		SpellLootTables:     grimoireStats.SpellLootTables,
		SkillFunctions:      skillStats.SkillFunctions,
		EnemySkillFunctions: enemySkillStats.EnemySkillFunctions,
		EnemyFunctions:      enemyStats.EnemyFunctions,
		EnemyLootTables:     enemyStats.EnemyLootTables,
		TreasureLootTables:  treasureStats.TreasureLootTables,
		LoottableLootTables: loottableStats.LoottableLootTables,
	}
	stats.TotalFiles = stats.ItemFunctions + stats.ItemLootTables + stats.SpellFunctions + stats.SpellLootTables + stats.SkillFunctions + stats.EnemySkillFunctions + stats.EnemyFunctions + stats.EnemyLootTables + stats.TreasureLootTables + stats.LoottableLootTables + debugGrimoireFunctions + tickDispatcherFiles + spawnTableStats.SpawnTableFunctions

	return SaveDataResponse{
		OK:         true,
		Message:    "datapack export completed",
		OutputRoot: settings.OutputRoot,
		Generated:  stats,
	}
}

func ValidateSettings(settingsPath string) error {
	_, err := loadExportSettings(settingsPath)
	return err
}

type itemOutputStats struct {
	ItemFunctions  int
	ItemLootTables int
}

type spellOutputStats struct {
	SpellFunctions  int
	SpellLootTables int
}

type skillOutputStats struct {
	SkillFunctions int
}

type enemySkillOutputStats struct {
	EnemySkillFunctions int
}

type enemyOutputStats struct {
	EnemyFunctions  int
	EnemyLootTables int
}

type spawnTableOutputStats struct {
	SpawnTableFunctions int
}

type treasureOutputStats struct {
	TreasureLootTables int
}

type loottableOutputStats struct {
	LoottableLootTables int
}

func exportFailure(code, message string, err error) SaveDataResponse {
	return SaveDataResponse{
		OK:      false,
		Code:    code,
		Message: message,
		Details: err.Error(),
	}
}
