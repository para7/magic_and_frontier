package export

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	ec "maf_command_editor/app/domain/export/convert"
	mc "maf_command_editor/app/minecraft"
)

type TreasureArtifact struct {
	Namespace string
	RelPath   string
	LootTable map[string]any
}

func BuildTreasureArtifacts(master DBMaster, lootTableSourceRoot, minecraftLootRoot string) ([]TreasureArtifact, error) {
	if master == nil {
		return []TreasureArtifact{}, nil
	}

	lootTableSourceRoot = filepath.Clean(strings.TrimSpace(lootTableSourceRoot))
	if lootTableSourceRoot == "" {
		return []TreasureArtifact{}, nil
	}

	if _, err := os.Stat(lootTableSourceRoot); err != nil {
		if os.IsNotExist(err) {
			return []TreasureArtifact{}, nil
		}
		return nil, err
	}
	lookups := buildMasterEntityLookups(master)

	artifacts := []TreasureArtifact{}
	err := filepath.WalkDir(lootTableSourceRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		namespace, relPath, err := namespaceAndRelPath(lootTableSourceRoot, path)
		if err != nil {
			return err
		}
		context := fmt.Sprintf("loot_table(%s:%s)", namespace, relPath)

		lootTable, err := readLootTableJSON(path)
		if err != nil {
			return fmt.Errorf("%s: %w", context, err)
		}

		rawPools, hasPools := lootTable["pools"]
		if !hasPools {
			return nil
		}
		pools, ok := rawPools.([]any)
		if !ok {
			return fmt.Errorf("%s: pools must be an array", context)
		}

		resolvedPools, err := ec.ResolveMafLootPools(pools, lookups.itemsByID, lookups.grimoiresByID, lookups.passivesByID, lookups.bowsByID, context)
		if err != nil {
			return err
		}

		var out map[string]any
		if namespace == "minecraft" {
			tablePath := namespace + ":" + relPath
			base, _, err := mc.LoadLootTable(minecraftLootRoot, tablePath)
			if err != nil {
				return fmt.Errorf("%s: failed to load base loot table %s: %w", context, tablePath, err)
			}
			merged := base
			for i, rawPool := range resolvedPools {
				pool, ok := rawPool.(map[string]any)
				if !ok {
					return fmt.Errorf("%s: resolved pools[%d] must be an object", context, i)
				}
				merged, err = ec.MergeLootTablePools(merged, pool, tablePath)
				if err != nil {
					return err
				}
			}
			out = merged
		} else {
			out = make(map[string]any, len(lootTable))
			for key, value := range lootTable {
				out[key] = value
			}
			out["pools"] = resolvedPools
		}

		artifacts = append(artifacts, TreasureArtifact{
			Namespace: namespace,
			RelPath:   relPath,
			LootTable: out,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(artifacts, func(i, j int) bool {
		if artifacts[i].Namespace == artifacts[j].Namespace {
			return artifacts[i].RelPath < artifacts[j].RelPath
		}
		return artifacts[i].Namespace < artifacts[j].Namespace
	})

	return artifacts, nil
}

func WriteTreasureArtifacts(outputRoot string, artifacts []TreasureArtifact) error {
	for _, entry := range artifacts {
		outputPath := filepath.Join(
			outputRoot,
			"data",
			entry.Namespace,
			"loot_table",
			filepath.FromSlash(entry.RelPath)+".json",
		)
		if err := writeJSON(outputPath, entry.LootTable); err != nil {
			return err
		}
	}
	return nil
}

func namespaceAndRelPath(root, filePath string) (string, string, error) {
	rel, err := filepath.Rel(root, filePath)
	if err != nil {
		return "", "", err
	}
	rel = filepath.ToSlash(rel)
	if strings.HasPrefix(rel, "../") {
		return "", "", fmt.Errorf("loot table source is outside root: %s", filePath)
	}
	if !strings.HasSuffix(rel, ".json") {
		return "", "", fmt.Errorf("loot table source must be json: %s", filePath)
	}
	rel = strings.TrimSuffix(rel, ".json")

	parts := strings.SplitN(rel, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("loot table source must be namespace/path.json: %s", filePath)
	}

	namespace := strings.TrimSpace(parts[0])
	relPath := strings.Trim(parts[1], "/")
	if namespace == "" || relPath == "" {
		return "", "", fmt.Errorf("loot table source must be namespace/path.json: %s", filePath)
	}
	return namespace, relPath, nil
}

func readLootTableJSON(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var table map[string]any
	if err := json.Unmarshal(data, &table); err != nil {
		return nil, fmt.Errorf("read loot table %s: %w", path, err)
	}
	return table, nil
}
