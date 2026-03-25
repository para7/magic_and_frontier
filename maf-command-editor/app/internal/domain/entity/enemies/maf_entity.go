package enemies

import (
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
)

type EntityDeps struct {
	Mutex         *sync.RWMutex
	State         *common.EntryState[EnemyEntry]
	Repo          entity.EntrySaver[EnemyEntry]
	Now           func() time.Time
	EnemySkillIDs func() map[string]struct{}
	ItemIDs       func() map[string]struct{}
	GrimoireIDs   func() map[string]struct{}
}

type Entity struct {
	entity.BaseEntity[EnemyEntry]
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
		BaseEntity: entity.BaseEntity[EnemyEntry]{
			Mu: deps.Mutex, State: deps.State, Repo: deps.Repo,
			IDOf: func(e EnemyEntry) string { return e.ID },
		},
		now: now, enemySkillIDs: enemySkillIDs, itemIDs: itemIDs, grimoireIDs: grimoireIDs,
	}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[EnemyEntry] {
	e.Mu.RLock()
	defer e.Mu.RUnlock()
	return ValidateSave(input, e.enemySkillIDs(), e.itemIDs(), e.grimoireIDs(), e.now())
}

func (e *Entity) Create(entry EnemyEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.CreateEntry(entry, "enemy"); err != nil {
		return err
	}
	result := ValidateSave(enemyToInput(entry), e.enemySkillIDs(), e.itemIDs(), e.grimoireIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.AppendEntry(entry)
	return nil
}

func (e *Entity) Update(entry EnemyEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.UpdateEntry(entry, "enemy"); err != nil {
		return err
	}
	result := ValidateSave(enemyToInput(entry), e.enemySkillIDs(), e.itemIDs(), e.grimoireIDs(), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.UpsertEntry(entry)
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	return e.DeleteEntry(id, "enemy")
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
