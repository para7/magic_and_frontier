package master

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/domain/mcsource"
	"tools2/app/internal/domain/store"
)

type Dependencies struct {
	ItemRepo               store.ItemStateRepository
	GrimoireRepo           store.GrimoireStateRepository
	SkillRepo              store.EntryStateRepository[skills.SkillEntry]
	EnemySkillRepo         store.EntryStateRepository[enemyskills.EnemySkillEntry]
	EnemyRepo              store.EntryStateRepository[enemies.EnemyEntry]
	SpawnTableRepo         store.EntryStateRepository[spawntables.SpawnTableEntry]
	TreasureRepo           store.EntryStateRepository[treasures.TreasureEntry]
	LootTableRepo          store.EntryStateRepository[loottables.LootTableEntry]
	MinecraftLootTableRoot string
	Now                    func() time.Time
}

type JSONMaster struct {
	mu sync.RWMutex

	itemRepo       store.ItemStateRepository
	grimoireRepo   store.GrimoireStateRepository
	skillRepo      store.EntryStateRepository[skills.SkillEntry]
	enemySkillRepo store.EntryStateRepository[enemyskills.EnemySkillEntry]
	enemyRepo      store.EntryStateRepository[enemies.EnemyEntry]
	spawnTableRepo store.EntryStateRepository[spawntables.SpawnTableEntry]
	treasureRepo   store.EntryStateRepository[treasures.TreasureEntry]
	lootTableRepo  store.EntryStateRepository[loottables.LootTableEntry]

	itemState       items.ItemState
	grimoireState   grimoire.GrimoireState
	skillState      common.EntryState[skills.SkillEntry]
	enemySkillState common.EntryState[enemyskills.EnemySkillEntry]
	enemyState      common.EntryState[enemies.EnemyEntry]
	spawnTableState common.EntryState[spawntables.SpawnTableEntry]
	treasureState   common.EntryState[treasures.TreasureEntry]
	lootTableState  common.EntryState[loottables.LootTableEntry]

	minecraftLootTableRoot string
	treasureSourcePaths    map[string]struct{}
	treasureSourceErr      error
	now                    func() time.Time
}

func NewJSONMaster(deps Dependencies) (*JSONMaster, error) {
	if deps.Now == nil {
		deps.Now = time.Now
	}
	if deps.ItemRepo == nil || deps.GrimoireRepo == nil || deps.SkillRepo == nil || deps.EnemySkillRepo == nil || deps.EnemyRepo == nil || deps.SpawnTableRepo == nil || deps.TreasureRepo == nil || deps.LootTableRepo == nil {
		return nil, fmt.Errorf("master dependencies are incomplete")
	}

	itemState, err := deps.ItemRepo.LoadItemState()
	if err != nil {
		return nil, fmt.Errorf("load items: %w", err)
	}
	grimoireState, err := deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		return nil, fmt.Errorf("load grimoire: %w", err)
	}
	skillState, err := deps.SkillRepo.LoadState()
	if err != nil {
		return nil, fmt.Errorf("load skills: %w", err)
	}
	enemySkillState, err := deps.EnemySkillRepo.LoadState()
	if err != nil {
		return nil, fmt.Errorf("load enemy skills: %w", err)
	}
	enemyState, err := deps.EnemyRepo.LoadState()
	if err != nil {
		return nil, fmt.Errorf("load enemies: %w", err)
	}
	spawnTableState, err := deps.SpawnTableRepo.LoadState()
	if err != nil {
		return nil, fmt.Errorf("load spawn tables: %w", err)
	}
	treasureState, err := deps.TreasureRepo.LoadState()
	if err != nil {
		return nil, fmt.Errorf("load treasures: %w", err)
	}
	lootTableState, err := deps.LootTableRepo.LoadState()
	if err != nil {
		return nil, fmt.Errorf("load loottables: %w", err)
	}

	m := &JSONMaster{
		itemRepo:               deps.ItemRepo,
		grimoireRepo:           deps.GrimoireRepo,
		skillRepo:              deps.SkillRepo,
		enemySkillRepo:         deps.EnemySkillRepo,
		enemyRepo:              deps.EnemyRepo,
		spawnTableRepo:         deps.SpawnTableRepo,
		treasureRepo:           deps.TreasureRepo,
		lootTableRepo:          deps.LootTableRepo,
		itemState:              itemState,
		grimoireState:          grimoireState,
		skillState:             skillState,
		enemySkillState:        enemySkillState,
		enemyState:             enemyState,
		spawnTableState:        spawnTableState,
		treasureState:          treasureState,
		lootTableState:         lootTableState,
		minecraftLootTableRoot: deps.MinecraftLootTableRoot,
		now:                    deps.Now,
		treasureSourcePaths:    map[string]struct{}{},
	}
	m.refreshTreasureSourcesLocked()
	return m, nil
}

