package export

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/skills"
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
	ItemState          items.ItemState
	GrimoireState      grimoire.GrimoireState
	Skills             []skills.SkillEntry
	EnemySkills        []enemyskills.EnemySkillEntry
	Enemies            []enemies.EnemyEntry
	Treasures          []treasures.TreasureEntry
	ExportSettingsPath string
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
	skillStats, err := generateSkillOutputs(settings, params.Skills)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	enemySkillStats, err := generateEnemySkillOutputs(settings, params.EnemySkills)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	enemyStats, err := generateEnemyOutputs(settings, params.Enemies, params.Treasures, params.ItemState.Items, params.GrimoireState.Entries)
	if err != nil {
		return exportFailure("EXPORT_FAILED", "Datapack export failed.", err)
	}
	treasureStats, err := generateTreasureOutputs(settings, params.Treasures, params.ItemState.Items, params.GrimoireState.Entries)
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
	}
	stats.TotalFiles = stats.ItemFunctions + stats.ItemLootTables + stats.SpellFunctions + stats.SpellLootTables + stats.SkillFunctions + stats.EnemySkillFunctions + stats.EnemyFunctions + stats.EnemyLootTables + stats.TreasureLootTables + 3

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
	skillFunctionDir, err := pathStringOrDefault(pathValues, "skillFunctionDir", filepath.ToSlash(filepath.Join("data", namespace, "function", "skill")))
	if err != nil {
		return ExportSettings{}, err
	}
	enemySkillFunctionDir, err := pathStringOrDefault(pathValues, "enemySkillFunctionDir", filepath.ToSlash(filepath.Join("data", namespace, "function", "enemy_skill")))
	if err != nil {
		return ExportSettings{}, err
	}
	enemyFunctionDir, err := pathStringOrDefault(pathValues, "enemyFunctionDir", filepath.ToSlash(filepath.Join("data", namespace, "function", "enemy", "spawn")))
	if err != nil {
		return ExportSettings{}, err
	}
	enemyLootDir, err := pathStringOrDefault(pathValues, "enemyLootDir", filepath.ToSlash(filepath.Join("data", namespace, "loot_table", "enemy")))
	if err != nil {
		return ExportSettings{}, err
	}
	treasureLootDir, err := pathStringOrDefault(pathValues, "treasureLootDir", filepath.ToSlash(filepath.Join("data", namespace, "loot_table", "treasure")))
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

	settings := ExportSettings{
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
			DebugFunctionDir:      debugFunctionDir,
			MinecraftTagDir:       minecraftTagDir,
		},
	}
	return settings, nil
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
		filepath.Join(settings.Paths.DebugFunctionDir, "item"),
		filepath.Join(settings.Paths.DebugFunctionDir, "grimoire"),
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

	loadTagPath := filepath.Join(settings.OutputRoot, settings.Paths.MinecraftTagDir, "load.json")
	if err := os.MkdirAll(filepath.Dir(loadTagPath), 0o755); err != nil {
		return err
	}
	loadTag := map[string]any{"values": []string{settings.Namespace + ":load"}}
	if err := writeJSON(loadTagPath, loadTag); err != nil {
		return err
	}

	loadFuncPath := filepath.Join(settings.OutputRoot, "data", settings.Namespace, "function", "load.mcfunction")
	if err := os.MkdirAll(filepath.Dir(loadFuncPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(loadFuncPath, []byte(fmt.Sprintf("tellraw @a [{\"text\":\"enabled datapack: %s\"}]\n", settings.Namespace)), 0o644)
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
		functionPath := filepath.Join(functionRoot, entry.ID+".mcfunction")
		lootPath := filepath.Join(lootRoot, entry.ID+".json")
		lines := []string{
			fmt.Sprintf("# itemId=%s sourceId=%s", entry.ItemID, entry.ID),
			fmt.Sprintf("loot give @s loot %s:item/%s", settings.Namespace, entry.ID),
		}
		if err := os.WriteFile(functionPath, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
			return itemOutputStats{}, err
		}
		lootTable := map[string]any{
			"type": "minecraft:generic",
			"pools": []any{
				map[string]any{
					"rolls": 1,
					"entries": []any{
						map[string]any{
							"type": "minecraft:item",
							"name": entry.ItemID,
							"functions": []any{
								map[string]any{
									"function": "minecraft:set_count",
									"count":    entry.Count,
								},
								map[string]any{
									"function": "minecraft:set_custom_data",
									"tag":      fmt.Sprintf("{maf:{item_id:%s,source_id:%s,nbt_snapshot:%s}}", jsonString(entry.ItemID), jsonString(entry.ID), jsonString(entry.NBT)),
								},
							},
						},
					},
				},
			},
		}
		if err := writeJSON(lootPath, lootTable); err != nil {
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
	count := 0
	for _, entry := range entries {
		for idx, variant := range entry.Variants {
			count++
			base := fmt.Sprintf("cast_%d_v%d", entry.CastID, idx+1)
			functionPath := filepath.Join(functionRoot, base+".mcfunction")
			lootPath := filepath.Join(lootRoot, base+".json")
			lines := []string{
				fmt.Sprintf("# castid=%d variant=%d sourceId=%s", entry.CastID, idx+1, entry.ID),
				toSpellGiveCommand(entry, variant),
			}
			if err := os.WriteFile(functionPath, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
				return spellOutputStats{}, err
			}
			if err := writeJSON(lootPath, toSpellLootTable(entry, variant)); err != nil {
				return spellOutputStats{}, err
			}
		}
	}
	return spellOutputStats{SpellFunctions: count, SpellLootTables: count}, nil
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

func generateEnemyOutputs(settings ExportSettings, entries []enemies.EnemyEntry, treasuresState []treasures.TreasureEntry, itemEntries []items.ItemEntry, grimoireEntries []grimoire.GrimoireEntry) (enemyOutputStats, error) {
	functionRoot := filepath.Join(settings.OutputRoot, settings.Paths.EnemyFunctionDir)
	lootRoot := filepath.Join(settings.OutputRoot, settings.Paths.EnemyLootDir)
	if err := os.MkdirAll(functionRoot, 0o755); err != nil {
		return enemyOutputStats{}, err
	}
	if err := os.MkdirAll(lootRoot, 0o755); err != nil {
		return enemyOutputStats{}, err
	}
	treasuresByID := map[string]treasures.TreasureEntry{}
	for _, entry := range treasuresState {
		treasuresByID[entry.ID] = entry
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
		functionPath := filepath.Join(functionRoot, entry.ID+".mcfunction")
		lootPath := filepath.Join(lootRoot, entry.ID+".json")
		drops, err := resolveEnemyDrops(entry, treasuresByID)
		if err != nil {
			return enemyOutputStats{}, err
		}
		if err := os.WriteFile(functionPath, []byte(strings.Join(toEnemyFunctionLines(entry), "\n")+"\n"), 0o644); err != nil {
			return enemyOutputStats{}, err
		}
		lootTable, err := buildDropLootTable(drops, itemsByID, grimoiresByID, "enemy("+entry.ID+")")
		if err != nil {
			return enemyOutputStats{}, err
		}
		if err := writeJSON(lootPath, lootTable); err != nil {
			return enemyOutputStats{}, err
		}
	}
	return enemyOutputStats{EnemyFunctions: len(entries), EnemyLootTables: len(entries)}, nil
}

func generateTreasureOutputs(settings ExportSettings, entries []treasures.TreasureEntry, itemEntries []items.ItemEntry, grimoireEntries []grimoire.GrimoireEntry) (treasureOutputStats, error) {
	lootRoot := filepath.Join(settings.OutputRoot, settings.Paths.TreasureLootDir)
	if err := os.MkdirAll(lootRoot, 0o755); err != nil {
		return treasureOutputStats{}, err
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
			return treasureOutputStats{}, fmt.Errorf("treasure(%s): lootPools must not be empty", entry.ID)
		}
		lootTable, err := buildDropLootTable(entry.LootPools, itemsByID, grimoiresByID, "treasure("+entry.ID+")")
		if err != nil {
			return treasureOutputStats{}, err
		}
		if err := writeJSON(filepath.Join(lootRoot, entry.ID+".json"), lootTable); err != nil {
			return treasureOutputStats{}, err
		}
	}
	return treasureOutputStats{TreasureLootTables: len(entries)}, nil
}

func resolveEnemyDrops(entry enemies.EnemyEntry, treasuresByID map[string]treasures.TreasureEntry) ([]treasures.DropRef, error) {
	if len(entry.DropTable) > 0 {
		return toTreasureDrops(entry.DropTable), nil
	}
	if treasure, ok := treasuresByID[entry.DropTableID]; ok && len(treasure.LootPools) > 0 {
		return treasure.LootPools, nil
	}
	return nil, fmt.Errorf("enemy(%s): drop table was not found", entry.ID)
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
	entries := make([]any, 0, len(drops))
	for _, drop := range drops {
		if drop.Kind == "item" {
			item, ok := itemsByID[drop.RefID]
			if !ok {
				return nil, fmt.Errorf("%s: referenced item not found (%s)", context, drop.RefID)
			}
			entries = append(entries, map[string]any{
				"type":   "minecraft:item",
				"name":   item.ItemID,
				"weight": toWeight(drop.Weight),
				"functions": []any{
					map[string]any{
						"function": "minecraft:set_count",
						"count":    toCountValue(drop.CountMin, drop.CountMax),
					},
					map[string]any{
						"function": "minecraft:set_custom_data",
						"tag":      fmt.Sprintf("{maf:{source_id:%s}}", jsonString(drop.RefID)),
					},
				},
			})
			continue
		}
		entry, ok := grimoiresByID[drop.RefID]
		if !ok {
			return nil, fmt.Errorf("%s: referenced grimoire not found (%s)", context, drop.RefID)
		}
		functions := []any{
			map[string]any{
				"function": "minecraft:set_count",
				"count":    toCountValue(drop.CountMin, drop.CountMax),
			},
			map[string]any{
				"function": "minecraft:set_name",
				"name":     map[string]any{"text": entry.Title},
			},
		}
		if lore := toLoreComponents(entry.Description); len(lore) > 0 {
			functions = append(functions, map[string]any{
				"function": "minecraft:set_lore",
				"lore":     lore,
			})
		}
		functions = append(functions, map[string]any{
			"function": "minecraft:set_custom_data",
			"tag":      fmt.Sprintf("{maf:{source_id:%s,%s}}", jsonString(drop.RefID), toPrimarySpellCustomData(entry)),
		})
		entries = append(entries, map[string]any{
			"type":      "minecraft:item",
			"name":      "minecraft:written_book",
			"weight":    toWeight(drop.Weight),
			"functions": functions,
		})
	}
	return map[string]any{
		"type": "minecraft:generic",
		"pools": []any{
			map[string]any{
				"rolls":   1,
				"entries": entries,
			},
		},
	}, nil
}

func toEnemyFunctionLines(entry enemies.EnemyEntry) []string {
	lines := []string{
		fmt.Sprintf("# enemyId=%s name=%s", entry.ID, firstNonEmpty(entry.Name, entry.ID)),
		fmt.Sprintf("# distance=%v..%v", entry.SpawnRule.Distance.Min, entry.SpawnRule.Distance.Max),
		fmt.Sprintf("execute positioned %v %v %v run summon minecraft:zombie ~ ~ ~", entry.SpawnRule.Origin.X, entry.SpawnRule.Origin.Y, entry.SpawnRule.Origin.Z),
	}
	if entry.SpawnRule.AxisBounds != nil {
		raw, _ := json.Marshal(entry.SpawnRule.AxisBounds)
		lines = append([]string{lines[0], lines[1], "# axisBounds=" + strings.ReplaceAll(string(raw), " ", "")}, lines[2:]...)
	}
	return lines
}

func firstNonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func normalizeFunctionBody(script string) string {
	if strings.HasSuffix(script, "\n") {
		return script
	}
	return script + "\n"
}

func toSpellLootTable(entry grimoire.GrimoireEntry, variant grimoire.Variant) map[string]any {
	functions := []any{
		map[string]any{
			"function": "minecraft:set_name",
			"name":     map[string]any{"text": entry.Title},
		},
	}
	if lore := toLoreComponents(entry.Description); len(lore) > 0 {
		functions = append(functions, map[string]any{
			"function": "minecraft:set_lore",
			"lore":     lore,
		})
	}
	functions = append(functions, map[string]any{
		"function": "minecraft:set_custom_data",
		"tag":      fmt.Sprintf("{%s}", toSpellCustomData(entry, variant)),
	})
	return map[string]any{
		"type": "minecraft:generic",
		"pools": []any{
			map[string]any{
				"rolls": 1,
				"entries": []any{
					map[string]any{
						"type":      "minecraft:item",
						"name":      "minecraft:written_book",
						"functions": functions,
					},
				},
			},
		},
	}
}

func toSpellGiveCommand(entry grimoire.GrimoireEntry, variant grimoire.Variant) string {
	customName := "'" + strings.ReplaceAll(string(mustJSON(map[string]string{"text": entry.Title})), "'", "\\'") + "'"
	loreParts := make([]string, 0)
	for _, line := range linesToLoreValues(entry.Description) {
		loreParts = append(loreParts, "'"+strings.ReplaceAll(string(mustJSON(map[string]string{"text": line})), "'", "\\'")+"'")
	}
	loreValue := ""
	if len(loreParts) > 0 {
		loreValue = ",lore:[" + strings.Join(loreParts, ",") + "]"
	}
	return fmt.Sprintf("give @s minecraft:written_book[custom_name=%s%s,custom_data={%s}] 1", customName, loreValue, toSpellCustomData(entry, variant))
}

func toSpellCustomData(entry grimoire.GrimoireEntry, variant grimoire.Variant) string {
	return fmt.Sprintf("maf:{spell:{castid:%d,cost:%d,cast:%d,title:%s,description:%s,script:%s}}", entry.CastID, variant.Cost, variant.Cast, jsonString(entry.Title), jsonString(entry.Description), jsonString(entry.Script))
}

func toPrimarySpellCustomData(entry grimoire.GrimoireEntry) string {
	primary := grimoire.Variant{}
	if len(entry.Variants) > 0 {
		primary = entry.Variants[0]
	}
	return fmt.Sprintf("spell:{castid:%d,cost:%d,cast:%d,title:%s,description:%s,script:%s}", entry.CastID, primary.Cost, primary.Cast, jsonString(entry.Title), jsonString(entry.Description), jsonString(entry.Script))
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

func toLoreComponents(value string) []any {
	lines := linesToLoreValues(value)
	out := make([]any, 0, len(lines))
	for _, line := range lines {
		out = append(out, map[string]any{"text": line})
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

func jsonString(value string) string {
	return string(mustJSON(value))
}

func mustJSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return data
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func exportFailure(code, message string, err error) SaveDataResponse {
	return SaveDataResponse{
		OK:      false,
		Code:    code,
		Message: message,
		Details: err.Error(),
	}
}
