package export

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/mcsource"
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

type treasureOverlayManifest struct {
	Paths []string `json:"paths"`
}

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

func writeDatapackScaffold(settings ExportSettings) error {
	if err := os.MkdirAll(settings.OutputRoot, 0o755); err != nil {
		return err
	}

	outputDirs := []string{
		settings.Paths.ItemFunctionDir,
		settings.Paths.ItemLootDir,
		settings.Paths.SpellFunctionDir,
		settings.Paths.SpellLootDir,
		settings.Paths.SkillFunctionDir,
		settings.Paths.EnemySkillFunctionDir,
		settings.Paths.EnemyFunctionDir,
		settings.Paths.EnemyLootDir,
		settings.Paths.TreasureLootDir,
		settings.Paths.LoottableLootDir,
		generatedTickFunctionDir(settings),
		generatedGrimoireDebugFunctionDir(settings),
	}
	cleanupDirs := make([]string, 0, len(outputDirs))
	for _, relative := range outputDirs {
		if root, ok := generatedRootDir(relative); ok {
			cleanupDirs = append(cleanupDirs, root)
			continue
		}
		cleanupDirs = append(cleanupDirs, filepath.Clean(relative))
	}
	for _, relative := range uniqueSortedPaths(cleanupDirs) {
		abs := filepath.Join(settings.OutputRoot, relative)
		if err := os.RemoveAll(abs); err != nil {
			return err
		}
	}
	for _, relative := range uniqueSortedPaths(outputDirs) {
		abs := filepath.Join(settings.OutputRoot, relative)
		if err := os.MkdirAll(abs, 0o755); err != nil {
			return err
		}
	}
	return nil
}

func generatedRootDir(relative string) (string, bool) {
	normalized := filepath.ToSlash(filepath.Clean(relative))
	if normalized == "generated" || strings.HasSuffix(normalized, "/generated") {
		return filepath.FromSlash(normalized), true
	}
	if idx := strings.Index(normalized, "/generated/"); idx >= 0 {
		return filepath.FromSlash(normalized[:idx+len("/generated")]), true
	}
	return "", false
}

func uniqueSortedPaths(paths []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(paths))
	for _, path := range paths {
		normalized := filepath.Clean(path)
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	sort.Strings(out)
	return out
}

func generatedGrimoireDebugFunctionDir(settings ExportSettings) string {
	return filepath.Join("data", settings.Namespace, "function", "generated", "debug", "grimoire")
}

func generatedTickFunctionDir(settings ExportSettings) string {
	return filepath.Join("data", settings.Namespace, "function", "generated")
}

func generatedVHFunctionDir(settings ExportSettings) string {
	return filepath.Join(generatedTickFunctionDir(settings), "vh")
}

func generatedVHReplacerFunctionDir(settings ExportSettings) string {
	return filepath.Join(generatedVHFunctionDir(settings), "replacer")
}

func generatedVHReplacerRuleFunctionDir(settings ExportSettings) string {
	return filepath.Join(generatedVHReplacerFunctionDir(settings), "rule")
}

func generateTickDispatcher(settings ExportSettings) (int, error) {
	path := filepath.Join(settings.OutputRoot, generatedTickFunctionDir(settings), "tick.mcfunction")
	body := strings.Join([]string{
		"# generated by tools2: data-driven tick entrypoint",
		fmt.Sprintf("execute as @e[type=#p7b:enemymob,tag=!maf_vh_checked] at @s run function %s", functionResourceID(settings, generatedVHReplacerFunctionDir(settings), "tick")),
		fmt.Sprintf("execute as @e[tag=EnemySkill] at @s run function %s", functionResourceID(settings, settings.Paths.EnemySkillFunctionDir, "main")),
		"",
	}, "\n")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return 0, err
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return 0, err
	}
	return 1, nil
}

