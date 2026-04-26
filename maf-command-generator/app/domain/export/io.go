package export

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var removeAllPath = os.RemoveAll

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

func cleanupConfiguredPaths(outputRoot string, cleanPaths []string) error {
	if len(cleanPaths) == 0 {
		return nil
	}
	resolvedPaths := make([]string, 0, len(cleanPaths))
	for i, cleanPath := range cleanPaths {
		resolvedPath, err := resolveCleanPath(outputRoot, cleanPath)
		if err != nil {
			return fmt.Errorf("cleanPaths[%d]: %w", i, err)
		}
		if err := validateCleanPathResolved(resolvedPath); err != nil {
			return fmt.Errorf("cleanPaths[%d]: %w", i, err)
		}
		resolvedPaths = append(resolvedPaths, resolvedPath)
	}
	for i, resolvedPath := range resolvedPaths {
		if err := removeAllPath(resolvedPath); err != nil {
			return fmt.Errorf("cleanPaths[%d]: remove failed: %w", i, err)
		}
	}
	return nil
}

func validateCleanPath(outputRoot, cleanPath string) error {
	resolvedPath, err := resolveCleanPath(outputRoot, cleanPath)
	if err != nil {
		return err
	}
	return validateCleanPathResolved(resolvedPath)
}

func resolveCleanPath(outputRoot, cleanPath string) (string, error) {
	if outputRoot == "" {
		return "", fmt.Errorf("outputRoot must not be empty")
	}
	rootAbs, err := filepath.Abs(filepath.Clean(outputRoot))
	if err != nil {
		return "", fmt.Errorf("failed to resolve outputRoot path %q: %w", outputRoot, err)
	}
	cleanPathAbs, err := filepath.Abs(filepath.Clean(cleanPath))
	if err != nil {
		return "", fmt.Errorf("failed to resolve cleanPath path %q: %w", cleanPath, err)
	}
	rel, err := filepath.Rel(rootAbs, cleanPathAbs)
	if err != nil {
		return "", fmt.Errorf("failed to compare paths (outputRoot=%q cleanPath=%q): %w", outputRoot, cleanPath, err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("must be inside outputRoot (outputRoot=%q cleanPath=%q)", outputRoot, cleanPath)
	}
	return cleanPathAbs, nil
}

func validateCleanPathResolved(cleanPathAbs string) error {
	info, err := os.Stat(cleanPathAbs)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("target must be directory: %q", cleanPathAbs)
	}

	return filepath.WalkDir(cleanPathAbs, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(d.Name())
		if ext != ".json" && ext != ".mcfunction" {
			return fmt.Errorf("unsupported extension %q in %q", ext, path)
		}
		return nil
	})
}
