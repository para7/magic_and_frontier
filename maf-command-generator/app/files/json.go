package files

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type entriesFile[T any] struct {
	Entries []T `json:"entries"`
}

type JsonStore[T any] struct {
	Path    string
	Entries []T
}

func NewJsonStore[T any](path string) JsonStore[T] {
	return JsonStore[T]{Path: path}
}

func (s *JsonStore[T]) Load() error {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}
	var f entriesFile[T]
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	s.Entries = f.Entries
	return nil
}

func (s *JsonStore[T]) Save() error {
	entries := s.Entries
	if entries == nil {
		entries = []T{}
	}
	data, err := json.MarshalIndent(entriesFile[T]{Entries: entries}, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(s.Path, data, 0o644)
}