func generateItemOutputs(settings ExportSettings, entries []items.ItemEntry) (itemOutputStats, error) {
	functionRoot := filepath.Join(settings.OutputRoot, settings.Paths.ItemFunctionDir)
	lootRoot := filepath.Join(settings.OutputRoot, settings.Paths.ItemLootDir)
	if err := os.MkdirAll(functionRoot, 0o755); err != nil {
		return itemOutputStats{}, err
	}
	if err := os.MkdirAll(lootRoot, 0o755); err != nil {
		return itemOutputStats{}, err
	}
	for _, entry := range entries {
		if err := os.WriteFile(filepath.Join(functionRoot, entry.ID+".mcfunction"), []byte(strings.Join([]string{
			fmt.Sprintf("# itemId=%s sourceId=%s", entry.ItemID, entry.ID),
			fmt.Sprintf("loot give @s loot %s", lootTableResourceID(settings, settings.Paths.ItemLootDir, entry.ID)),
			"",
		}, "\n")), 0o644); err != nil {
			return itemOutputStats{}, err
		}
		if err := writeJSON(filepath.Join(lootRoot, entry.ID+".json"), toItemLootTable(entry, entry.Count)); err != nil {
			return itemOutputStats{}, err
		}
	}
	return itemOutputStats{ItemFunctions: len(entries), ItemLootTables: len(entries)}, nil
}

func generateGrimoireOutputs(settings ExportSettings, entries []grimoire.GrimoireEntry) (spellOutputStats, error) {
	functionRoot := filepath.Join(settings.OutputRoot, settings.Paths.SpellFunctionDir)
	lootRoot := filepath.Join(settings.OutputRoot, settings.Paths.SpellLootDir)
	if err := os.MkdirAll(functionRoot, 0o755); err != nil {
		return spellOutputStats{}, err
	}
	if err := os.MkdirAll(lootRoot, 0o755); err != nil {
		return spellOutputStats{}, err
	}

	dispatchLines := make([]string, 0, len(entries))
	for _, entry := range entries {
		if err := os.WriteFile(filepath.Join(functionRoot, entry.ID+".mcfunction"), []byte(normalizeFunctionBody(entry.Script)), 0o644); err != nil {
			return spellOutputStats{}, err
		}
		if err := writeJSON(filepath.Join(lootRoot, entry.ID+".json"), toSpellLootTable(entry)); err != nil {
			return spellOutputStats{}, err
		}
		dispatchLines = append(dispatchLines, fmt.Sprintf("execute if entity @s[scores={mafEffectID=%d}] run function %s", entry.CastID, functionResourceID(settings, settings.Paths.SpellFunctionDir, entry.ID)))
	}
	if err := os.WriteFile(filepath.Join(functionRoot, "selectexec.mcfunction"), []byte(strings.Join(dispatchLines, "\n")+"\n"), 0o644); err != nil {
		return spellOutputStats{}, err
	}
	return spellOutputStats{SpellFunctions: len(entries) + 1, SpellLootTables: len(entries)}, nil
}

func generateGrimoireDebugFunctions(settings ExportSettings, entries []grimoire.GrimoireEntry) (int, error) {
	root := filepath.Join(settings.OutputRoot, generatedGrimoireDebugFunctionDir(settings))
	if err := os.MkdirAll(root, 0o755); err != nil {
		return 0, err
	}
	for _, entry := range entries {
		if err := os.WriteFile(filepath.Join(root, entry.ID+".mcfunction"), []byte(grimoireDebugGiveCommand(entry)+"\n"), 0o644); err != nil {
			return 0, err
		}
	}
	return len(entries), nil
}

func generateSkillOutputs(settings ExportSettings, entries []skills.SkillEntry) (skillOutputStats, error) {
	root := filepath.Join(settings.OutputRoot, settings.Paths.SkillFunctionDir)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return skillOutputStats{}, err
	}
	for _, entry := range entries {
		if err := os.WriteFile(filepath.Join(root, entry.ID+".mcfunction"), []byte(normalizeFunctionBody(entry.Script)), 0o644); err != nil {
			return skillOutputStats{}, err
		}
	}
	return skillOutputStats{SkillFunctions: len(entries)}, nil
}