func (m *JSONMaster) refreshTreasureSourcesLocked() {
	m.treasureSourcePaths = map[string]struct{}{}
	m.treasureSourceErr = nil
	if strings.TrimSpace(m.minecraftLootTableRoot) == "" {
		m.treasureSourceErr = fmt.Errorf("minecraft loot table root is empty")
		return
	}
	sources, err := mcsource.ListLootTables(m.minecraftLootTableRoot)
	if err != nil {
		m.treasureSourceErr = err
		return
	}
	for _, source := range sources {
		if !treasures.IsSupportedTablePath(source.TablePath) {
			continue
		}
		m.treasureSourcePaths[source.TablePath] = struct{}{}
	}
}

func (m *JSONMaster) HasItem(id string) bool {
	_, ok := m.Items().FindByID(id)
	return ok
}

func (m *JSONMaster) HasGrimoire(id string) bool {
	_, ok := m.Grimoires().FindByID(id)
	return ok
}

func (m *JSONMaster) HasSkill(id string) bool {
	_, ok := m.Skills().FindByID(id)
	return ok
}

func (m *JSONMaster) HasEnemySkill(id string) bool {
	_, ok := m.EnemySkills().FindByID(id)
	return ok
}

func (m *JSONMaster) HasEnemy(id string) bool {
	_, ok := m.Enemies().FindByID(id)
	return ok
}

func (m *JSONMaster) HasTreasure(id string) bool {
	_, ok := m.Treasures().FindByID(id)
	return ok
}

func (m *JSONMaster) HasLootTable(id string) bool {
	_, ok := m.LootTables().FindByID(id)
	return ok
}

func (m *JSONMaster) HasSpawnTable(id string) bool {
	_, ok := m.SpawnTables().FindByID(id)
	return ok
}

func (m *JSONMaster) Items() MafEntity[items.SaveInput, items.ItemEntry] {
	return itemEntity{m: m}
}

func (m *JSONMaster) Grimoires() MafEntity[grimoire.SaveInput, grimoire.GrimoireEntry] {
	return grimoireEntity{m: m}
}

func (m *JSONMaster) Skills() MafEntity[skills.SaveInput, skills.SkillEntry] {
	return skillEntity{m: m}
}

func (m *JSONMaster) EnemySkills() MafEntity[enemyskills.SaveInput, enemyskills.EnemySkillEntry] {
	return enemySkillEntity{m: m}
}

func (m *JSONMaster) Enemies() MafEntity[enemies.SaveInput, enemies.EnemyEntry] {
	return enemyEntity{m: m}
}

func (m *JSONMaster) Treasures() MafEntity[treasures.SaveInput, treasures.TreasureEntry] {
	return treasureEntity{m: m}
}

func (m *JSONMaster) LootTables() MafEntity[loottables.SaveInput, loottables.LootTableEntry] {
	return loottableEntity{m: m}
}

func (m *JSONMaster) SpawnTables() MafEntity[spawntables.SaveInput, spawntables.SpawnTableEntry] {
	return spawnTableEntity{m: m}
}

func (m *JSONMaster) SaveAll() error {
	if err := m.Items().Save(); err != nil {
		return err
	}
	if err := m.Grimoires().Save(); err != nil {
		return err
	}
	if err := m.Skills().Save(); err != nil {
		return err
	}
	if err := m.EnemySkills().Save(); err != nil {
		return err
	}
	if err := m.Enemies().Save(); err != nil {
		return err
	}
	if err := m.SpawnTables().Save(); err != nil {
		return err
	}
	if err := m.Treasures().Save(); err != nil {
		return err
	}
	if err := m.LootTables().Save(); err != nil {
		return err
	}
	return nil
}

func (m *JSONMaster) nowUTC() time.Time {
	return m.now()
}

