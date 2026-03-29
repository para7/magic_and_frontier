package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func normalizeFunctionBody(script string) string {
	if strings.HasSuffix(script, "\n") {
		return script
	}
	return script + "\n"
}

func resourceRefName(namespace, logicalDir, baseName string) string {
	dir := strings.Trim(filepath.ToSlash(logicalDir), "/")
	if dir == "" {
		return namespace + ":" + baseName
	}
	return namespace + ":" + dir + "/" + baseName
}

// minecraft で認識される function の名前を取得する
func functionRefName(logicalDir, baseName string) string {
	return resourceRefName("maf", logicalDir, baseName)
}

func writeFunctionFile(path string, script string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(normalizeFunctionBody(script)), 0o755)
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
