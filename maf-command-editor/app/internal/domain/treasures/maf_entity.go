package treasures

import (
	"fmt"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
)

type entrySaver interface {
	SaveState(common.EntryState[TreasureEntry]) error
}

type EntityDeps struct {
	Mutex               *sync.RWMutex
	State               *common.EntryState[TreasureEntry]
	Repo                entrySaver
	Now                 func() time.Time
	ItemIDs             func() map[string]struct{}
	GrimoireIDs         func() map[string]struct{}
	TreasureSourceErr   func() error
	TreasureSourcePaths func() map[string]struct{}
}

type Entity struct {
	mu                  *sync.RWMutex
	state               *common.EntryState[TreasureEntry]
	repo                entrySaver
	now                 func() time.Time
	itemIDs             func() map[string]struct{}
	grimoireIDs         func() map[string]struct{}
	treasureSourceErr   func() error
	treasureSourcePaths func() map[string]struct{}
}

func NewEntity(deps EntityDeps) *Entity {
	now := deps.Now
	if now == nil {
		now = time.Now
	}
	emptySet := func() map[string]struct{} { return map[string]struct{}{} }
	itemIDs := deps.ItemIDs
	if itemIDs == nil {
		itemIDs = emptySet
	}
	grimoireIDs := deps.GrimoireIDs
	if grimoireIDs == nil {
		grimoireIDs = emptySet
	}
	sourcePaths := deps.TreasureSourcePaths
	if sourcePaths == nil {
		sourcePaths = emptySet
	}
	sourceErr := deps.TreasureSourceErr
	if sourceErr == nil {
		sourceErr = func() error { return nil }
	}
	return &Entity{
		mu: deps.Mutex, state: deps.State, repo: deps.Repo, now: now,
		itemIDs: itemIDs, grimoireIDs: grimoireIDs,
		treasureSourceErr: sourceErr, treasureSourcePaths: sourcePaths,
	}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[TreasureEntry] {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.validateLocked(input)
}

func (e *Entity) validateLocked(input SaveInput) common.SaveResult[TreasureEntry] {
	if err := e.treasureSourceErr(); err != nil {
		return common.SaveValidationError[TreasureEntry](
			common.FieldErrors{"tablePath": err.Error()},
			"Validation failed. Fix the highlighted fields.",
		)
	}
	result := ValidateSave(input, e.itemIDs(), e.grimoireIDs(), e.treasureSourcePaths(), e.now())
	if !result.OK {
		return result
	}
	if conflictID := e.duplicateTablePathLocked(result.Entry.ID, result.Entry.TablePath); conflictID != "" {
		return common.SaveValidationError[TreasureEntry](
			common.FieldErrors{"tablePath": "Loot table path is already used by " + conflictID + "."},
			"Validation failed. Fix the highlighted fields.",
		)
	}
	return result
}

func (e *Entity) duplicateTablePathLocked(entryID, tablePath string) string {
	tablePath = entity.TrimID(tablePath)
	for _, entry := range e.state.Entries {
		if entry.ID != entryID && entity.TrimID(entry.TablePath) == tablePath {
			return entry.ID
		}
	}
	return ""
}

func (e *Entity) Create(entry TreasureEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if entity.HasID(e.state.Entries, entry.ID, func(it TreasureEntry) string { return it.ID }) {
		return fmt.Errorf("%w: treasure %s", entity.ErrDuplicateID, entry.ID)
	}
	result := e.validateLocked(treasureToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.state.Entries = append(entity.CopyEntries(e.state.Entries), entry)
	return nil
}

func (e *Entity) Update(entry TreasureEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !entity.HasID(e.state.Entries, entry.ID, func(it TreasureEntry) string { return it.ID }) {
		return fmt.Errorf("%w: treasure %s", entity.ErrNotFound, entry.ID)
	}
	result := e.validateLocked(treasureToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	next, _ := entity.UpsertByID(e.state.Entries, entry, func(it TreasureEntry) string { return it.ID })
	e.state.Entries = next
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next, ok := entity.DeleteByID(e.state.Entries, id, func(it TreasureEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: treasure %s", entity.ErrNotFound, id)
	}
	e.state.Entries = next
	return nil
}

func (e *Entity) Save() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next := entity.CopyEntries(e.state.Entries)
	entity.SortByID(next, func(it TreasureEntry) string { return it.ID })
	e.state.Entries = next
	return e.repo.SaveState(*e.state)
}

func (e *Entity) ListAll() []TreasureEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.CopyEntries(e.state.Entries)
}

func (e *Entity) FindByID(id string) (TreasureEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.FindByID(e.state.Entries, id, func(it TreasureEntry) string { return it.ID })
}

func (e *Entity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func treasureToInput(entry TreasureEntry) SaveInput {
	return SaveInput{
		ID:        entry.ID,
		TablePath: entry.TablePath,
		LootPools: append([]DropRef{}, entry.LootPools...),
	}
}

func SaveInputFromEntry(entry TreasureEntry) SaveInput {
	return treasureToInput(entry)
}
