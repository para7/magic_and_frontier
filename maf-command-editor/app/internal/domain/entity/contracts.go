package entity

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

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

type EntrySaver[T any] interface {
	SaveState(common.EntryState[T]) error
}

type DropRef struct {
	Kind     string   `json:"kind" validate:"trimmed_required,trimmed_oneof=minecraft_item item grimoire"`
	RefID    string   `json:"refId" validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	Weight   float64  `json:"weight" validate:"gte=1,lte=100000"`
	CountMin *float64 `json:"countMin,omitempty" validate:"omitempty,gte=1,lte=64"`
	CountMax *float64 `json:"countMax,omitempty" validate:"omitempty,gte=1,lte=64"`
}

type BaseEntity[E any] struct {
	Mu    *sync.RWMutex
	State *common.EntryState[E]
	Repo  EntrySaver[E]
	IDOf  func(E) string
}

func (b *BaseEntity[E]) Save() error {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	next := CopyEntries(b.State.Entries)
	SortByID(next, b.IDOf)
	b.State.Entries = next
	return b.Repo.SaveState(*b.State)
}

func (b *BaseEntity[E]) ListAll() []E {
	b.Mu.RLock()
	defer b.Mu.RUnlock()
	return CopyEntries(b.State.Entries)
}

func (b *BaseEntity[E]) FindByID(id string) (E, bool) {
	b.Mu.RLock()
	defer b.Mu.RUnlock()
	return FindByID(b.State.Entries, id, b.IDOf)
}

func (b *BaseEntity[E]) HasID(id string) bool {
	_, ok := b.FindByID(id)
	return ok
}

func (b *BaseEntity[E]) CreateEntry(entry E, label string) error {
	if HasID(b.State.Entries, b.IDOf(entry), b.IDOf) {
		return fmt.Errorf("%w: %s %s", ErrDuplicateID, label, b.IDOf(entry))
	}
	return nil
}

func (b *BaseEntity[E]) AppendEntry(entry E) {
	b.State.Entries = append(CopyEntries(b.State.Entries), entry)
}

func (b *BaseEntity[E]) UpdateEntry(entry E, label string) error {
	if !HasID(b.State.Entries, b.IDOf(entry), b.IDOf) {
		return fmt.Errorf("%w: %s %s", ErrNotFound, label, b.IDOf(entry))
	}
	return nil
}

func (b *BaseEntity[E]) UpsertEntry(entry E) {
	next, _ := UpsertByID(b.State.Entries, entry, b.IDOf)
	b.State.Entries = next
}

func (b *BaseEntity[E]) DeleteEntry(id string, label string) error {
	next, ok := DeleteByID(b.State.Entries, id, b.IDOf)
	if !ok {
		return fmt.Errorf("%w: %s %s", ErrNotFound, label, id)
	}
	b.State.Entries = next
	return nil
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
