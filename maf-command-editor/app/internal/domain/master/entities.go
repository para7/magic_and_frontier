package master

import (
	"fmt"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
)

type itemEntity struct{ m *JSONMaster }
type grimoireEntity struct{ m *JSONMaster }
type skillEntity struct{ m *JSONMaster }
type enemySkillEntity struct{ m *JSONMaster }
type enemyEntity struct{ m *JSONMaster }
type treasureEntity struct{ m *JSONMaster }
type loottableEntity struct{ m *JSONMaster }
type spawnTableEntity struct{ m *JSONMaster }

func (e itemEntity) Validate(input items.SaveInput, _ DBMaster) common.SaveResult[items.ItemEntry] {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return items.ValidateSave(input, e.m.skillIDSetLocked(), e.m.nowUTC())
}

func (e itemEntity) Create(entry items.ItemEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if hasID(e.m.itemState.Items, entry.ID, func(it items.ItemEntry) string { return it.ID }) {
		return fmt.Errorf("%w: item %s", ErrDuplicateID, entry.ID)
	}
	result := items.ValidateSave(itemToInput(entry), e.m.skillIDSetLocked(), e.m.nowUTC())
	if !result.OK {
		return relationErrFromResult(result)
	}
	e.m.itemState.Items = append(copyEntries(e.m.itemState.Items), entry)
	return nil
}

func (e itemEntity) Update(entry items.ItemEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if !hasID(e.m.itemState.Items, entry.ID, func(it items.ItemEntry) string { return it.ID }) {
		return fmt.Errorf("%w: item %s", ErrNotFound, entry.ID)
	}
	result := items.ValidateSave(itemToInput(entry), e.m.skillIDSetLocked(), e.m.nowUTC())
	if !result.OK {
		return relationErrFromResult(result)
	}
	next, _ := upsertByID(e.m.itemState.Items, entry, func(it items.ItemEntry) string { return it.ID })
	e.m.itemState.Items = next
	return nil
}

func (e itemEntity) Delete(id string, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next, ok := deleteByID(e.m.itemState.Items, id, func(it items.ItemEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: item %s", ErrNotFound, id)
	}
	e.m.itemState.Items = next
	return nil
}

func (e itemEntity) Save() error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next := copyEntries(e.m.itemState.Items)
	sortByID(next, func(it items.ItemEntry) string { return it.ID })
	e.m.itemState.Items = next
	return e.m.itemRepo.SaveItemState(e.m.itemState)
}

func (e itemEntity) ListAll() []items.ItemEntry {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return copyEntries(e.m.itemState.Items)
}

func (e itemEntity) FindByID(id string) (items.ItemEntry, bool) {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return findByID(e.m.itemState.Items, id, func(it items.ItemEntry) string { return it.ID })
}

