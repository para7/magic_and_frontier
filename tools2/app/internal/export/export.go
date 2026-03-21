package export

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
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
	treasureStats, err := generateTreasureOutputs(settings, params.MinecraftLootTableRoot, params.Treasures, params.ItemState.Items, params.GrimoireState.Entries)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	loottableStats, err := generateLootTableOutputs(settings, params.LootTables, params.ItemState.Items, params.GrimoireState.Entries)
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
	stats.TotalFiles = stats.ItemFunctions + stats.ItemLootTables + stats.SpellFunctions + stats.SpellLootTables + stats.SkillFunctions + stats.EnemySkillFunctions + stats.EnemyFunctions + stats.EnemyLootTables + stats.TreasureLootTables + stats.LoottableLootTables + debugGrimoireFunctions + 3

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

type treasureOutputStats struct {
	TreasureLootTables int
}

type loottableOutputStats struct {
	LoottableLootTables int
}

func loadExportSettings(settingsPath string) (ExportSettings, error) {
	var parsed map[string]any
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return ExportSettings{}, fmt.Errorf("failed to read export settings at %s: %w", settingsPath, err)
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return ExportSettings{}, fmt.Errorf("failed to read export settings at %s: %w", settingsPath, err)
	}

	outputRoot, err := requireString(parsed["outputRoot"], "outputRoot")
	if err != nil {
		return ExportSettings{}, err
	}
	namespace, err := requireString(parsed["namespace"], "namespace")
	if err != nil {
		return ExportSettings{}, err
	}
	if !namespacePattern(namespace) {
		return ExportSettings{}, fmt.Errorf("namespace must match [a-z0-9_.-]+")
	}
	templatePackPath, err := requireString(parsed["templatePackPath"], "templatePackPath")
	if err != nil {
		return ExportSettings{}, err
	}

	pathValues, ok := parsed["paths"].(map[string]any)
	if !ok {
		return ExportSettings{}, fmt.Errorf("paths must be an object")
	}

	baseDir := filepath.Dir(settingsPath)
	itemFunctionDir, err := requirePathString(pathValues, "itemFunctionDir")
	if err != nil {
		return ExportSettings{}, err
	}
	itemLootDir, err := requirePathString(pathValues, "itemLootDir")
	if err != nil {
		return ExportSettings{}, err
	}
	spellFunctionDir, err := requirePathString(pathValues, "spellFunctionDir")
	if err != nil {
		return ExportSettings{}, err
	}
	spellLootDir, err := requirePathString(pathValues, "spellLootDir")
	if err != nil {
		return ExportSettings{}, err
	}
	skillFunctionDir, err := pathStringOrDefault(pathValues, "skillFunctionDir", defaultGeneratedPath(namespace, "function", "skill"))
	if err != nil {
		return ExportSettings{}, err
	}
	enemySkillFunctionDir, err := pathStringOrDefault(pathValues, "enemySkillFunctionDir", defaultGeneratedPath(namespace, "function", "enemy_skill"))
	if err != nil {
		return ExportSettings{}, err
	}
	enemyFunctionDir, err := pathStringOrDefault(pathValues, "enemyFunctionDir", defaultGeneratedPath(namespace, "function", "enemy"))
	if err != nil {
		return ExportSettings{}, err
	}
	enemyLootDir, err := pathStringOrDefault(pathValues, "enemyLootDir", defaultGeneratedPath(namespace, "loot_table", "enemy"))
	if err != nil {
		return ExportSettings{}, err
	}
	treasureLootDir, err := pathStringOrDefault(pathValues, "treasureLootDir", defaultGeneratedPath(namespace, "loot_table", "treasure"))
	if err != nil {
		return ExportSettings{}, err
	}
	loottableLootDir, err := pathStringOrDefault(pathValues, "loottableLootDir", defaultGeneratedPath(namespace, "loot_table", "loottable"))
	if err != nil {
		return ExportSettings{}, err
	}
	debugFunctionDir, err := pathStringOrDefault(pathValues, "debugFunctionDir", filepath.ToSlash(filepath.Join("data", namespace, "function", "debug", "give")))
	if err != nil {
		return ExportSettings{}, err
	}
	minecraftTagDir, err := requirePathString(pathValues, "minecraftTagDir")
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

func requireString(value any, key string) (string, error) {
	text, ok := value.(string)
	if !ok || strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	return text, nil
}

func requirePathString(values map[string]any, key string) (string, error) {
	return requireString(values[key], "paths."+key)
}

func pathStringOrDefault(values map[string]any, key, fallback string) (string, error) {
	if _, ok := values[key]; !ok {
		return fallback, nil
	}
	return requireString(values[key], "paths."+key)
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
	templateData, err := os.ReadFile(settings.TemplatePackPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(settings.OutputRoot, "pack.mcmeta"), templateData, 0o644); err != nil {
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
		filepath.Join("data", "minecraft", "loot_table"),
		filepath.Join(settings.Paths.DebugFunctionDir, "item"),
		filepath.Join(settings.Paths.DebugFunctionDir, "grimoire"),
		generatedGrimoireDebugFunctionDir(settings),
	}
	for _, relative := range outputDirs {
		abs := filepath.Join(settings.OutputRoot, relative)
		if err := os.RemoveAll(abs); err != nil {
			return err
		}
		if err := os.MkdirAll(abs, 0o755); err != nil {
			return err
		}
	}

	if err := writeFunctionTag(settings, "load", settings.Namespace+":load"); err != nil {
		return err
	}
	return writeFunctionTag(settings, "tick", settings.Namespace+":tick")
}

func generatedGrimoireDebugFunctionDir(settings ExportSettings) string {
	return filepath.Join("data", settings.Namespace, "function", "generated", "debug", "grimoire")
}

func writeFunctionTag(settings ExportSettings, tagName string, values ...string) error {
	tagPath := filepath.Join(settings.OutputRoot, settings.Paths.MinecraftTagDir, tagName+".json")
	if err := os.MkdirAll(filepath.Dir(tagPath), 0o755); err != nil {
		return err
	}
	return writeJSON(tagPath, map[string]any{"values": values})
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
	for _, entry := range entries {
		if err := os.WriteFile(filepath.Join(root, entry.ID+".mcfunction"), []byte(normalizeFunctionBody(entry.Script)), 0o644); err != nil {
			return enemySkillOutputStats{}, err
		}
	}
	return enemySkillOutputStats{EnemySkillFunctions: len(entries)}, nil
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
	tags := []string{jsonString("maf_enemy"), jsonString("maf_enemy_" + entry.ID)}
	for _, skillID := range entry.EnemySkillIDs {
		tags = append(tags, jsonString("maf_enemy_skill_"+skillID))
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

func exportFailure(code, message string, err error) SaveDataResponse {
	return SaveDataResponse{
		OK:      false,
		Code:    code,
		Message: message,
		Details: err.Error(),
	}
}
