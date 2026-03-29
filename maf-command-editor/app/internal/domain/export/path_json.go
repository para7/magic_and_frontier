package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"maf-command-editor/app/internal/domain/common"
)

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