func (e itemEntity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func (e grimoireEntity) Validate(input grimoire.SaveInput, _ DBMaster) common.SaveResult[grimoire.GrimoireEntry] {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return e.m.validateGrimoireLocked(input)
}

func (m *JSONMaster) validateGrimoireLocked(input grimoire.SaveInput) common.SaveResult[grimoire.GrimoireEntry] {
	result := grimoire.ValidateSave(input, m.nowUTC())
	if !result.OK {
		return result
	}
	if conflictID := m.duplicateCastIDLocked(result.Entry.ID, result.Entry.CastID); conflictID != "" {
		return common.SaveValidationError[grimoire.GrimoireEntry](
			common.FieldErrors{"castid": "Cast ID is already used by " + conflictID + "."},
			"Validation failed. Fix the highlighted fields.",
		)
	}
	return result
}

func (e grimoireEntity) Create(entry grimoire.GrimoireEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if hasID(e.m.grimoireState.Entries, entry.ID, func(it grimoire.GrimoireEntry) string { return it.ID }) {
		return fmt.Errorf("%w: grimoire %s", ErrDuplicateID, entry.ID)
	}
	result := e.m.validateGrimoireLocked(grimoireToInput(entry))
	if !result.OK {
		return relationErrFromResult(result)
	}
	e.m.grimoireState.Entries = append(copyEntries(e.m.grimoireState.Entries), entry)
	return nil
}

func (e grimoireEntity) Update(entry grimoire.GrimoireEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if !hasID(e.m.grimoireState.Entries, entry.ID, func(it grimoire.GrimoireEntry) string { return it.ID }) {
		return fmt.Errorf("%w: grimoire %s", ErrNotFound, entry.ID)
	}
	result := e.m.validateGrimoireLocked(grimoireToInput(entry))
	if !result.OK {
		return relationErrFromResult(result)
	}
	next, _ := upsertByID(e.m.grimoireState.Entries, entry, func(it grimoire.GrimoireEntry) string { return it.ID })
	e.m.grimoireState.Entries = next
	return nil
}

func (e grimoireEntity) Delete(id string, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next, ok := deleteByID(e.m.grimoireState.Entries, id, func(it grimoire.GrimoireEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: grimoire %s", ErrNotFound, id)
	}
	e.m.grimoireState.Entries = next
	return nil
}

func (e grimoireEntity) Save() error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next := copyEntries(e.m.grimoireState.Entries)
	sortByID(next, func(it grimoire.GrimoireEntry) string { return it.ID })
	e.m.grimoireState.Entries = next
	return e.m.grimoireRepo.SaveGrimoireState(e.m.grimoireState)
}

func (e grimoireEntity) ListAll() []grimoire.GrimoireEntry {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return copyEntries(e.m.grimoireState.Entries)
}

func (e grimoireEntity) FindByID(id string) (grimoire.GrimoireEntry, bool) {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return findByID(e.m.grimoireState.Entries, id, func(it grimoire.GrimoireEntry) string { return it.ID })
}

func (e grimoireEntity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func (e skillEntity) Validate(input skills.SaveInput, _ DBMaster) common.SaveResult[skills.SkillEntry] {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return skills.ValidateSave(input, e.m.nowUTC())
}

func (e skillEntity) Create(entry skills.SkillEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if hasID(e.m.skillState.Entries, entry.ID, func(it skills.SkillEntry) string { return it.ID }) {
		return fmt.Errorf("%w: skill %s", ErrDuplicateID, entry.ID)
	}
	result := skills.ValidateSave(skillToInput(entry), e.m.nowUTC())
	if !result.OK {
		return relationErrFromResult(result)
	}
	e.m.skillState.Entries = append(copyEntries(e.m.skillState.Entries), entry)
	return nil
}

func (e skillEntity) Update(entry skills.SkillEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if !hasID(e.m.skillState.Entries, entry.ID, func(it skills.SkillEntry) string { return it.ID }) {
		return fmt.Errorf("%w: skill %s", ErrNotFound, entry.ID)
	}
	result := skills.ValidateSave(skillToInput(entry), e.m.nowUTC())
	if !result.OK {
		return relationErrFromResult(result)
	}
	next, _ := upsertByID(e.m.skillState.Entries, entry, func(it skills.SkillEntry) string { return it.ID })
	e.m.skillState.Entries = next
	return nil
}

func (e skillEntity) Delete(id string, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	for _, item := range e.m.itemState.Items {
		if strings.TrimSpace(item.SkillID) == strings.TrimSpace(id) {
			return fmt.Errorf("%w: skill is referenced by item %s", ErrRelation, item.ID)
		}
	}
	next, ok := deleteByID(e.m.skillState.Entries, id, func(it skills.SkillEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: skill %s", ErrNotFound, id)
	}
	e.m.skillState.Entries = next
	return nil
}

func (e skillEntity) Save() error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next := copyEntries(e.m.skillState.Entries)
	sortByID(next, func(it skills.SkillEntry) string { return it.ID })
	e.m.skillState.Entries = next
	return e.m.skillRepo.SaveState(e.m.skillState)
}

func (e skillEntity) ListAll() []skills.SkillEntry {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return copyEntries(e.m.skillState.Entries)
}

func (e skillEntity) FindByID(id string) (skills.SkillEntry, bool) {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return findByID(e.m.skillState.Entries, id, func(it skills.SkillEntry) string { return it.ID })
}

func (e skillEntity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func (e enemySkillEntity) Validate(input enemyskills.SaveInput, _ DBMaster) common.SaveResult[enemyskills.EnemySkillEntry] {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return enemyskills.ValidateSave(input, e.m.nowUTC())
}

func (e enemySkillEntity) Create(entry enemyskills.EnemySkillEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if hasID(e.m.enemySkillState.Entries, entry.ID, func(it enemyskills.EnemySkillEntry) string { return it.ID }) {
		return fmt.Errorf("%w: enemy skill %s", ErrDuplicateID, entry.ID)
	}
	result := enemyskills.ValidateSave(enemySkillToInput(entry), e.m.nowUTC())
	if !result.OK {
		return relationErrFromResult(result)
	}
	e.m.enemySkillState.Entries = append(copyEntries(e.m.enemySkillState.Entries), entry)
	return nil
}

func (e enemySkillEntity) Update(entry enemyskills.EnemySkillEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if !hasID(e.m.enemySkillState.Entries, entry.ID, func(it enemyskills.EnemySkillEntry) string { return it.ID }) {
		return fmt.Errorf("%w: enemy skill %s", ErrNotFound, entry.ID)
	}
	result := enemyskills.ValidateSave(enemySkillToInput(entry), e.m.nowUTC())
	if !result.OK {
		return relationErrFromResult(result)
	}
	next, _ := upsertByID(e.m.enemySkillState.Entries, entry, func(it enemyskills.EnemySkillEntry) string { return it.ID })
	e.m.enemySkillState.Entries = next
	return nil
}

func (e enemySkillEntity) Delete(id string, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	for _, enemy := range e.m.enemyState.Entries {
		for _, skillID := range enemy.EnemySkillIDs {
			if strings.TrimSpace(skillID) == strings.TrimSpace(id) {
				return fmt.Errorf("%w: enemy skill is referenced by enemy %s", ErrRelation, enemy.ID)
			}
		}
	}
	next, ok := deleteByID(e.m.enemySkillState.Entries, id, func(it enemyskills.EnemySkillEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: enemy skill %s", ErrNotFound, id)
	}
	e.m.enemySkillState.Entries = next
	return nil
}

func (e enemySkillEntity) Save() error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next := copyEntries(e.m.enemySkillState.Entries)
	sortByID(next, func(it enemyskills.EnemySkillEntry) string { return it.ID })
	e.m.enemySkillState.Entries = next
	return e.m.enemySkillRepo.SaveState(e.m.enemySkillState)
}

func (e enemySkillEntity) ListAll() []enemyskills.EnemySkillEntry {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return copyEntries(e.m.enemySkillState.Entries)
}

func (e enemySkillEntity) FindByID(id string) (enemyskills.EnemySkillEntry, bool) {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return findByID(e.m.enemySkillState.Entries, id, func(it enemyskills.EnemySkillEntry) string { return it.ID })
}

func (e enemySkillEntity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func (e enemyEntity) Validate(input enemies.SaveInput, _ DBMaster) common.SaveResult[enemies.EnemyEntry] {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return enemies.ValidateSave(input, e.m.enemySkillIDSetLocked(), e.m.itemIDSetLocked(), e.m.grimoireIDSetLocked(), e.m.nowUTC())
}

func (e enemyEntity) Create(entry enemies.EnemyEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if hasID(e.m.enemyState.Entries, entry.ID, func(it enemies.EnemyEntry) string { return it.ID }) {
		return fmt.Errorf("%w: enemy %s", ErrDuplicateID, entry.ID)
	}
	result := enemies.ValidateSave(enemyToInput(entry), e.m.enemySkillIDSetLocked(), e.m.itemIDSetLocked(), e.m.grimoireIDSetLocked(), e.m.nowUTC())
	if !result.OK {
		return relationErrFromResult(result)
	}
	e.m.enemyState.Entries = append(copyEntries(e.m.enemyState.Entries), entry)
	return nil
}

func (e enemyEntity) Update(entry enemies.EnemyEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if !hasID(e.m.enemyState.Entries, entry.ID, func(it enemies.EnemyEntry) string { return it.ID }) {
		return fmt.Errorf("%w: enemy %s", ErrNotFound, entry.ID)
	}
	result := enemies.ValidateSave(enemyToInput(entry), e.m.enemySkillIDSetLocked(), e.m.itemIDSetLocked(), e.m.grimoireIDSetLocked(), e.m.nowUTC())
	if !result.OK {
		return relationErrFromResult(result)
	}
	next, _ := upsertByID(e.m.enemyState.Entries, entry, func(it enemies.EnemyEntry) string { return it.ID })
	e.m.enemyState.Entries = next
	return nil
}

func (e enemyEntity) Delete(id string, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next, ok := deleteByID(e.m.enemyState.Entries, id, func(it enemies.EnemyEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: enemy %s", ErrNotFound, id)
	}
	e.m.enemyState.Entries = next
	return nil
}

func (e enemyEntity) Save() error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next := copyEntries(e.m.enemyState.Entries)
	sortByID(next, func(it enemies.EnemyEntry) string { return it.ID })
	e.m.enemyState.Entries = next
	return e.m.enemyRepo.SaveState(e.m.enemyState)
}

func (e enemyEntity) ListAll() []enemies.EnemyEntry {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return copyEntries(e.m.enemyState.Entries)
}

func (e enemyEntity) FindByID(id string) (enemies.EnemyEntry, bool) {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return findByID(e.m.enemyState.Entries, id, func(it enemies.EnemyEntry) string { return it.ID })
}

func (e enemyEntity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func (e treasureEntity) Validate(input treasures.SaveInput, _ DBMaster) common.SaveResult[treasures.TreasureEntry] {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return e.m.validateTreasureLocked(input)
}

func (m *JSONMaster) validateTreasureLocked(input treasures.SaveInput) common.SaveResult[treasures.TreasureEntry] {
	if m.treasureSourceErr != nil {
		return common.SaveValidationError[treasures.TreasureEntry](
			common.FieldErrors{"tablePath": m.treasureSourceErr.Error()},
			"Validation failed. Fix the highlighted fields.",
		)
	}
	result := treasures.ValidateSave(input, m.itemIDSetLocked(), m.grimoireIDSetLocked(), m.treasureSourcePaths, m.nowUTC())
	if !result.OK {
		return result
	}
	if conflictID := m.duplicateTreasureTablePathLocked(result.Entry.ID, result.Entry.TablePath); conflictID != "" {
		return common.SaveValidationError[treasures.TreasureEntry](
			common.FieldErrors{"tablePath": "Loot table path is already used by " + conflictID + "."},
			"Validation failed. Fix the highlighted fields.",
		)
	}
	return result
}

func (e treasureEntity) Create(entry treasures.TreasureEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if hasID(e.m.treasureState.Entries, entry.ID, func(it treasures.TreasureEntry) string { return it.ID }) {
		return fmt.Errorf("%w: treasure %s", ErrDuplicateID, entry.ID)
	}
	result := e.m.validateTreasureLocked(treasureToInput(entry))
	if !result.OK {
		return relationErrFromResult(result)
	}
	e.m.treasureState.Entries = append(copyEntries(e.m.treasureState.Entries), entry)
	return nil
}

func (e treasureEntity) Update(entry treasures.TreasureEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if !hasID(e.m.treasureState.Entries, entry.ID, func(it treasures.TreasureEntry) string { return it.ID }) {
		return fmt.Errorf("%w: treasure %s", ErrNotFound, entry.ID)
	}
	result := e.m.validateTreasureLocked(treasureToInput(entry))
	if !result.OK {
		return relationErrFromResult(result)
	}
	next, _ := upsertByID(e.m.treasureState.Entries, entry, func(it treasures.TreasureEntry) string { return it.ID })
	e.m.treasureState.Entries = next
	return nil
}

func (e treasureEntity) Delete(id string, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next, ok := deleteByID(e.m.treasureState.Entries, id, func(it treasures.TreasureEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: treasure %s", ErrNotFound, id)
	}
	e.m.treasureState.Entries = next
	return nil
}

func (e treasureEntity) Save() error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next := copyEntries(e.m.treasureState.Entries)
	sortByID(next, func(it treasures.TreasureEntry) string { return it.ID })
	e.m.treasureState.Entries = next
	return e.m.treasureRepo.SaveState(e.m.treasureState)
}

func (e treasureEntity) ListAll() []treasures.TreasureEntry {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return copyEntries(e.m.treasureState.Entries)
}

func (e treasureEntity) FindByID(id string) (treasures.TreasureEntry, bool) {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return findByID(e.m.treasureState.Entries, id, func(it treasures.TreasureEntry) string { return it.ID })
}

func (e treasureEntity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func (e loottableEntity) Validate(input loottables.SaveInput, _ DBMaster) common.SaveResult[loottables.LootTableEntry] {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return loottables.ValidateSave(input, e.m.itemIDSetLocked(), e.m.grimoireIDSetLocked(), e.m.nowUTC())
}

func (e loottableEntity) Create(entry loottables.LootTableEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if hasID(e.m.lootTableState.Entries, entry.ID, func(it loottables.LootTableEntry) string { return it.ID }) {
		return fmt.Errorf("%w: loottable %s", ErrDuplicateID, entry.ID)
	}
	result := loottables.ValidateSave(loottableToInput(entry), e.m.itemIDSetLocked(), e.m.grimoireIDSetLocked(), e.m.nowUTC())
	if !result.OK {
		return relationErrFromResult(result)
	}
	e.m.lootTableState.Entries = append(copyEntries(e.m.lootTableState.Entries), entry)
	return nil
}

func (e loottableEntity) Update(entry loottables.LootTableEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if !hasID(e.m.lootTableState.Entries, entry.ID, func(it loottables.LootTableEntry) string { return it.ID }) {
		return fmt.Errorf("%w: loottable %s", ErrNotFound, entry.ID)
	}
	result := loottables.ValidateSave(loottableToInput(entry), e.m.itemIDSetLocked(), e.m.grimoireIDSetLocked(), e.m.nowUTC())
	if !result.OK {
		return relationErrFromResult(result)
	}
	next, _ := upsertByID(e.m.lootTableState.Entries, entry, func(it loottables.LootTableEntry) string { return it.ID })
	e.m.lootTableState.Entries = next
	return nil
}

func (e loottableEntity) Delete(id string, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next, ok := deleteByID(e.m.lootTableState.Entries, id, func(it loottables.LootTableEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: loottable %s", ErrNotFound, id)
	}
	e.m.lootTableState.Entries = next
	return nil
}

func (e loottableEntity) Save() error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next := copyEntries(e.m.lootTableState.Entries)
	sortByID(next, func(it loottables.LootTableEntry) string { return it.ID })
	e.m.lootTableState.Entries = next
	return e.m.lootTableRepo.SaveState(e.m.lootTableState)
}

func (e loottableEntity) ListAll() []loottables.LootTableEntry {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return copyEntries(e.m.lootTableState.Entries)
}

func (e loottableEntity) FindByID(id string) (loottables.LootTableEntry, bool) {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return findByID(e.m.lootTableState.Entries, id, func(it loottables.LootTableEntry) string { return it.ID })
}

func (e loottableEntity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func (e spawnTableEntity) Validate(input spawntables.SaveInput, _ DBMaster) common.SaveResult[spawntables.SpawnTableEntry] {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return e.m.validateSpawnTableLocked(input)
}

func (m *JSONMaster) validateSpawnTableLocked(input spawntables.SaveInput) common.SaveResult[spawntables.SpawnTableEntry] {
	result := spawntables.ValidateSave(input, m.enemyIDSetLocked(), m.nowUTC())
	if !result.OK {
		return result
	}
	if conflictID, ok := m.firstSpawnOverlapLocked(*result.Entry); ok {
		return common.SaveValidationError[spawntables.SpawnTableEntry](
			common.FieldErrors{"range": "Range overlaps with " + conflictID + "."},
			"Validation failed. Fix the highlighted fields.",
		)
	}
	return result
}

func (e spawnTableEntity) Create(entry spawntables.SpawnTableEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if hasID(e.m.spawnTableState.Entries, entry.ID, func(it spawntables.SpawnTableEntry) string { return it.ID }) {
		return fmt.Errorf("%w: spawn table %s", ErrDuplicateID, entry.ID)
	}
	result := e.m.validateSpawnTableLocked(spawnTableToInput(entry))
	if !result.OK {
		return relationErrFromResult(result)
	}
	e.m.spawnTableState.Entries = append(copyEntries(e.m.spawnTableState.Entries), entry)
	return nil
}

func (e spawnTableEntity) Update(entry spawntables.SpawnTableEntry, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	if !hasID(e.m.spawnTableState.Entries, entry.ID, func(it spawntables.SpawnTableEntry) string { return it.ID }) {
		return fmt.Errorf("%w: spawn table %s", ErrNotFound, entry.ID)
	}
	result := e.m.validateSpawnTableLocked(spawnTableToInput(entry))
	if !result.OK {
		return relationErrFromResult(result)
	}
	next, _ := upsertByID(e.m.spawnTableState.Entries, entry, func(it spawntables.SpawnTableEntry) string { return it.ID })
	e.m.spawnTableState.Entries = next
	return nil
}

func (e spawnTableEntity) Delete(id string, _ DBMaster) error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next, ok := deleteByID(e.m.spawnTableState.Entries, id, func(it spawntables.SpawnTableEntry) string { return it.ID })
	if !ok {
		return fmt.Errorf("%w: spawn table %s", ErrNotFound, id)
	}
	e.m.spawnTableState.Entries = next
	return nil
}

func (e spawnTableEntity) Save() error {
	e.m.mu.Lock()
	defer e.m.mu.Unlock()
	next := copyEntries(e.m.spawnTableState.Entries)
	sortByID(next, func(it spawntables.SpawnTableEntry) string { return it.ID })
	e.m.spawnTableState.Entries = next
	return e.m.spawnTableRepo.SaveState(e.m.spawnTableState)
}

func (e spawnTableEntity) ListAll() []spawntables.SpawnTableEntry {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return copyEntries(e.m.spawnTableState.Entries)
}

func (e spawnTableEntity) FindByID(id string) (spawntables.SpawnTableEntry, bool) {
	e.m.mu.RLock()
	defer e.m.mu.RUnlock()
	return findByID(e.m.spawnTableState.Entries, id, func(it spawntables.SpawnTableEntry) string { return it.ID })
}

func (e spawnTableEntity) HasID(id string) bool {
	_, ok := e.FindByID(id)
	return ok
}

func itemToInput(entry items.ItemEntry) items.SaveInput {
	return items.SaveInput{
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

func grimoireToInput(entry grimoire.GrimoireEntry) grimoire.SaveInput {
	return grimoire.SaveInput{
		ID:          entry.ID,
		CastID:      entry.CastID,
		CastTime:    entry.CastTime,
		MPCost:      entry.MPCost,
		Script:      entry.Script,
		Title:       entry.Title,
		Description: entry.Description,
	}
}

func skillToInput(entry skills.SkillEntry) skills.SaveInput {
	return skills.SaveInput{
		ID:          entry.ID,
		Name:        entry.Name,
		SkillType:   entry.SkillType,
		Description: entry.Description,
		Script:      entry.Script,
	}
}

func enemySkillToInput(entry enemyskills.EnemySkillEntry) enemyskills.SaveInput {
	return enemyskills.SaveInput{
		ID:          entry.ID,
		Name:        entry.Name,
		Description: entry.Description,
		Script:      entry.Script,
	}
}

func treasureToInput(entry treasures.TreasureEntry) treasures.SaveInput {
	return treasures.SaveInput{
		ID:        entry.ID,
		TablePath: entry.TablePath,
		LootPools: append([]treasures.DropRef{}, entry.LootPools...),
	}
}

func loottableToInput(entry loottables.LootTableEntry) loottables.SaveInput {
	return loottables.SaveInput{
		ID:        entry.ID,
		Memo:      entry.Memo,
		LootPools: append([]treasures.DropRef{}, entry.LootPools...),
	}
}

func enemyToInput(entry enemies.EnemyEntry) enemies.SaveInput {
	return enemies.SaveInput{
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
		Drops:         append([]enemies.DropRef{}, entry.Drops...),
	}
}

func spawnTableToInput(entry spawntables.SpawnTableEntry) spawntables.SaveInput {
	return spawntables.SaveInput{
		ID:            entry.ID,
		SourceMobType: entry.SourceMobType,
		Dimension:     entry.Dimension,
		MinX:          entry.MinX,
		MaxX:          entry.MaxX,
		MinY:          entry.MinY,
		MaxY:          entry.MaxY,
		MinZ:          entry.MinZ,
		MaxZ:          entry.MaxZ,
		BaseMobWeight: entry.BaseMobWeight,
		Replacements:  append([]spawntables.ReplacementEntry{}, entry.Replacements...),
	}
}
