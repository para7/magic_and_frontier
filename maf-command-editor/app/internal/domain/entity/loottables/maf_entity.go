package loottables

import (
	"sync"
	"time"

	"maf-command-editor/app/internal/domain/common"
	"maf-command-editor/app/internal/domain/entity"
	"maf-command-editor/app/internal/domain/entity/treasures"
)

type EntityDeps struct {
	Mutex       *sync.RWMutex
	State       *common.EntryState[LootTableEntry]
	Repo        entity.EntrySaver[LootTableEntry]
	Now         func() time.Time
	ItemIDs     func() map[string]struct{}
	GrimoireIDs func() map[string]struct{}
}

type Entity struct {
	entity.BaseEntity[LootTableEntry]
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
	return &Entity{
		BaseEntity: entity.BaseEntity[LootTableEntry]{
			Mu: deps.Mutex, State: deps.State, Repo: deps.Repo,
			IDOf: func(e LootTableEntry) string { return e.ID },
		},
		now: now, itemIDs: itemIDs, grimoireIDs: grimoireIDs,
	}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[LootTableEntry] {
	e.Mu.RLock()
	defer e.Mu.RUnlock()
	return ValidateSave(input, e.itemIDs(), e.grimoireIDs(), e.now())
}

func (e *Entity) Create(entry LootTableEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.CreateEntry(entry, "loottable"); err != nil {
		return err
	}
	result := ValidateSave(loottableToInput(entry), e.itemIDs(), e.grimoireIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.AppendEntry(entry)
	return nil
}

func (e *Entity) Update(entry LootTableEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.UpdateEntry(entry, "loottable"); err != nil {
		return err
	}
	result := ValidateSave(loottableToInput(entry), e.itemIDs(), e.grimoireIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.UpsertEntry(entry)
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	return e.DeleteEntry(id, "loottable")
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
