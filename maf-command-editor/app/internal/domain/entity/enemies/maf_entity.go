package enemies

import (
	"fmt"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
)

type entrySaver interface {
	SaveState(common.EntryState[EnemyEntry]) error
}

type EntityDeps struct {
	Mutex         *sync.RWMutex
	State         *common.EntryState[EnemyEntry]
	Repo          entrySaver
	Now           func() time.Time
	EnemySkillIDs func() map[string]struct{}
	ItemIDs       func() map[string]struct{}
	GrimoireIDs   func() map[string]struct{}
}

type Entity struct {
	mu            *sync.RWMutex
	state         *common.EntryState[EnemyEntry]
	repo          entrySaver
	now           func() time.Time
	enemySkillIDs func() map[string]struct{}
	itemIDs       func() map[string]struct{}
	grimoireIDs   func() map[string]struct{}
}

func NewEntity(deps EntityDeps) *Entity {
	now := deps.Now
	if now == nil {
		now = time.Now
	}
	emptySet := func() map[string]struct{} { return map[string]struct{}{} }
	enemySkillIDs := deps.EnemySkillIDs
	if enemySkillIDs == nil {
		enemySkillIDs = emptySet
	}
	itemIDs := deps.ItemIDs
	if itemIDs == nil {
		itemIDs = emptySet
	}
	grimoireIDs := deps.GrimoireIDs
	if grimoireIDs == nil {
		grimoireIDs = emptySet
	}
	return &Entity{
		mu: deps.Mutex, state: deps.State, repo: deps.Repo, now: now,
		enemySkillIDs: enemySkillIDs, itemIDs: itemIDs, grimoireIDs: grimoireIDs,
	}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[EnemyEntry] {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return ValidateSave(input, e.enemySkillIDs(), e.itemIDs(), e.grimoireIDs(), e.now())
}

func (e *Entity) Create(entry EnemyEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if entity.HasID(e.state.Entries, entry.ID, func(it EnemyEntry) string { return it.ID }) {
		return fmt.Errorf("%w: enemy %s", entity.ErrDuplicateID, entry.ID)
	}
	result := ValidateSave(enemyToInput(entry), e.enemySkillIDs(), e.itemIDs(), e.grimoireIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.state.Entries = append(entity.CopyEntries(e.state.Entries), entry)
	return nil
}

func (e *Entity) Update(entry EnemyEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !entity.HasID(e.state.Entries, entry.ID, func(it EnemyEntry) string { return it.ID }) {
		return fmt.Errorf("%w: enemy %s", entity.ErrNotFound, entry.ID)
	}
	result := ValidateSave(enemyToInput(entry), e.enemySkillIDs(), e.itemIDs(), e.grimoireIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	next, _ := entity.UpsertByID(e.state.Entries, entry, func(it EnemyEntry) string { return it.ID })
	e.state.Entries = next
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next, ok := entity.DeleteByID(e.state.Entries, id, func(it EnemyEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: enemy %s", entity.ErrNotFound, id)
	}
	e.state.Entries = next
	return nil
}

func (e *Entity) Save() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next := entity.CopyEntries(e.state.Entries)
	entity.SortByID(next, func(it EnemyEntry) string { return it.ID })
	e.state.Entries = next
	return e.repo.SaveState(*e.state)
}

func (e *Entity) ListAll() []EnemyEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.CopyEntries(e.state.Entries)
}

func (e *Entity) FindByID(id string) (EnemyEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.FindByID(e.state.Entries, id, func(it EnemyEntry) string { return it.ID })
}

func (e *Entity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func enemyToInput(entry EnemyEntry) SaveInput {
	return SaveInput{
		ID:            entry.ID,
		MobType:       entry.MobType,
		Name:          entry.Name,
		HP:            entry.HP,
		Memo:          entry.Memo,
		Attack:        entry.Attack,
		Defense:       entry.Defense,
		MoveSpeed:     entry.MoveSpeed,
		Equipment:     entry.Equipment,
		EnemySkillIDs: append([]string{}, entry.EnemySkillIDs...),
		DropMode:      entry.DropMode,
		Drops:         append([]DropRef{}, entry.Drops...),
	}
}

func SaveInputFromEntry(entry EnemyEntry) SaveInput {
	return enemyToInput(entry)
}
