package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/idseq"
)

type ItemStateRepository interface {
	LoadItemState() (items.ItemState, error)
	SaveItemState(items.ItemState) error
}

type GrimoireStateRepository interface {
	LoadGrimoireState() (grimoire.GrimoireState, error)
	SaveGrimoireState(grimoire.GrimoireState) error
}

type EntryStateRepository[T any] interface {
	LoadState() (common.EntryState[T], error)
	SaveState(common.EntryState[T]) error
}

type CounterRepository interface {
	LoadCounterState() (idseq.CounterState, error)
	SaveCounterState(idseq.CounterState) error
}

type itemRepository struct {
	path string
}

type grimoireRepository struct {
	path string
}

type entryRepository[T any] struct {
	path string
}

type counterRepository struct {
	path string
}

func NewItemStateRepository(path string) ItemStateRepository {
	return itemRepository{path: path}
}

func NewGrimoireStateRepository(path string) GrimoireStateRepository {
	return grimoireRepository{path: path}
}

func NewEntryStateRepository[T any](path string) EntryStateRepository[T] {
	return entryRepository[T]{path: path}
}

func NewCounterRepository(path string) CounterRepository {
	return counterRepository{path: path}
}

func (r itemRepository) LoadItemState() (items.ItemState, error) {
	var state items.ItemState
	err := readJSON(r.path, &state)
	if err == nil {
		if state.Items == nil {
			state.Items = []items.ItemEntry{}
		}
		return state, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return items.ItemState{Items: []items.ItemEntry{}}, nil
	}
	return items.ItemState{}, err
}

func (r itemRepository) SaveItemState(state items.ItemState) error {
	if state.Items == nil {
		state.Items = []items.ItemEntry{}
	}
	return writeJSON(r.path, state)
}

func (r grimoireRepository) LoadGrimoireState() (grimoire.GrimoireState, error) {
	var state grimoire.GrimoireState
	err := readJSON(r.path, &state)
	if err == nil {
		if state.Entries == nil {
			state.Entries = []grimoire.GrimoireEntry{}
		}
		return state, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return grimoire.GrimoireState{Entries: []grimoire.GrimoireEntry{}}, nil
	}
	return grimoire.GrimoireState{}, err
}

func (r grimoireRepository) SaveGrimoireState(state grimoire.GrimoireState) error {
	if state.Entries == nil {
		state.Entries = []grimoire.GrimoireEntry{}
	}
	return writeJSON(r.path, state)
}

func (r entryRepository[T]) LoadState() (common.EntryState[T], error) {
	var state common.EntryState[T]
	err := readJSON(r.path, &state)
	if err == nil {
		if state.Entries == nil {
			state.Entries = []T{}
		}
		return state, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return common.EntryState[T]{Entries: []T{}}, nil
	}
	return common.EntryState[T]{}, err
}

func (r entryRepository[T]) SaveState(state common.EntryState[T]) error {
	if state.Entries == nil {
		state.Entries = []T{}
	}
	return writeJSON(r.path, state)
}

func (r counterRepository) LoadCounterState() (idseq.CounterState, error) {
	var state idseq.CounterState
	err := readJSON(r.path, &state)
	if err == nil {
		return state, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return idseq.EmptyCounterState(), nil
	}
	return idseq.CounterState{}, err
}

func (r counterRepository) SaveCounterState(state idseq.CounterState) error {
	return writeJSON(r.path, state)
}

func readJSON(path string, dest any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func writeJSON(path string, value any) error {
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