func generateEnemySkillOutputs(settings ExportSettings, entries []enemyskills.EnemySkillEntry) (enemySkillOutputStats, error) {
	root := filepath.Join(settings.OutputRoot, settings.Paths.EnemySkillFunctionDir)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return enemySkillOutputStats{}, err
	}
	dispatchLines := []string{"# generated by tools2: enemy skill dispatcher"}
	for _, entry := range entries {
		if err := os.WriteFile(filepath.Join(root, entry.ID+".mcfunction"), []byte(normalizeFunctionBody(entry.Script)), 0o644); err != nil {
			return enemySkillOutputStats{}, err
		}
		dispatchLines = append(dispatchLines, fmt.Sprintf("execute if entity @s[tag=%s] run function %s", entry.ID, functionResourceID(settings, settings.Paths.EnemySkillFunctionDir, entry.ID)))
	}
	if len(entries) == 0 {
		dispatchLines = append(dispatchLines, "# no enemy skill entries")
	}
	if err := os.WriteFile(filepath.Join(root, "main.mcfunction"), []byte(strings.Join(dispatchLines, "\n")+"\n"), 0o644); err != nil {
		return enemySkillOutputStats{}, err
	}
	return enemySkillOutputStats{EnemySkillFunctions: len(entries) + 1}, nil
}

func generateSpawnTableOutputs(settings ExportSettings, entries []spawntables.SpawnTableEntry) (spawnTableOutputStats, error) {
	replacerRoot := filepath.Join(settings.OutputRoot, generatedVHReplacerFunctionDir(settings))
	ruleRoot := filepath.Join(settings.OutputRoot, generatedVHReplacerRuleFunctionDir(settings))
	if err := os.MkdirAll(replacerRoot, 0o755); err != nil {
		return spawnTableOutputStats{}, err
	}
	if err := os.MkdirAll(ruleRoot, 0o755); err != nil {
		return spawnTableOutputStats{}, err
	}

	tickLines := []string{"# generated by tools2: spawn replacement selector"}
	for _, entry := range entries {
		if err := os.WriteFile(filepath.Join(ruleRoot, entry.ID+".mcfunction"), []byte(strings.Join(spawnTableRuleLines(settings, entry), "\n")+"\n"), 0o644); err != nil {
			return spawnTableOutputStats{}, err
		}
		dx := entry.MaxX - entry.MinX
		dy := entry.MaxY - entry.MinY
		dz := entry.MaxZ - entry.MinZ
		tickLines = append(tickLines,
			fmt.Sprintf("execute if entity @s[type=%s] if dimension %s if entity @s[x=%d,y=%d,z=%d,dx=%d,dy=%d,dz=%d] run function %s",
				entry.SourceMobType,
				entry.Dimension,
				entry.MinX, entry.MinY, entry.MinZ,
				dx, dy, dz,
				functionResourceID(settings, generatedVHReplacerRuleFunctionDir(settings), entry.ID),
			),
		)
	}
	tickLines = append(tickLines, "tag @s add maf_vh_checked")
	tickPath := filepath.Join(replacerRoot, "tick.mcfunction")
	if err := os.WriteFile(tickPath, []byte(strings.Join(tickLines, "\n")+"\n"), 0o644); err != nil {
		return spawnTableOutputStats{}, err
	}

	return spawnTableOutputStats{SpawnTableFunctions: len(entries) + 1}, nil
}

func spawnTableRuleLines(settings ExportSettings, entry spawntables.SpawnTableEntry) []string {
	lines := []string{
		fmt.Sprintf("# spawn table %s source=%s dimension=%s", entry.ID, entry.SourceMobType, entry.Dimension),
	}
	totalWeight := entry.BaseMobWeight
	for _, replacement := range entry.Replacements {
		totalWeight += replacement.Weight
	}
	if totalWeight <= 0 {
		return append(lines, "tag @s add maf_vh_checked")
	}
	lines = append(lines, fmt.Sprintf("execute store result score rand p7_Rand1 run random value 0..%d", totalWeight-1))

	start := entry.BaseMobWeight
	for _, replacement := range entry.Replacements {
		end := start + replacement.Weight - 1
		if replacement.Weight > 0 {
			lines = append(lines,
				fmt.Sprintf("execute if score rand p7_Rand1 matches %d..%d at @s run function %s", start, end, functionResourceID(settings, settings.Paths.EnemyFunctionDir, replacement.EnemyID)),
				fmt.Sprintf("execute if score rand p7_Rand1 matches %d..%d run kill @s", start, end),
			)
		}
		start = end + 1
	}
	lines = append(lines, "tag @s add maf_vh_checked")
	return lines
}

