package skills

import (
	"fmt"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
	"tools2/app/internal/domain/entity/items"
)

type EntityDeps struct {
	Mutex      *sync.RWMutex
	State      *common.EntryState[SkillEntry]
	Repo       entity.EntrySaver[SkillEntry]
	Now        func() time.Time
	ItemStates *[]items.ItemEntry
}

type Entity struct {
	entity.BaseEntity[SkillEntry]
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
	return &Entity{
		BaseEntity: entity.BaseEntity[SkillEntry]{
			Mu: deps.Mutex, State: deps.State, Repo: deps.Repo,
			IDOf: func(e SkillEntry) string { return e.ID },
		},
		now: now, itemStates: itemsRef,
	}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[SkillEntry] {
	e.Mu.RLock()
	defer e.Mu.RUnlock()
	return ValidateSave(input, e.now())
}

func (e *Entity) Create(entry SkillEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.CreateEntry(entry, "skill"); err != nil {
		return err
	}
	result := ValidateSave(skillToInput(entry), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.AppendEntry(entry)
	return nil
}

func (e *Entity) Update(entry SkillEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.UpdateEntry(entry, "skill"); err != nil {
		return err
	}
	result := ValidateSave(skillToInput(entry), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.UpsertEntry(entry)
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	trimmedID := entity.TrimID(id)
	for _, item := range *e.itemStates {
		if entity.TrimID(item.SkillID) == trimmedID {
			return fmt.Errorf("%w: skill is referenced by item %s", entity.ErrRelation, item.ID)
		}
	}
	return e.DeleteEntry(id, "skill")
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
