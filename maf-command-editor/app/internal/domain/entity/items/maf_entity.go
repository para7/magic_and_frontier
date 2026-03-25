package items

import (
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
)

type EntityDeps struct {
	Mutex    *sync.RWMutex
	State    *common.EntryState[ItemEntry]
	Repo     entity.EntrySaver[ItemEntry]
	Now      func() time.Time
	SkillIDs func() map[string]struct{}
}

type Entity struct {
	entity.BaseEntity[ItemEntry]
	now      func() time.Time
	skillIDs func() map[string]struct{}
}

func NewEntity(deps EntityDeps) *Entity {
	now := deps.Now
	if now == nil {
		now = time.Now
	}
	skillIDs := deps.SkillIDs
	if skillIDs == nil {
		skillIDs = func() map[string]struct{} { return map[string]struct{}{} }
	}
	return &Entity{
		BaseEntity: entity.BaseEntity[ItemEntry]{
			Mu: deps.Mutex, State: deps.State, Repo: deps.Repo,
			IDOf: func(e ItemEntry) string { return e.ID },
		},
		now: now, skillIDs: skillIDs,
	}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[ItemEntry] {
	e.Mu.RLock()
	defer e.Mu.RUnlock()
	return ValidateSave(input, e.skillIDs(), e.now())
}

func (e *Entity) Create(entry ItemEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.CreateEntry(entry, "item"); err != nil {
		return err
	}
	result := ValidateSave(itemToInput(entry), e.skillIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.AppendEntry(entry)
	return nil
}

func (e *Entity) Update(entry ItemEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.UpdateEntry(entry, "item"); err != nil {
		return err
	}
	result := ValidateSave(itemToInput(entry), e.skillIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.UpsertEntry(entry)
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	return e.DeleteEntry(id, "item")
}

func itemToInput(entry ItemEntry) SaveInput {
	return SaveInput{
		ID:                  entry.ID,
		ItemID:              entry.ItemID,
		SkillID:             entry.SkillID,
		CustomName:          entry.CustomName,
		Lore:                entry.Lore,
		Enchantments:        entry.Enchantments,
		Unbreakable:         entry.Unbreakable,
		CustomModelData:     entry.CustomModelData,
		RepairCost:          entry.RepairCost,
		HideFlags:           entry.HideFlags,
		PotionID:            entry.PotionID,
		CustomPotionColor:   entry.CustomPotionColor,
		CustomPotionEffects: entry.CustomPotionEffects,
		AttributeModifiers:  entry.AttributeModifiers,
		CustomNBT:           entry.CustomNBT,
	}
}

func SaveInputFromEntry(entry ItemEntry) SaveInput {
	return itemToInput(entry)
}
