package treasures

import (
	"sync"
	"time"

	"maf-command-editor/app/internal/domain/common"
	"maf-command-editor/app/internal/domain/entity"
)

type EntityDeps struct {
	Mutex               *sync.RWMutex
	State               *common.EntryState[TreasureEntry]
	Repo                entity.EntrySaver[TreasureEntry]
	Now                 func() time.Time
	ItemIDs             func() map[string]struct{}
	GrimoireIDs         func() map[string]struct{}
	TreasureSourceErr   func() error
	TreasureSourcePaths func() map[string]struct{}
}

type Entity struct {
	entity.BaseEntity[TreasureEntry]
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
		BaseEntity: entity.BaseEntity[TreasureEntry]{
			Mu: deps.Mutex, State: deps.State, Repo: deps.Repo,
			IDOf: func(e TreasureEntry) string { return e.ID },
		},
		now: now, itemIDs: itemIDs, grimoireIDs: grimoireIDs,
		treasureSourceErr: sourceErr, treasureSourcePaths: sourcePaths,
	}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[TreasureEntry] {
	e.Mu.RLock()
	defer e.Mu.RUnlock()
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
	for _, entry := range e.State.Entries {
		if entry.ID != entryID && entity.TrimID(entry.TablePath) == tablePath {
			return entry.ID
		}
	}
	return ""
}

func (e *Entity) Create(entry TreasureEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.CreateEntry(entry, "treasure"); err != nil {
		return err
	}
	result := e.validateLocked(treasureToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.AppendEntry(entry)
	return nil
}

func (e *Entity) Update(entry TreasureEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.UpdateEntry(entry, "treasure"); err != nil {
		return err
	}
	result := e.validateLocked(treasureToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.UpsertEntry(entry)
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	return e.DeleteEntry(id, "treasure")
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
