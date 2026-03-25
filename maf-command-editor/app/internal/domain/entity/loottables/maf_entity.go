package loottables

import (
	"fmt"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
	"tools2/app/internal/domain/entity/treasures"
)

type entrySaver interface {
	SaveState(common.EntryState[LootTableEntry]) error
}

type EntityDeps struct {
	Mutex       *sync.RWMutex
	State       *common.EntryState[LootTableEntry]
	Repo        entrySaver
	Now         func() time.Time
	ItemIDs     func() map[string]struct{}
	GrimoireIDs func() map[string]struct{}
}

type Entity struct {
	mu          *sync.RWMutex
	state       *common.EntryState[LootTableEntry]
	repo        entrySaver
	now         func() time.Time
	itemIDs     func() map[string]struct{}
	grimoireIDs func() map[string]struct{}
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
	return &Entity{mu: deps.Mutex, state: deps.State, repo: deps.Repo, now: now, itemIDs: itemIDs, grimoireIDs: grimoireIDs}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[LootTableEntry] {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return ValidateSave(input, e.itemIDs(), e.grimoireIDs(), e.now())
}

func (e *Entity) Create(entry LootTableEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if entity.HasID(e.state.Entries, entry.ID, func(it LootTableEntry) string { return it.ID }) {
		return fmt.Errorf("%w: loottable %s", entity.ErrDuplicateID, entry.ID)
	}
	result := ValidateSave(loottableToInput(entry), e.itemIDs(), e.grimoireIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.state.Entries = append(entity.CopyEntries(e.state.Entries), entry)
	return nil
}

func (e *Entity) Update(entry LootTableEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !entity.HasID(e.state.Entries, entry.ID, func(it LootTableEntry) string { return it.ID }) {
		return fmt.Errorf("%w: loottable %s", entity.ErrNotFound, entry.ID)
	}
	result := ValidateSave(loottableToInput(entry), e.itemIDs(), e.grimoireIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	next, _ := entity.UpsertByID(e.state.Entries, entry, func(it LootTableEntry) string { return it.ID })
	e.state.Entries = next
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next, ok := entity.DeleteByID(e.state.Entries, id, func(it LootTableEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: loottable %s", entity.ErrNotFound, id)
	}
	e.state.Entries = next
	return nil
}

func (e *Entity) Save() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next := entity.CopyEntries(e.state.Entries)
	entity.SortByID(next, func(it LootTableEntry) string { return it.ID })
	e.state.Entries = next
	return e.repo.SaveState(*e.state)
}

func (e *Entity) ListAll() []LootTableEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.CopyEntries(e.state.Entries)
}

func (e *Entity) FindByID(id string) (LootTableEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.FindByID(e.state.Entries, id, func(it LootTableEntry) string { return it.ID })
}

func (e *Entity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func loottableToInput(entry LootTableEntry) SaveInput {
	return SaveInput{
		ID:        entry.ID,
		Memo:      entry.Memo,
		LootPools: append([]treasures.DropRef{}, entry.LootPools...),
	}
}

func SaveInputFromEntry(entry LootTableEntry) SaveInput {
	return loottableToInput(entry)
}
