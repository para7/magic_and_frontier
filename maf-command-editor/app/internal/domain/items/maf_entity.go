package items

import (
	"fmt"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
)

type entrySaver interface {
	SaveState(common.EntryState[ItemEntry]) error
}

type EntityDeps struct {
	Mutex    *sync.RWMutex
	State    *common.EntryState[ItemEntry]
	Repo     entrySaver
	Now      func() time.Time
	SkillIDs func() map[string]struct{}
}

type Entity struct {
	mu       *sync.RWMutex
	state    *common.EntryState[ItemEntry]
	repo     entrySaver
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
	return &Entity{mu: deps.Mutex, state: deps.State, repo: deps.Repo, now: now, skillIDs: skillIDs}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[ItemEntry] {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return ValidateSave(input, e.skillIDs(), e.now())
}

func (e *Entity) Create(entry ItemEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if entity.HasID(e.state.Entries, entry.ID, func(it ItemEntry) string { return it.ID }) {
		return fmt.Errorf("%w: item %s", entity.ErrDuplicateID, entry.ID)
	}
	result := ValidateSave(itemToInput(entry), e.skillIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.state.Entries = append(entity.CopyEntries(e.state.Entries), entry)
	return nil
}

func (e *Entity) Update(entry ItemEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !entity.HasID(e.state.Entries, entry.ID, func(it ItemEntry) string { return it.ID }) {
		return fmt.Errorf("%w: item %s", entity.ErrNotFound, entry.ID)
	}
	result := ValidateSave(itemToInput(entry), e.skillIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	next, _ := entity.UpsertByID(e.state.Entries, entry, func(it ItemEntry) string { return it.ID })
	e.state.Entries = next
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next, ok := entity.DeleteByID(e.state.Entries, id, func(it ItemEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: item %s", entity.ErrNotFound, id)
	}
	e.state.Entries = next
	return nil
}

func (e *Entity) Save() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next := entity.CopyEntries(e.state.Entries)
	entity.SortByID(next, func(it ItemEntry) string { return it.ID })
	e.state.Entries = next
	return e.repo.SaveState(*e.state)
}

func (e *Entity) ListAll() []ItemEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.CopyEntries(e.state.Entries)
}

func (e *Entity) FindByID(id string) (ItemEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.FindByID(e.state.Entries, id, func(it ItemEntry) string { return it.ID })
}

func (e *Entity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
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
