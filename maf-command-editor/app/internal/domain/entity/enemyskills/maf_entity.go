package enemyskills

import (
	"fmt"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/entity"
	"tools2/app/internal/domain/entity/enemies"
)

type entrySaver interface {
	SaveState(common.EntryState[EnemySkillEntry]) error
}

type EntityDeps struct {
	Mutex       *sync.RWMutex
	State       *common.EntryState[EnemySkillEntry]
	Repo        entrySaver
	Now         func() time.Time
	EnemyStates *[]enemies.EnemyEntry
}

type Entity struct {
	mu          *sync.RWMutex
	state       *common.EntryState[EnemySkillEntry]
	repo        entrySaver
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
	return &Entity{mu: deps.Mutex, state: deps.State, repo: deps.Repo, now: now, enemyStates: enemiesRef}
}

func (e *Entity) Validate(input SaveInput, _ entity.MasterRef) common.SaveResult[EnemySkillEntry] {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return ValidateSave(input, e.now())
}

func (e *Entity) Create(entry EnemySkillEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if entity.HasID(e.state.Entries, entry.ID, func(it EnemySkillEntry) string { return it.ID }) {
		return fmt.Errorf("%w: enemy skill %s", entity.ErrDuplicateID, entry.ID)
	}
	result := ValidateSave(enemySkillToInput(entry), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	e.state.Entries = append(entity.CopyEntries(e.state.Entries), entry)
	return nil
}

func (e *Entity) Update(entry EnemySkillEntry, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !entity.HasID(e.state.Entries, entry.ID, func(it EnemySkillEntry) string { return it.ID }) {
		return fmt.Errorf("%w: enemy skill %s", entity.ErrNotFound, entry.ID)
	}
	result := ValidateSave(enemySkillToInput(entry), e.now())
	if !result.OK {
		return entity.RelationErrFromResult(result)
	}
	next, _ := entity.UpsertByID(e.state.Entries, entry, func(it EnemySkillEntry) string { return it.ID })
	e.state.Entries = next
	return nil
}

func (e *Entity) Delete(id string, _ entity.MasterRef) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	trimmedID := entity.TrimID(id)
	for _, enemy := range *e.enemyStates {
		for _, skillID := range enemy.EnemySkillIDs {
			if entity.TrimID(skillID) == trimmedID {
				return fmt.Errorf("%w: enemy skill is referenced by enemy %s", entity.ErrRelation, enemy.ID)
			}
		}
	}
	next, ok := entity.DeleteByID(e.state.Entries, id, func(it EnemySkillEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: enemy skill %s", entity.ErrNotFound, id)
	}
	e.state.Entries = next
	return nil
}

func (e *Entity) Save() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	next := entity.CopyEntries(e.state.Entries)
	entity.SortByID(next, func(it EnemySkillEntry) string { return it.ID })
	e.state.Entries = next
	return e.repo.SaveState(*e.state)
}

func (e *Entity) ListAll() []EnemySkillEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.CopyEntries(e.state.Entries)
}

func (e *Entity) FindByID(id string) (EnemySkillEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return entity.FindByID(e.state.Entries, id, func(it EnemySkillEntry) string { return it.ID })
}

func (e *Entity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
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
