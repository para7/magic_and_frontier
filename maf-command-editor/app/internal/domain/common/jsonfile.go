package common

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func ReadJSON(path string, dest any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func WriteJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

type StateRepository[T any] struct {
	Path string
}

func (r StateRepository[T]) LoadState() (EntryState[T], error) {
	var state EntryState[T]
	err := ReadJSON(r.Path, &state)
	if err == nil {
		if state.Entries == nil {
			state.Entries = []T{}
		}
		return state, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return EntryState[T]{Entries: []T{}}, nil
	}
	return EntryState[T]{}, err
}

func (r StateRepository[T]) SaveState(state EntryState[T]) error {
	if state.Entries == nil {
		state.Entries = []T{}
	}
	return WriteJSON(r.Path, state)
}
