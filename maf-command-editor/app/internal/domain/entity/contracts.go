package entity

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"tools2/app/internal/domain/common"
)

var (
	ErrDuplicateID = errors.New("duplicate id")
	ErrNotFound    = errors.New("entry not found")
	ErrRelation    = errors.New("relation validation failed")
)

type MasterRef interface {
	HasItem(id string) bool
	HasGrimoire(id string) bool
	HasSkill(id string) bool
	HasEnemySkill(id string) bool
	HasEnemy(id string) bool
	HasTreasure(id string) bool
	HasLootTable(id string) bool
	HasSpawnTable(id string) bool
}

type MafEntity[I any, E any] interface {
	Validate(input I, master MasterRef) common.SaveResult[E]
	Create(entry E, master MasterRef) error
	Update(entry E, master MasterRef) error
	Delete(id string, master MasterRef) error
	Save() error
	ListAll() []E
	FindByID(id string) (E, bool)
	HasID(id string) bool
}

func TrimID(value string) string {
	return strings.TrimSpace(value)
}

func CopyEntries[T any](entries []T) []T {
	return append([]T{}, entries...)
}

func HasID[T any](entries []T, id string, idOf func(T) string) bool {
	id = TrimID(id)
	if id == "" {
		return false
	}
	for _, entry := range entries {
		if TrimID(idOf(entry)) == id {
			return true
		}
	}
	return false
}

func FindByID[T any](entries []T, id string, idOf func(T) string) (T, bool) {
	var zero T
	id = TrimID(id)
	if id == "" {
		return zero, false
	}
	for _, entry := range entries {
		if TrimID(idOf(entry)) == id {
			return entry, true
		}
	}
	return zero, false
}

func UpsertByID[T any](entries []T, entry T, idOf func(T) string) ([]T, bool) {
	entryID := TrimID(idOf(entry))
	next := CopyEntries(entries)
	for i := range next {
		if TrimID(idOf(next[i])) == entryID {
			next[i] = entry
			return next, true
		}
	}
	next = append(next, entry)
	return next, false
}

func DeleteByID[T any](entries []T, id string, idOf func(T) string) ([]T, bool) {
	id = TrimID(id)
	next := make([]T, 0, len(entries))
	found := false
	for _, entry := range entries {
		if TrimID(idOf(entry)) == id {
			found = true
			continue
		}
		next = append(next, entry)
	}
	return next, found
}

func SortByID[T any](entries []T, idOf func(T) string) {
	sort.Slice(entries, func(i, j int) bool {
		return TrimID(idOf(entries[i])) < TrimID(idOf(entries[j]))
	})
}

func IDSet[T any](entries []T, idOf func(T) string) map[string]struct{} {
	out := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		id := TrimID(idOf(entry))
		if id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}

func RelationErrFromResult[T any](result common.SaveResult[T]) error {
	if len(result.FieldErrors) > 0 {
		return fmt.Errorf("%w: %s", ErrRelation, firstFieldError(result.FieldErrors))
	}
	if strings.TrimSpace(result.FormError) != "" {
		return fmt.Errorf("%w: %s", ErrRelation, result.FormError)
	}
	return fmt.Errorf("%w", ErrRelation)
}

func firstFieldError(errs common.FieldErrors) string {
	if len(errs) == 0 {
		return "validation failed"
	}
	keys := make([]string, 0, len(errs))
	for key := range errs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	first := keys[0]
	return first + ": " + errs[first]
}
