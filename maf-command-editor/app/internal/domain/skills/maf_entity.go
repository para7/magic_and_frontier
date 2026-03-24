package skills

import (
	"fmt"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
	"tools2/app/internal/domain/items"
)

type entrySaver interface {
	SaveState(common.EntryState[SkillEntry]) error
}

type EntityDeps struct {
	Mutex      *sync.RWMutex
	State      *common.EntryState[SkillEntry]
	Repo       entrySaver
	Now        func() time.Time
	ItemStates *[]items.ItemEntry
}

type Entity struct {
	mu         *sync.RWMutex
	state      *common.EntryState[SkillEntry]
	repo       entrySaver
	now        func() time.Time
	itemStates *[]items.ItemEntry
}

func NewEntity(deps EntityDeps) *Entity {
	now := deps.Now
	if now == nil {
		now = time.Now
	}
	itemsRef := deps.ItemStates
	if itemsRef == nil {
		empty := []items.ItemEntry{}
		itemsRef = &empty
	}
	return &Entity{mu: deps.Mutex, state: deps.State, repo: deps.Repo, now: now, itemStates: itemsRef}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[SkillEntry] {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return ValidateSave(input, e.now())
}

func (e *Entity) Create(entry SkillEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if entity.HasID(e.state.Entries, entry.ID, func(it SkillEntry) string { return it.ID }) {
		return fmt.Errorf("%w: skill %s", entity.ErrDuplicateID, entry.ID)
	}
	result := ValidateSave(skillToInput(entry), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.state.Entries = append(entity.CopyEntries(e.state.Entries), entry)
	return nil
}

func (e *Entity) Update(entry SkillEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !entity.HasID(e.state.Entries, entry.ID, func(it SkillEntry) string { return it.ID }) {
		return fmt.Errorf("%w: skill %s", entity.ErrNotFound, entry.ID)
	}
	result := ValidateSave(skillToInput(entry), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	next, _ := entity.UpsertByID(e.state.Entries, entry, func(it SkillEntry) string { return it.ID })
	e.state.Entries = next
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	trimmedID := entity.TrimID(id)
	for _, item := range *e.itemStates {
		if entity.TrimID(item.SkillID) == trimmedID {
			return fmt.Errorf("%w: skill is referenced by item %s", entity.ErrRelation, item.ID)
		}
	}
	next, ok := entity.DeleteByID(e.state.Entries, id, func(it SkillEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: skill %s", entity.ErrNotFound, id)
	}
	e.state.Entries = next
	return nil
}

func (e *Entity) Save() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next := entity.CopyEntries(e.state.Entries)
	entity.SortByID(next, func(it SkillEntry) string { return it.ID })
	e.state.Entries = next
	return e.repo.SaveState(*e.state)
}

func (e *Entity) ListAll() []SkillEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.CopyEntries(e.state.Entries)
}

func (e *Entity) FindByID(id string) (SkillEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.FindByID(e.state.Entries, id, func(it SkillEntry) string { return it.ID })
}

func (e *Entity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func skillToInput(entry SkillEntry) SaveInput {
	return SaveInput{
		ID:          entry.ID,
		Name:        entry.Name,
		SkillType:   entry.SkillType,
		Description: entry.Description,
		Script:      entry.Script,
	}
}

func SaveInputFromEntry(entry SkillEntry) SaveInput {
	return skillToInput(entry)
}
