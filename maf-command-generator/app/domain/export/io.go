package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

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
	body := script
	if !strings.HasSuffix(body, "\n") {
		body += "\n" // mcfunction は末尾改行が必要
	}
	return os.WriteFile(path, []byte(body), 0o755)
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

func removeFileIfExists(path string) error {
	err := os.Remove(path)
	if err == nil || os.IsNotExist(err) {
		return nil
	}
	return err
}
