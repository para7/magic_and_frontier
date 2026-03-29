package export

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type treasureOverlayManifest struct {
	Paths []string `json:"paths"`
}

func treasureOverlayManifestPath(settings ExportSettings) string {
	return filepath.Join(settings.OutputRoot, ".maf-command-editor", "treasure-overrides.json")
}

func cleanupTreasureOverlayOutputs(settings ExportSettings) error {
	manifest, err := readTreasureOverlayManifest(settings)
	if err != nil {
		return err
	}
	for _, rel := range manifest.Paths {
		abs, err := safeOutputPath(settings.OutputRoot, rel)
		if err != nil {
			return err
		}
		if err := os.Remove(abs); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func readTreasureOverlayManifest(settings ExportSettings) (treasureOverlayManifest, error) {
	path := treasureOverlayManifestPath(settings)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return treasureOverlayManifest{}, nil
		}
		return treasureOverlayManifest{}, err
	}
	var manifest treasureOverlayManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return treasureOverlayManifest{}, fmt.Errorf("invalid treasure overlay manifest: %w", err)
	}
	return manifest, nil
}

func writeTreasureOverlayManifest(settings ExportSettings, absPaths []string) error {
	manifest := treasureOverlayManifest{Paths: make([]string, 0, len(absPaths))}
	for _, absPath := range absPaths {
		rel, err := filepath.Rel(settings.OutputRoot, absPath)
		if err != nil {
			return err
		}
		manifest.Paths = append(manifest.Paths, filepath.ToSlash(rel))
	}
	sort.Strings(manifest.Paths)
	path := treasureOverlayManifestPath(settings)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func safeOutputPath(outputRoot, relPath string) (string, error) {
	normalized := filepath.Clean(filepath.FromSlash(relPath))
	if normalized == "." {
		return "", fmt.Errorf("invalid relative path in treasure overlay manifest: %q", relPath)
	}
	joined := filepath.Clean(filepath.Join(outputRoot, normalized))
	root := filepath.Clean(outputRoot)
	if joined == root {
		return "", fmt.Errorf("invalid relative path in treasure overlay manifest: %q", relPath)
	}
	prefix := root + string(os.PathSeparator)
	if !strings.HasPrefix(joined, prefix) {
		return "", fmt.Errorf("path escapes output root in treasure overlay manifest: %q", relPath)
	}
	return joined, nil
}
