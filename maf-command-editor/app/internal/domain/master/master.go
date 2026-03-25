package master

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/entity"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/mcsource"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
)

type entryRepo[T any] interface {
	LoadState() (common.EntryState[T], error)
	SaveState(common.EntryState[T]) error
}

type Dependencies struct {
	ItemRepo               entryRepo[items.ItemEntry]
	GrimoireRepo           entryRepo[grimoire.GrimoireEntry]
	SkillRepo              entryRepo[skills.SkillEntry]
	EnemySkillRepo         entryRepo[enemyskills.EnemySkillEntry]
	EnemyRepo              entryRepo[enemies.EnemyEntry]
	SpawnTableRepo         entryRepo[spawntables.SpawnTableEntry]
	TreasureRepo           entryRepo[treasures.TreasureEntry]
	LootTableRepo          entryRepo[loottables.LootTableEntry]
	MinecraftLootTableRoot string
	Now                    func() time.Time
}

type JSONMaster struct {
	mu sync.RWMutex

	itemRepo       entryRepo[items.ItemEntry]
	grimoireRepo   entryRepo[grimoire.GrimoireEntry]
	skillRepo      entryRepo[skills.SkillEntry]
	enemySkillRepo entryRepo[enemyskills.EnemySkillEntry]
	enemyRepo      entryRepo[enemies.EnemyEntry]
	spawnTableRepo entryRepo[spawntables.SpawnTableEntry]
	treasureRepo   entryRepo[treasures.TreasureEntry]
	lootTableRepo  entryRepo[loottables.LootTableEntry]

	itemState       common.EntryState[items.ItemEntry]
	grimoireState   common.EntryState[grimoire.GrimoireEntry]
	skillState      common.EntryState[skills.SkillEntry]
	enemySkillState common.EntryState[enemyskills.EnemySkillEntry]
	enemyState      common.EntryState[enemies.EnemyEntry]
	spawnTableState common.EntryState[spawntables.SpawnTableEntry]
	treasureState   common.EntryState[treasures.TreasureEntry]
	lootTableState  common.EntryState[loottables.LootTableEntry]

	itemsEntity       entity.MafEntity[items.SaveInput, items.ItemEntry]
	grimoiresEntity   entity.MafEntity[grimoire.SaveInput, grimoire.GrimoireEntry]
	skillsEntity      entity.MafEntity[skills.SaveInput, skills.SkillEntry]
	enemySkillsEntity entity.MafEntity[enemyskills.SaveInput, enemyskills.EnemySkillEntry]
	enemiesEntity     entity.MafEntity[enemies.SaveInput, enemies.EnemyEntry]
	treasuresEntity   entity.MafEntity[treasures.SaveInput, treasures.TreasureEntry]
	lootTablesEntity  entity.MafEntity[loottables.SaveInput, loottables.LootTableEntry]
	spawnTablesEntity entity.MafEntity[spawntables.SaveInput, spawntables.SpawnTableEntry]

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

	itemState, err := deps.ItemRepo.LoadState()
	if err != nil {
		return nil, fmt.Errorf("load items: %w", err)
	}
	grimoireState, err := deps.GrimoireRepo.LoadState()
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
	m.initEntities()
	return m, nil
}

