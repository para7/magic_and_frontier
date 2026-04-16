package export

import (
	"fmt"
	"path/filepath"
	"strings"

	ec "maf_command_editor/app/domain/export/convert"
	model "maf_command_editor/app/domain/model"
	enemyModel "maf_command_editor/app/domain/model/enemy"
	mc "maf_command_editor/app/minecraft"
)

type EnemyArtifact struct {
	ID           string
	SummonScript string
	LootTable    map[string]any
}

func BuildEnemyArtifacts(master DBMaster, enemyLootLogicalDir, minecraftLootRoot string) ([]EnemyArtifact, error) {
	if master == nil {
		return []EnemyArtifact{}, nil
	}
	lookups := buildMasterEntityLookups(master)

	enemies := master.ListEnemies()
	artifacts := make([]EnemyArtifact, 0, len(enemies))
	for _, entry := range enemies {
		pool, err := ec.BuildDropLootPool(entry.Drops, lookups.itemsByID, lookups.grimoiresByID, lookups.passivesByID, lookups.bowsByID, "enemy("+entry.ID+")")
		if err != nil {
			return nil, err
		}

		lootTable, err := buildEnemyLootTable(entry, pool, minecraftLootRoot)
		if err != nil {
			return nil, err
		}

		lootID := resourceRefName("maf", enemyLootLogicalDir, entry.ID)
		lines := ec.ToEnemyFunctionLines(entry, lootID, lookups.itemsByID)
		artifacts = append(artifacts, EnemyArtifact{
			ID:           entry.ID,
			SummonScript: strings.Join(lines, "\n"),
			LootTable:    lootTable,
		})
	}
	return artifacts, nil
}

func WriteEnemyArtifacts(enemyDir string, enemyLootDir string, enemies []EnemyArtifact) error {
	for _, entry := range enemies {
		functionPath := filepath.Join(enemyDir, entry.ID+".mcfunction")
		if err := writeFunctionFile(functionPath, entry.SummonScript); err != nil {
			return err
		}
		lootPath := filepath.Join(enemyLootDir, entry.ID+".json")
		if err := writeJSON(lootPath, entry.LootTable); err != nil {
			return err
		}
	}
	return nil
}

func buildEnemyLootTable(entry enemyModel.Enemy, customPool map[string]any, minecraftLootRoot string) (map[string]any, error) {
	dropMode := strings.TrimSpace(entry.DropMode)
	switch dropMode {
	case "replace":
		return map[string]any{
			"type":  "minecraft:generic",
			"pools": []any{customPool},
		}, nil
	case "append":
		tablePath, err := enemyBaseLootTablePath(entry.MobType)
		if err != nil {
			return nil, err
		}
		base, _, err := mc.LoadLootTable(minecraftLootRoot, tablePath)
		if err != nil {
			return nil, fmt.Errorf("enemy(%s): failed to load base loot table %s: %w", entry.ID, tablePath, err)
		}
		merged, err := ec.MergeLootTablePools(base, customPool, tablePath)
		if err != nil {
			return nil, err
		}
		return merged, nil
	default:
		return nil, fmt.Errorf("enemy(%s): unsupported dropMode %q", entry.ID, entry.DropMode)
	}
}

func enemyBaseLootTablePath(mobType string) (string, error) {
	mobType = strings.TrimSpace(mobType)
	if !model.IsNamespacedResourceID(mobType) {
		return "", fmt.Errorf("invalid mobType: %s", mobType)
	}
	parts := strings.SplitN(mobType, ":", 2)
	if len(parts) != 2 || parts[0] != "minecraft" {
		return "", fmt.Errorf("mobType for append must be minecraft namespace: %s", mobType)
	}
	entityID := model.NormalizeResourcePath(parts[1])
	if entityID == "" {
		return "", fmt.Errorf("invalid mobType: %s", mobType)
	}
	return "minecraft:entities/" + entityID, nil
}
