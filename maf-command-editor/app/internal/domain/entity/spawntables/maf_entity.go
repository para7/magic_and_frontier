package spawntables

import (
	"sync"
	"time"

	"maf-command-editor/app/internal/domain/common"
	"maf-command-editor/app/internal/domain/entity"
)

type EntityDeps struct {
	Mutex    *sync.RWMutex
	State    *common.EntryState[SpawnTableEntry]
	Repo     entity.EntrySaver[SpawnTableEntry]
	Now      func() time.Time
	EnemyIDs func() map[string]struct{}
}

type Entity struct {
	entity.BaseEntity[SpawnTableEntry]
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
	return &Entity{
		BaseEntity: entity.BaseEntity[SpawnTableEntry]{
			Mu: deps.Mutex, State: deps.State, Repo: deps.Repo,
			IDOf: func(e SpawnTableEntry) string { return e.ID },
		},
		now: now, enemyIDs: enemyIDs,
	}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[SpawnTableEntry] {
	e.Mu.RLock()
	defer e.Mu.RUnlock()
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
	entries := make([]SpawnTableEntry, 0, len(e.State.Entries))
	for _, it := range e.State.Entries {
		if it.ID == entry.ID {
			continue
		}
		entries = append(entries, it)
	}
	return FirstOverlap(entries, entry)
}

func (e *Entity) Create(entry SpawnTableEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.CreateEntry(entry, "spawn table"); err != nil {
		return err
	}
	result := e.validateLocked(spawnTableToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.AppendEntry(entry)
	return nil
}

func (e *Entity) Update(entry SpawnTableEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.UpdateEntry(entry, "spawn table"); err != nil {
		return err
	}
	result := e.validateLocked(spawnTableToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.UpsertEntry(entry)
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	return e.DeleteEntry(id, "spawn table")
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
