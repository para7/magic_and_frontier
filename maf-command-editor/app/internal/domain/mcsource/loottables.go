package mcsource

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"tools2/app/internal/domain/common"
)

type LootTableSource struct {
	TablePath string
	FilePath  string
}

func ListLootTables(root string) ([]LootTableSource, error) {
	root = filepath.Clean(root)
	entries := []LootTableSource{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		tablePath, convErr := tablePathFromFile(root, path)
		if convErr != nil {
			return convErr
		}
		entries = append(entries, LootTableSource{
			TablePath: tablePath,
			FilePath:  path,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].TablePath < entries[j].TablePath
	})
	return entries, nil
}

func LoadLootTable(root, tablePath string) (map[string]any, string, error) {
	filePath, err := FilePathForTable(root, tablePath)
	if err != nil {
		return nil, "", err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, filePath, err
	}
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, filePath, fmt.Errorf("read minecraft loot table %s: %w", filePath, err)
	}
	return parsed, filePath, nil
}

func Exists(root, tablePath string) (bool, string, error) {
	filePath, err := FilePathForTable(root, tablePath)
	if err != nil {
		return false, "", err
	}
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false, filePath, nil
		}
		return false, filePath, err
	}
	return true, filePath, nil
}

func FilePathForTable(root, tablePath string) (string, error) {
	if !common.IsSafeNamespacedResourcePath(tablePath) {
		return "", fmt.Errorf("invalid loot table path: %s", tablePath)
	}
	parts := strings.SplitN(tablePath, ":", 2)
	if len(parts) != 2 || parts[0] != "minecraft" || parts[1] == "" {
		return "", fmt.Errorf("treasure tablePath must target minecraft namespace: %s", tablePath)
	}
	return filepath.Join(filepath.Clean(root), filepath.FromSlash(common.NormalizeResourcePath(parts[1]))+".json"), nil
}

func tablePathFromFile(root, path string) (string, error) {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}
	rel = filepath.ToSlash(rel)
	if !strings.HasSuffix(rel, ".json") {
		return "", fmt.Errorf("expected json file: %s", path)
	}
	rel = strings.TrimSuffix(rel, ".json")
	rel = common.NormalizeResourcePath(rel)
	if rel == "" {
		return "", fmt.Errorf("invalid minecraft loot table file: %s", path)
	}
	return "minecraft:" + rel, nil
}
