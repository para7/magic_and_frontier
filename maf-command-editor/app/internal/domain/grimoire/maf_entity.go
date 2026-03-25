package grimoire

import (
	"fmt"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
)

type entrySaver interface {
	SaveState(common.EntryState[GrimoireEntry]) error
}

type EntityDeps struct {
	Mutex *sync.RWMutex
	State *common.EntryState[GrimoireEntry]
	Repo  entrySaver
	Now   func() time.Time
}

type Entity struct {
	mu    *sync.RWMutex
	state *common.EntryState[GrimoireEntry]
	repo  entrySaver
	now   func() time.Time
}

func NewEntity(deps EntityDeps) *Entity {
	now := deps.Now
	if now == nil {
		now = time.Now
	}
	return &Entity{mu: deps.Mutex, state: deps.State, repo: deps.Repo, now: now}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[GrimoireEntry] {
	e.mu.RLock()
	defer e.mu.RUnlock()
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
	for _, entry := range e.state.Entries {
		if entry.ID != entryID && entry.CastID == castID {
			return entry.ID
		}
	}
	return ""
}

func (e *Entity) Create(entry GrimoireEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if entity.HasID(e.state.Entries, entry.ID, func(it GrimoireEntry) string { return it.ID }) {
		return fmt.Errorf("%w: grimoire %s", entity.ErrDuplicateID, entry.ID)
	}
	result := e.validateLocked(grimoireToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.state.Entries = append(entity.CopyEntries(e.state.Entries), entry)
	return nil
}

func (e *Entity) Update(entry GrimoireEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !entity.HasID(e.state.Entries, entry.ID, func(it GrimoireEntry) string { return it.ID }) {
		return fmt.Errorf("%w: grimoire %s", entity.ErrNotFound, entry.ID)
	}
	result := e.validateLocked(grimoireToInput(entry))
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	next, _ := entity.UpsertByID(e.state.Entries, entry, func(it GrimoireEntry) string { return it.ID })
	e.state.Entries = next
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next, ok := entity.DeleteByID(e.state.Entries, id, func(it GrimoireEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: grimoire %s", entity.ErrNotFound, id)
	}
	e.state.Entries = next
	return nil
}

func (e *Entity) Save() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next := entity.CopyEntries(e.state.Entries)
	entity.SortByID(next, func(it GrimoireEntry) string { return it.ID })
	e.state.Entries = next
	return e.repo.SaveState(*e.state)
}

func (e *Entity) ListAll() []GrimoireEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.CopyEntries(e.state.Entries)
}

func (e *Entity) FindByID(id string) (GrimoireEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.FindByID(e.state.Entries, id, func(it GrimoireEntry) string { return it.ID })
}

func (e *Entity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
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