func generateEnemyOutputs(settings ExportSettings, entries []enemies.EnemyEntry, itemEntries []items.ItemEntry, grimoireEntries []grimoire.GrimoireEntry) (enemyOutputStats, error) {
	functionRoot := filepath.Join(settings.OutputRoot, settings.Paths.EnemyFunctionDir)
	lootRoot := filepath.Join(settings.OutputRoot, settings.Paths.EnemyLootDir)
	if err := os.MkdirAll(functionRoot, 0o755); err != nil {
		return enemyOutputStats{}, err
	}
	if err := os.MkdirAll(lootRoot, 0o755); err != nil {
		return enemyOutputStats{}, err
	}

	itemsByID := map[string]items.ItemEntry{}
	for _, entry := range itemEntries {
		itemsByID[entry.ID] = entry
	}
	grimoiresByID := map[string]grimoire.GrimoireEntry{}
	for _, entry := range grimoireEntries {
		grimoiresByID[entry.ID] = entry
	}

	for _, entry := range entries {
		lootTable, err := buildDropLootTable(toTreasureDrops(entry.Drops), itemsByID, grimoiresByID, "enemy("+entry.ID+")")
		if err != nil {
			return enemyOutputStats{}, err
		}
		if err := writeJSON(filepath.Join(lootRoot, entry.ID+".json"), lootTable); err != nil {
			return enemyOutputStats{}, err
		}
		if err := os.WriteFile(filepath.Join(functionRoot, entry.ID+".mcfunction"), []byte(strings.Join(toEnemyFunctionLines(settings, entry, itemsByID), "\n")+"\n"), 0o644); err != nil {
			return enemyOutputStats{}, err
		}
	}
	return enemyOutputStats{EnemyFunctions: len(entries), EnemyLootTables: len(entries)}, nil
}

func generateTreasureOutputs(settings ExportSettings, minecraftLootTableRoot string, entries []treasures.TreasureEntry, itemEntries []items.ItemEntry, grimoireEntries []grimoire.GrimoireEntry) (treasureOutputStats, error) {
	if err := cleanupTreasureOverlayOutputs(settings); err != nil {
		return treasureOutputStats{}, err
	}
	writtenPaths := make([]string, 0, len(entries))

	itemsByID := map[string]items.ItemEntry{}
	for _, entry := range itemEntries {
		itemsByID[entry.ID] = entry
	}
	grimoiresByID := map[string]grimoire.GrimoireEntry{}
	for _, entry := range grimoireEntries {
		grimoiresByID[entry.ID] = entry
	}
	for _, entry := range entries {
		if len(entry.LootPools) == 0 {
			return treasureOutputStats{}, fmt.Errorf("treasure(%s): lootPools must not be empty", entry.ID)
		}
		baseLootTable, _, err := mcsource.LoadLootTable(minecraftLootTableRoot, entry.TablePath)
		if err != nil {
			return treasureOutputStats{}, err
		}
		pool, err := buildDropLootPool(entry.LootPools, itemsByID, grimoiresByID, "treasure("+entry.ID+")")
		if err != nil {
			return treasureOutputStats{}, err
		}
		lootTable, err := mergeLootTablePools(baseLootTable, pool, entry.TablePath)
		if err != nil {
			return treasureOutputStats{}, err
		}
		outPath, err := lootTableOutputPath(settings, entry.TablePath)
		if err != nil {
			return treasureOutputStats{}, err
		}
		if err := writeJSON(outPath, lootTable); err != nil {
			return treasureOutputStats{}, err
		}
		writtenPaths = append(writtenPaths, outPath)
	}
	if err := writeTreasureOverlayManifest(settings, writtenPaths); err != nil {
		return treasureOutputStats{}, err
	}
	return treasureOutputStats{TreasureLootTables: len(entries)}, nil
}

