package files

import (
	"encoding/json"
	"os"
)

type entriesFile[T any] struct {
	Entries []T `json:"entries"`
}

type JsonStore[T any] struct {
	Path string
	// Entries []T
}

func NewJsonStore[T any](path string) JsonStore[T] {
	return JsonStore[T]{Path: path}
}

func (s *JsonStore[T]) Load() ([]T, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, err
	}
	var f entriesFile[T]
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	// s.Entries = f.Entries
	// return nil
	return f.Entries, nil
}