func copyEntries[T any](entries []T) []T {
	return append([]T{}, entries...)
}

func trimID(value string) string {
	return strings.TrimSpace(value)
}

func hasID[T any](entries []T, id string, idOf func(T) string) bool {
	id = trimID(id)
	if id == "" {
		return false
	}
	for _, entry := range entries {
		if trimID(idOf(entry)) == id {
			return true
		}
	}
	return false
}

func findByID[T any](entries []T, id string, idOf func(T) string) (T, bool) {
	var zero T
	id = trimID(id)
	if id == "" {
		return zero, false
	}
	for _, entry := range entries {
		if trimID(idOf(entry)) == id {
			return entry, true
		}
	}
	return zero, false
}

func upsertByID[T any](entries []T, entry T, idOf func(T) string) ([]T, bool) {
	entryID := trimID(idOf(entry))
	next := copyEntries(entries)
	for i := range next {
		if trimID(idOf(next[i])) == entryID {
			next[i] = entry
			return next, true
		}
	}
	next = append(next, entry)
	return next, false
}

func deleteByID[T any](entries []T, id string, idOf func(T) string) ([]T, bool) {
	id = trimID(id)
	next := make([]T, 0, len(entries))
	found := false
	for _, entry := range entries {
		if trimID(idOf(entry)) == id {
			found = true
			continue
		}
		next = append(next, entry)
	}
	return next, found
}

func sortByID[T any](entries []T, idOf func(T) string) {
	sort.Slice(entries, func(i, j int) bool {
		return trimID(idOf(entries[i])) < trimID(idOf(entries[j]))
	})
}

func idSet[T any](entries []T, idOf func(T) string) map[string]struct{} {
	out := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		id := trimID(idOf(entry))
		if id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}

func firstFieldError(errs common.FieldErrors) string {
	if len(errs) == 0 {
		return "validation failed"
	}
	keys := make([]string, 0, len(errs))
	for key := range errs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	first := keys[0]
	return first + ": " + errs[first]
}

func relationErrFromResult[T any](result common.SaveResult[T]) error {
	if len(result.FieldErrors) > 0 {
		return fmt.Errorf("%w: %s", ErrRelation, firstFieldError(result.FieldErrors))
	}
	if strings.TrimSpace(result.FormError) != "" {
		return fmt.Errorf("%w: %s", ErrRelation, result.FormError)
	}
	return fmt.Errorf("%w", ErrRelation)
}

func (m *JSONMaster) duplicateCastIDLocked(entryID string, castID int) string {
	for _, entry := range m.grimoireState.Entries {
		if entry.ID != entryID && entry.CastID == castID {
			return entry.ID
		}
	}
	return ""
}

func (m *JSONMaster) duplicateTreasureTablePathLocked(entryID, tablePath string) string {
	tablePath = strings.TrimSpace(tablePath)
	for _, entry := range m.treasureState.Entries {
		if entry.ID != entryID && strings.TrimSpace(entry.TablePath) == tablePath {
			return entry.ID
		}
	}
	return ""
}

func (m *JSONMaster) firstSpawnOverlapLocked(entry spawntables.SpawnTableEntry) (string, bool) {
	entries := make([]spawntables.SpawnTableEntry, 0, len(m.spawnTableState.Entries))
	for _, it := range m.spawnTableState.Entries {
		if it.ID == entry.ID {
			continue
		}
		entries = append(entries, it)
	}
	return spawntables.FirstOverlap(entries, entry)
}

func (m *JSONMaster) itemIDSetLocked() map[string]struct{} {
	return idSet(m.itemState.Items, func(entry items.ItemEntry) string { return entry.ID })
}

func (m *JSONMaster) grimoireIDSetLocked() map[string]struct{} {
	return idSet(m.grimoireState.Entries, func(entry grimoire.GrimoireEntry) string { return entry.ID })
}

func (m *JSONMaster) skillIDSetLocked() map[string]struct{} {
	return idSet(m.skillState.Entries, func(entry skills.SkillEntry) string { return entry.ID })
}

func (m *JSONMaster) enemySkillIDSetLocked() map[string]struct{} {
	return idSet(m.enemySkillState.Entries, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
}

func (m *JSONMaster) enemyIDSetLocked() map[string]struct{} {
	return idSet(m.enemyState.Entries, func(entry enemies.EnemyEntry) string { return entry.ID })
}
