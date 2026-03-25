package export

import (
	"fmt"
	"os"
	"path/filepath"

	"maf-command-editor/app/internal/domain/entity/grimoire"
	"maf-command-editor/app/internal/domain/entity/items"
	"maf-command-editor/app/internal/domain/entity/loottables"
	"maf-command-editor/app/internal/domain/entity/treasures"
	"maf-command-editor/app/internal/domain/mcsource"
)

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
