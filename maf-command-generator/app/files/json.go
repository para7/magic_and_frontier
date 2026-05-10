package files

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type entriesFile[T any] struct {
	Entries []T `json:"entries"`
}

// JsonStore[T] reads a directory tree of {"entries":[...]} JSON files.
// All *.json files under Path are merged into a single slice on Load.
type JsonStore[T any] struct {
	Path string
}

func NewJsonStore[T any](path string) JsonStore[T] {
	return JsonStore[T]{Path: path}
}

// Load reads all *.json files under Path and returns their merged entries.
func (s *JsonStore[T]) Load() ([]T, error) {
	files, err := JSONFilePaths(s.Path)
	if err != nil {
		return nil, err
	}
	var merged []T
	for _, fpath := range files {
		data, err := os.ReadFile(fpath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", filepath.Base(fpath), err)
		}
		var f entriesFile[T]
		if err := json.Unmarshal(data, &f); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", filepath.Base(fpath), err)
		}
		merged = append(merged, f.Entries...)
	}
	return merged, nil
}

func JSONFilePaths(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return paths, nil
}