func (m *JSONMaster) initEntities() {
	m.itemsEntity = items.NewEntity(items.EntityDeps{
		Mutex: &m.mu,
		State: &m.itemState,
		Repo:  m.itemRepo,
		Now:   m.nowUTC,
		SkillIDs: func() map[string]struct{} {
			return m.skillIDSetLocked()
		},
	})
	m.grimoiresEntity = grimoire.NewEntity(grimoire.EntityDeps{Mutex: &m.mu, State: &m.grimoireState, Repo: m.grimoireRepo, Now: m.nowUTC})
	m.skillsEntity = skills.NewEntity(skills.EntityDeps{Mutex: &m.mu, State: &m.skillState, Repo: m.skillRepo, Now: m.nowUTC, ItemStates: &m.itemState.Entries})
	m.enemySkillsEntity = enemyskills.NewEntity(enemyskills.EntityDeps{Mutex: &m.mu, State: &m.enemySkillState, Repo: m.enemySkillRepo, Now: m.nowUTC, EnemyStates: &m.enemyState.Entries})
	m.enemiesEntity = enemies.NewEntity(enemies.EntityDeps{
		Mutex: &m.mu, State: &m.enemyState, Repo: m.enemyRepo, Now: m.nowUTC,
		EnemySkillIDs: func() map[string]struct{} { return m.enemySkillIDSetLocked() },
		ItemIDs:       func() map[string]struct{} { return m.itemIDSetLocked() },
		GrimoireIDs:   func() map[string]struct{} { return m.grimoireIDSetLocked() },
	})
	m.treasuresEntity = treasures.NewEntity(treasures.EntityDeps{
		Mutex: &m.mu, State: &m.treasureState, Repo: m.treasureRepo, Now: m.nowUTC,
		ItemIDs:             func() map[string]struct{} { return m.itemIDSetLocked() },
		GrimoireIDs:         func() map[string]struct{} { return m.grimoireIDSetLocked() },
		TreasureSourceErr:   func() error { return m.treasureSourceErr },
		TreasureSourcePaths: func() map[string]struct{} { return m.treasureSourcePathsLocked() },
	})
	m.lootTablesEntity = loottables.NewEntity(loottables.EntityDeps{
		Mutex: &m.mu, State: &m.lootTableState, Repo: m.lootTableRepo, Now: m.nowUTC,
		ItemIDs:     func() map[string]struct{} { return m.itemIDSetLocked() },
		GrimoireIDs: func() map[string]struct{} { return m.grimoireIDSetLocked() },
	})
	m.spawnTablesEntity = spawntables.NewEntity(spawntables.EntityDeps{
		Mutex: &m.mu, State: &m.spawnTableState, Repo: m.spawnTableRepo, Now: m.nowUTC,
		EnemyIDs: func() map[string]struct{} { return m.enemyIDSetLocked() },
	})
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

func (m *JSONMaster) Items() entity.MafEntity[items.SaveInput, items.ItemEntry] {
	return m.itemsEntity
}

func (m *JSONMaster) Grimoires() entity.MafEntity[grimoire.SaveInput, grimoire.GrimoireEntry] {
	return m.grimoiresEntity
}

func (m *JSONMaster) Skills() entity.MafEntity[skills.SaveInput, skills.SkillEntry] {
	return m.skillsEntity
}

func (m *JSONMaster) EnemySkills() entity.MafEntity[enemyskills.SaveInput, enemyskills.EnemySkillEntry] {
	return m.enemySkillsEntity
}

func (m *JSONMaster) Enemies() entity.MafEntity[enemies.SaveInput, enemies.EnemyEntry] {
	return m.enemiesEntity
}

func (m *JSONMaster) Treasures() entity.MafEntity[treasures.SaveInput, treasures.TreasureEntry] {
	return m.treasuresEntity
}

func (m *JSONMaster) LootTables() entity.MafEntity[loottables.SaveInput, loottables.LootTableEntry] {
	return m.lootTablesEntity
}

func (m *JSONMaster) SpawnTables() entity.MafEntity[spawntables.SaveInput, spawntables.SpawnTableEntry] {
	return m.spawnTablesEntity
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

func idSet[T any](entries []T, idOf func(T) string) map[string]struct{} {
	return entity.IDSet(entries, idOf)
}

func (m *JSONMaster) itemIDSetLocked() map[string]struct{} {
	return entity.IDSet(m.itemState.Entries, func(entry items.ItemEntry) string { return entry.ID })
}

func (m *JSONMaster) grimoireIDSetLocked() map[string]struct{} {
	return entity.IDSet(m.grimoireState.Entries, func(entry grimoire.GrimoireEntry) string { return entry.ID })
}

func (m *JSONMaster) skillIDSetLocked() map[string]struct{} {
	return entity.IDSet(m.skillState.Entries, func(entry skills.SkillEntry) string { return entry.ID })
}

func (m *JSONMaster) enemySkillIDSetLocked() map[string]struct{} {
	return entity.IDSet(m.enemySkillState.Entries, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
}

func (m *JSONMaster) enemyIDSetLocked() map[string]struct{} {
	return entity.IDSet(m.enemyState.Entries, func(entry enemies.EnemyEntry) string { return entry.ID })
}

func (m *JSONMaster) treasureSourcePathsLocked() map[string]struct{} {
	out := make(map[string]struct{}, len(m.treasureSourcePaths))
	for key := range m.treasureSourcePaths {
		out[key] = struct{}{}
	}
	return out
}
