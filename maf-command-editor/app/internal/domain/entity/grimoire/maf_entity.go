package grimoire

import (
	"sync"
	"time"

	"maf-command-editor/app/internal/domain/common"
	"maf-command-editor/app/internal/domain/entity"
)

type EntityDeps struct {
	Mutex *sync.RWMutex
	State *common.EntryState[GrimoireEntry]
	Repo  entity.EntrySaver[GrimoireEntry]
	Now   func() time.Time
}

type Entity struct {
	entity.BaseEntity[GrimoireEntry]
	now func() time.Time
}

func NewEntity(deps EntityDeps) *Entity {
	now := deps.Now
	if now == nil {
		now = time.Now
	}
	return &Entity{
		BaseEntity: entity.BaseEntity[GrimoireEntry]{
			Mu: deps.Mutex, State: deps.State, Repo: deps.Repo,
			IDOf: func(e GrimoireEntry) string { return e.ID },
		},
		now: now,
	}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[GrimoireEntry] {
	e.Mu.RLock()
	defer e.Mu.RUnlock()
	return e.validateLocked(input)
}

func (e *Entity) validateLocked(input SaveInput) common.SaveResult[GrimoireEntry] {
	result := ValidateSave(input, e.now())
	if !result.OK {
		return result
	}
	if conflictID := e.duplicateCastIDLocked(result.Entry.ID, result.Entry.CastID); conflictID != "" {
		return common.SaveValidationError[GrimoireEntry](
			common.FieldErrors{"castid": "Cast ID is already used by " + conflictID + "."},
			"Validation failed. Fix the highlighted fields.",
		)
	}
	return result
}

func (e *Entity) duplicateCastIDLocked(entryID string, castID int) string {
	for _, entry := range e.State.Entries {
		if entry.ID != entryID && entry.CastID == castID {
			return entry.ID
		}
	}
	return ""
}

func (e *Entity) Create(entry GrimoireEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.CreateEntry(entry, "grimoire"); err != nil {
		return err
	}
	result := e.validateLocked(grimoireToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.AppendEntry(entry)
	return nil
}

func (e *Entity) Update(entry GrimoireEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.UpdateEntry(entry, "grimoire"); err != nil {
		return err
	}
	result := e.validateLocked(grimoireToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.UpsertEntry(entry)
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	return e.DeleteEntry(id, "grimoire")
}

func grimoireToInput(entry GrimoireEntry) SaveInput {
	return SaveInput{
		ID:          entry.ID,
		CastID:      entry.CastID,
		CastTime:    entry.CastTime,
		MPCost:      entry.MPCost,
		Script:      entry.Script,
		Title:       entry.Title,
		Description: entry.Description,
	}
}

func SaveInputFromEntry(entry GrimoireEntry) SaveInput {
	return grimoireToInput(entry)
}
