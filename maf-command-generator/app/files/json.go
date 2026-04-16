package files

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type entriesFile[T any] struct {
	Entries []T `json:"entries"`
}

// JsonStore[T] reads a directory of {"entries":[...]} JSON files.
// All *.json files in Path are merged into a single slice on Load.
type JsonStore[T any] struct {
	Path string
}

func NewJsonStore[T any](path string) JsonStore[T] {
	return JsonStore[T]{Path: path}
}

// Load reads all *.json files in Path and returns their merged entries.
func (s *JsonStore[T]) Load() ([]T, error) {
	files, err := filepath.Glob(filepath.Join(s.Path, "*.json"))
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