func generateLootTableOutputs(settings ExportSettings, entries []loottables.LootTableEntry, itemEntries []items.ItemEntry, grimoireEntries []grimoire.GrimoireEntry) (loottableOutputStats, error) {
	lootRoot := filepath.Join(settings.OutputRoot, settings.Paths.LoottableLootDir)
	if err := os.MkdirAll(lootRoot, 0o755); err != nil {
		return loottableOutputStats{}, err
	}
	itemsByID := map[string]items.ItemEntry{}
	for _, entry := range itemEntries {
		itemsByID[entry.ID] = entry
	}
	grimoiresByID := map[string]grimoire.GrimoireEntry{}
	for _, entry := range grimoireEntries {
		grimoiresByID[entry.ID] = entry
	}
	for _, entry := range entries {
		if len(entry.LootPools) == 0 {
			return loottableOutputStats{}, fmt.Errorf("loottable(%s): lootPools must not be empty", entry.ID)
		}
		lootTable, err := buildDropLootTable(entry.LootPools, itemsByID, grimoiresByID, "loottable("+entry.ID+")")
		if err != nil {
			return loottableOutputStats{}, err
		}
		if err := writeJSON(filepath.Join(lootRoot, entry.ID+".json"), lootTable); err != nil {
			return loottableOutputStats{}, err
		}
	}
	return loottableOutputStats{LoottableLootTables: len(entries)}, nil
}

func toTreasureDrops(drops []enemies.DropRef) []treasures.DropRef {
	out := make([]treasures.DropRef, 0, len(drops))
	for _, drop := range drops {
		out = append(out, treasures.DropRef{
			Kind:     drop.Kind,
			RefID:    drop.RefID,
			Weight:   drop.Weight,
			CountMin: drop.CountMin,
			CountMax: drop.CountMax,
		})
	}
	return out
}

func buildDropLootTable(drops []treasures.DropRef, itemsByID map[string]items.ItemEntry, grimoiresByID map[string]grimoire.GrimoireEntry, context string) (map[string]any, error) {
	pool, err := buildDropLootPool(drops, itemsByID, grimoiresByID, context)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"type":  "minecraft:generic",
		"pools": []any{pool},
	}, nil
}

func buildDropLootPool(drops []treasures.DropRef, itemsByID map[string]items.ItemEntry, grimoiresByID map[string]grimoire.GrimoireEntry, context string) (map[string]any, error) {
	entries := make([]any, 0, len(drops))
	for _, drop := range drops {
		switch drop.Kind {
		case "minecraft_item":
			entries = append(entries, map[string]any{
				"type":   "minecraft:item",
				"name":   drop.RefID,
				"weight": toWeight(drop.Weight),
				"functions": []any{
					map[string]any{"function": "minecraft:set_count", "count": toCountValue(drop.CountMin, drop.CountMax)},
				},
			})
		case "item":
			item, ok := itemsByID[drop.RefID]
			if !ok {
				return nil, fmt.Errorf("%s: referenced item not found (%s)", context, drop.RefID)
			}
			entry := toItemLootEntry(item, drop.CountMin, drop.CountMax)
			entry["weight"] = toWeight(drop.Weight)
			entries = append(entries, entry)
		case "grimoire":
			entry, ok := grimoiresByID[drop.RefID]
			if !ok {
				return nil, fmt.Errorf("%s: referenced grimoire not found (%s)", context, drop.RefID)
			}
			out := toSpellLootEntry(entry, drop.CountMin, drop.CountMax)
			out["weight"] = toWeight(drop.Weight)
			entries = append(entries, out)
		default:
			return nil, fmt.Errorf("%s: unsupported drop kind (%s)", context, drop.Kind)
		}
	}
	return map[string]any{
		"rolls":   1,
		"entries": entries,
	}, nil
}

func mergeLootTablePools(base map[string]any, pool map[string]any, tablePath string) (map[string]any, error) {
	if base == nil {
		base = map[string]any{}
	}
	if rawPools, ok := base["pools"]; ok && rawPools != nil {
		pools, ok := rawPools.([]any)
		if !ok {
			return nil, fmt.Errorf("treasure(%s): base loot table pools must be an array", tablePath)
		}
		base["pools"] = append(pools, pool)
		return base, nil
	}
	base["pools"] = []any{pool}
	return base, nil
}

func toEnemyFunctionLines(settings ExportSettings, entry enemies.EnemyEntry, itemsByID map[string]items.ItemEntry) []string {
	lootID := lootTableResourceID(settings, settings.Paths.EnemyLootDir, entry.ID)
	lines := []string{
		fmt.Sprintf("# enemyId=%s mobType=%s", entry.ID, entry.MobType),
		fmt.Sprintf("# dropMode=%s", entry.DropMode),
		fmt.Sprintf("summon %s ~ ~ ~ %s", entry.MobType, enemySummonNBT(lootID, entry, itemsByID)),
	}
	return lines
}

