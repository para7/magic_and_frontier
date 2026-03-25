package enemyskills

import (
	"fmt"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
	"tools2/app/internal/domain/entity/enemies"
)

type EntityDeps struct {
	Mutex       *sync.RWMutex
	State       *common.EntryState[EnemySkillEntry]
	Repo        entity.EntrySaver[EnemySkillEntry]
	Now         func() time.Time
	EnemyStates *[]enemies.EnemyEntry
}

type Entity struct {
	entity.BaseEntity[EnemySkillEntry]
	now         func() time.Time
	enemyStates *[]enemies.EnemyEntry
}

func NewEntity(deps EntityDeps) *Entity {
	now := deps.Now
	if now == nil {
		now = time.Now
	}
	enemiesRef := deps.EnemyStates
	if enemiesRef == nil {
		empty := []enemies.EnemyEntry{}
		enemiesRef = &empty
	}
	return &Entity{
		BaseEntity: entity.BaseEntity[EnemySkillEntry]{
			Mu: deps.Mutex, State: deps.State, Repo: deps.Repo,
			IDOf: func(e EnemySkillEntry) string { return e.ID },
		},
		now: now, enemyStates: enemiesRef,
	}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[EnemySkillEntry] {
	e.Mu.RLock()
	defer e.Mu.RUnlock()
	return ValidateSave(input, e.now())
}

func (e *Entity) Create(entry EnemySkillEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.CreateEntry(entry, "enemy skill"); err != nil {
		return err
	}
	result := ValidateSave(enemySkillToInput(entry), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.AppendEntry(entry)
	return nil
}

func (e *Entity) Update(entry EnemySkillEntry, _ entity.MasterRef) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if err := e.UpdateEntry(entry, "enemy skill"); err != nil {
		return err
	}
	result := ValidateSave(enemySkillToInput(entry), e.now())
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
	for _, enemy := range *e.enemyStates {
		for _, skillID := range enemy.EnemySkillIDs {
			if entity.TrimID(skillID) == trimmedID {
				return fmt.Errorf("%w: enemy skill is referenced by enemy %s", entity.ErrRelation, enemy.ID)
			}
		}
	}
	return e.DeleteEntry(id, "enemy skill")
}

func enemySkillToInput(entry EnemySkillEntry) SaveInput {
	return SaveInput{
		ID:          entry.ID,
		Name:        entry.Name,
		Description: entry.Description,
		Script:      entry.Script,
	}
}

func SaveInputFromEntry(entry EnemySkillEntry) SaveInput {
	return enemySkillToInput(entry)
}
