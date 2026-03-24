package spawntables

import (
	"fmt"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
)

type entrySaver interface {
	SaveState(common.EntryState[SpawnTableEntry]) error
}

type EntityDeps struct {
	Mutex    *sync.RWMutex
	State    *common.EntryState[SpawnTableEntry]
	Repo     entrySaver
	Now      func() time.Time
	EnemyIDs func() map[string]struct{}
}

type Entity struct {
	mu       *sync.RWMutex
	state    *common.EntryState[SpawnTableEntry]
	repo     entrySaver
	now      func() time.Time
	enemyIDs func() map[string]struct{}
}

func NewEntity(deps EntityDeps) *Entity {
	now := deps.Now
	if now == nil {
		now = time.Now
	}
	enemyIDs := deps.EnemyIDs
	if enemyIDs == nil {
		enemyIDs = func() map[string]struct{} { return map[string]struct{}{} }
	}
	return &Entity{mu: deps.Mutex, state: deps.State, repo: deps.Repo, now: now, enemyIDs: enemyIDs}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[SpawnTableEntry] {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.validateLocked(input)
}

func (e *Entity) validateLocked(input SaveInput) common.SaveResult[SpawnTableEntry] {
	result := ValidateSave(input, e.enemyIDs(), e.now())
	if !result.OK {
		return result
	}
	if conflictID, ok := e.firstOverlapLocked(*result.Entry); ok {
		return common.SaveValidationError[SpawnTableEntry](
			common.FieldErrors{"range": "Range overlaps with " + conflictID + "."},
			"Validation failed. Fix the highlighted fields.",
		)
	}
	return result
}

func (e *Entity) firstOverlapLocked(entry SpawnTableEntry) (string, bool) {
	entries := make([]SpawnTableEntry, 0, len(e.state.Entries))
	for _, it := range e.state.Entries {
		if it.ID == entry.ID {
			continue
		}
		entries = append(entries, it)
	}
	return FirstOverlap(entries, entry)
}

func (e *Entity) Create(entry SpawnTableEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if entity.HasID(e.state.Entries, entry.ID, func(it SpawnTableEntry) string { return it.ID }) {
		return fmt.Errorf("%w: spawn table %s", entity.ErrDuplicateID, entry.ID)
	}
	result := e.validateLocked(spawnTableToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.state.Entries = append(entity.CopyEntries(e.state.Entries), entry)
	return nil
}

func (e *Entity) Update(entry SpawnTableEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !entity.HasID(e.state.Entries, entry.ID, func(it SpawnTableEntry) string { return it.ID }) {
		return fmt.Errorf("%w: spawn table %s", entity.ErrNotFound, entry.ID)
	}
	result := e.validateLocked(spawnTableToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	next, _ := entity.UpsertByID(e.state.Entries, entry, func(it SpawnTableEntry) string { return it.ID })
	e.state.Entries = next
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next, ok := entity.DeleteByID(e.state.Entries, id, func(it SpawnTableEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: spawn table %s", entity.ErrNotFound, id)
	}
	e.state.Entries = next
	return nil
}

func (e *Entity) Save() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next := entity.CopyEntries(e.state.Entries)
	entity.SortByID(next, func(it SpawnTableEntry) string { return it.ID })
	e.state.Entries = next
	return e.repo.SaveState(*e.state)
}

func (e *Entity) ListAll() []SpawnTableEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.CopyEntries(e.state.Entries)
}

func (e *Entity) FindByID(id string) (SpawnTableEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.FindByID(e.state.Entries, id, func(it SpawnTableEntry) string { return it.ID })
}

func (e *Entity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func spawnTableToInput(entry SpawnTableEntry) SaveInput {
	return SaveInput{
		ID:            entry.ID,
		SourceMobType: entry.SourceMobType,
		Dimension:     entry.Dimension,
		MinX:          entry.MinX,
		MaxX:          entry.MaxX,
		MinY:          entry.MinY,
		MaxY:          entry.MaxY,
		MinZ:          entry.MinZ,
		MaxZ:          entry.MaxZ,
		BaseMobWeight: entry.BaseMobWeight,
		Replacements:  append([]ReplacementEntry{}, entry.Replacements...),
	}
}

func SaveInputFromEntry(entry SpawnTableEntry) SaveInput {
	return spawnTableToInput(entry)
}