func enemySummonNBT(lootID string, entry enemies.EnemyEntry, itemsByID map[string]items.ItemEntry) string {
	parts := []string{
		fmt.Sprintf("Health:%sf", formatFloat(entry.HP)),
		fmt.Sprintf("DeathLootTable:%s", jsonString(lootID)),
	}
	if entry.Name != "" {
		parts = append(parts, fmt.Sprintf("CustomName:%s", singleQuotedJSON(map[string]string{"text": entry.Name})))
	}
	if tags := enemyTags(entry); len(tags) > 0 {
		parts = append(parts, fmt.Sprintf("Tags:[%s]", strings.Join(tags, ",")))
	}
	if attrs := enemyAttributes(entry); len(attrs) > 0 {
		parts = append(parts, fmt.Sprintf("Attributes:[%s]", strings.Join(attrs, ",")))
	}
	if handItems, handDrops := equipmentArray(itemsByID, entry.Equipment.Mainhand, entry.Equipment.Offhand); handItems != "" {
		parts = append(parts, "HandItems:["+handItems+"]", "HandDropChances:["+handDrops+"]")
	}
	if armorItems, armorDrops := equipmentArray(itemsByID, entry.Equipment.Feet, entry.Equipment.Legs, entry.Equipment.Chest, entry.Equipment.Head); armorItems != "" {
		parts = append(parts, "ArmorItems:["+armorItems+"]", "ArmorDropChances:["+armorDrops+"]")
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func enemyTags(entry enemies.EnemyEntry) []string {
	tags := []string{jsonString("maf_enemy"), jsonString("maf_enemy_" + entry.ID), jsonString("maf_vh_checked")}
	if len(entry.EnemySkillIDs) > 0 {
		tags = append(tags, jsonString("EnemySkill"))
	}
	for _, skillID := range entry.EnemySkillIDs {
		tags = append(tags, jsonString(skillID), jsonString("maf_enemy_skill_"+skillID))
	}
	return tags
}

func enemyAttributes(entry enemies.EnemyEntry) []string {
	attrs := []string{
		fmt.Sprintf("{Name:generic.max_health,Base:%s}", formatFloat(entry.HP)),
	}
	if entry.Attack != nil {
		attrs = append(attrs, fmt.Sprintf("{Name:generic.attack_damage,Base:%s}", formatFloat(*entry.Attack)))
	}
	if entry.Defense != nil {
		attrs = append(attrs, fmt.Sprintf("{Name:generic.armor,Base:%s}", formatFloat(*entry.Defense)))
	}
	if entry.MoveSpeed != nil {
		attrs = append(attrs, fmt.Sprintf("{Name:generic.movement_speed,Base:%s}", formatFloat(*entry.MoveSpeed)))
	}
	return attrs
}

func equipmentArray(itemsByID map[string]items.ItemEntry, slots ...*enemies.EquipmentSlot) (string, string) {
	itemsOut := make([]string, 0, len(slots))
	dropsOut := make([]string, 0, len(slots))
	for _, slot := range slots {
		if slot == nil {
			itemsOut = append(itemsOut, "{}")
			dropsOut = append(dropsOut, "0.085F")
			continue
		}
		itemsOut = append(itemsOut, fmt.Sprintf("{id:%s,Count:%db}", jsonString(resolveEquipmentItemID(slot, itemsByID)), slot.Count))
		dropChance := 0.085
		if slot.DropChance != nil {
			dropChance = *slot.DropChance
		}
		dropsOut = append(dropsOut, formatFloat(dropChance)+"F")
	}
	return strings.Join(itemsOut, ","), strings.Join(dropsOut, ",")
}

func resolveEquipmentItemID(slot *enemies.EquipmentSlot, itemsByID map[string]items.ItemEntry) string {
	if slot == nil {
		return ""
	}
	if slot.Kind == "item" {
		if entry, ok := itemsByID[slot.RefID]; ok && entry.ItemID != "" {
			return entry.ItemID
		}
	}
	return slot.RefID
}

func normalizeFunctionBody(script string) string {
	if strings.HasSuffix(script, "\n") {
		return script
	}
	return script + "\n"
}

func toSpellLootTable(entry grimoire.GrimoireEntry) map[string]any {
	return map[string]any{
		"type": "minecraft:generic",
		"pools": []any{
			map[string]any{
				"rolls": 1,
				"entries": []any{
					toSpellLootEntry(entry, ptrFloat(1), ptrFloat(1)),
				},
			},
		},
	}
}

func toItemLootTable(entry items.ItemEntry, count int) map[string]any {
	value := float64(count)
	return map[string]any{
		"type": "minecraft:generic",
		"pools": []any{
			map[string]any{
				"rolls": 1,
				"entries": []any{
					toItemLootEntry(entry, &value, &value),
				},
			},
		},
	}
}

func toItemLootEntry(entry items.ItemEntry, min, max *float64) map[string]any {
	return map[string]any{
		"type": "minecraft:item",
		"name": entry.ItemID,
		"functions": []any{
			map[string]any{"function": "minecraft:set_count", "add": false, "count": toCountValue(min, max)},
			map[string]any{"function": "minecraft:set_custom_data", "tag": itemCustomData(entry)},
		},
	}
}

func toSpellLootEntry(entry grimoire.GrimoireEntry, min, max *float64) map[string]any {
	functions := []any{
		map[string]any{"function": "minecraft:set_count", "add": false, "count": toCountValue(min, max)},
		map[string]any{"function": "minecraft:set_name", "name": map[string]any{"text": entry.Title}, "target": "item_name"},
	}
	if lore := toLoreComponents(entry.Description); len(lore) > 0 {
		functions = append(functions, map[string]any{"function": "minecraft:set_lore", "mode": "append", "lore": lore})
	}
	functions = append(functions, map[string]any{"function": "minecraft:set_custom_data", "tag": spellCustomData(entry)})
	return map[string]any{
		"type":      "minecraft:item",
		"name":      "minecraft:written_book",
		"functions": functions,
	}
}

func itemCustomData(entry items.ItemEntry) string {
	parts := []string{
		fmt.Sprintf("item_id:%s", jsonString(entry.ItemID)),
		fmt.Sprintf("source_id:%s", jsonString(entry.ID)),
		fmt.Sprintf("nbt_snapshot:%s", jsonString(entry.NBT)),
	}
	if entry.SkillID != "" {
		parts = append(parts, "maf_skill:1b", fmt.Sprintf("maf_skill_id:%s", jsonString(entry.SkillID)))
	}
	return "{maf:{" + strings.Join(parts, ",") + "}}"
}

func spellCustomData(entry grimoire.GrimoireEntry) string {
	return fmt.Sprintf("{maf:{grimoire_id:%s,spell:{castid:%d,cost:%d,cast:%d,title:%s,description:%s}}}", jsonString(entry.ID), entry.CastID, entry.MPCost, entry.CastTime, jsonString(entry.Title), jsonString(entry.Description))
}

func grimoireDebugGiveCommand(entry grimoire.GrimoireEntry) string {
	parts := []string{
		fmt.Sprintf("item_name=%s", singleQuotedJSON(map[string]any{"text": entry.Title})),
	}
	if loreLines := linesToLoreValues(entry.Description); len(loreLines) > 0 {
		loreParts := make([]string, 0, len(loreLines))
		for _, line := range loreLines {
			loreParts = append(loreParts, singleQuotedJSON(map[string]any{"text": line}))
		}
		parts = append(parts, "lore=["+strings.Join(loreParts, ",")+"]")
	}
	parts = append(parts, "custom_data="+spellCustomData(entry))
	return fmt.Sprintf("give @s minecraft:written_book[%s] 1", strings.Join(parts, ","))
}

func toLoreComponents(value string) []any {
	lines := linesToLoreValues(value)
	out := make([]any, 0, len(lines))
	for _, line := range lines {
		out = append(out, map[string]any{"text": line})
	}
	return out
}

func linesToLoreValues(value string) []string {
	raw := strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func toCountValue(min, max *float64) any {
	minValue := 1.0
	maxValue := 1.0
	if min != nil {
		minValue = *min
	}
	if max != nil {
		maxValue = *max
	}
	if minValue == maxValue {
		return minValue
	}
	return map[string]any{
		"type": "minecraft:uniform",
		"min":  minValue,
		"max":  maxValue,
	}
}

func toWeight(weight float64) int {
	if !isFinite(weight) || weight <= 0 {
		return 1
	}
	return int(math.Floor(weight))
}

func isFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func functionResourceID(settings ExportSettings, relativeDir, baseName string) string {
	resourcePath := strings.TrimPrefix(filepath.ToSlash(relativeDir), "data/"+settings.Namespace+"/function/")
	resourcePath = strings.Trim(resourcePath, "/")
	if resourcePath == "" {
		return settings.Namespace + ":" + baseName
	}
	return settings.Namespace + ":" + resourcePath + "/" + baseName
}

func lootTableResourceID(settings ExportSettings, relativeDir, baseName string) string {
	prefix := "data/" + settings.Namespace + "/loot_table/"
	resourcePath := strings.TrimPrefix(filepath.ToSlash(relativeDir), prefix)
	resourcePath = strings.Trim(resourcePath, "/")
	if resourcePath == "" {
		return settings.Namespace + ":" + baseName
	}
	return settings.Namespace + ":" + resourcePath + "/" + baseName
}

func lootTableOutputPath(settings ExportSettings, tablePath string) (string, error) {
	if !common.IsSafeNamespacedResourcePath(tablePath) {
		return "", fmt.Errorf("invalid loot table path: %s", tablePath)
	}
	parts := strings.SplitN(tablePath, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("invalid loot table path: %s", tablePath)
	}
	return filepath.Join(settings.OutputRoot, "data", parts[0], "loot_table", filepath.FromSlash(common.NormalizeResourcePath(parts[1]))+".json"), nil
}

func jsonString(value string) string {
	return string(mustJSON(value))
}

func singleQuotedJSON(value any) string {
	return "'" + strings.ReplaceAll(string(mustJSON(value)), "'", "\\'") + "'"
}

func mustJSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return data
}

func ptrFloat(v float64) *float64 {
	return &v
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func treasureOverlayManifestPath(settings ExportSettings) string {
	return filepath.Join(settings.OutputRoot, ".tools2", "treasure-overrides.json")
}

func cleanupTreasureOverlayOutputs(settings ExportSettings) error {
	manifest, err := readTreasureOverlayManifest(settings)
	if err != nil {
		return err
	}
	for _, rel := range manifest.Paths {
		abs, err := safeOutputPath(settings.OutputRoot, rel)
		if err != nil {
			return err
		}
		if err := os.Remove(abs); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func readTreasureOverlayManifest(settings ExportSettings) (treasureOverlayManifest, error) {
	path := treasureOverlayManifestPath(settings)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return treasureOverlayManifest{}, nil
		}
		return treasureOverlayManifest{}, err
	}
	var manifest treasureOverlayManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return treasureOverlayManifest{}, fmt.Errorf("invalid treasure overlay manifest: %w", err)
	}
	return manifest, nil
}

func writeTreasureOverlayManifest(settings ExportSettings, absPaths []string) error {
	manifest := treasureOverlayManifest{Paths: make([]string, 0, len(absPaths))}
	for _, absPath := range absPaths {
		rel, err := filepath.Rel(settings.OutputRoot, absPath)
		if err != nil {
			return err
		}
		manifest.Paths = append(manifest.Paths, filepath.ToSlash(rel))
	}
	sort.Strings(manifest.Paths)
	path := treasureOverlayManifestPath(settings)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func safeOutputPath(outputRoot, relPath string) (string, error) {
	normalized := filepath.Clean(filepath.FromSlash(relPath))
	if normalized == "." {
		return "", fmt.Errorf("invalid relative path in treasure overlay manifest: %q", relPath)
	}
	joined := filepath.Clean(filepath.Join(outputRoot, normalized))
	root := filepath.Clean(outputRoot)
	if joined == root {
		return "", fmt.Errorf("invalid relative path in treasure overlay manifest: %q", relPath)
	}
	prefix := root + string(os.PathSeparator)
	if !strings.HasPrefix(joined, prefix) {
		return "", fmt.Errorf("path escapes output root in treasure overlay manifest: %q", relPath)
	}
	return joined, nil
}

func exportFailure(code, message string, err error) SaveDataResponse {
	return SaveDataResponse{
		OK:      false,
		Code:    code,
		Message: message,
		Details: err.Error(),
	}
}
