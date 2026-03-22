package application

import (
	"fmt"
	"time"

	"tools2/app/internal/config"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/export"
	"tools2/app/internal/idseq"
	"tools2/app/internal/store"
)

type Dependencies struct {
	ItemRepo           store.ItemStateRepository
	GrimoireRepo       store.GrimoireStateRepository
	SkillRepo          store.EntryStateRepository[skills.SkillEntry]
	EnemySkillRepo     store.EntryStateRepository[enemyskills.EnemySkillEntry]
	EnemyRepo          store.EntryStateRepository[enemies.EnemyEntry]
	SpawnTableRepo     store.EntryStateRepository[spawntables.SpawnTableEntry]
	TreasureRepo       store.EntryStateRepository[treasures.TreasureEntry]
	LootTableRepo      store.EntryStateRepository[loottables.LootTableEntry]
	CounterRepo        store.CounterRepository
	ExportSettingsPath string
	Now                func() time.Time
}

type Service struct {
	cfg  config.Config
	deps Dependencies
}

type StateBundle struct {
	ItemState       items.ItemState
	GrimoireState   grimoire.GrimoireState
	SkillState      common.EntryState[skills.SkillEntry]
	EnemySkillState common.EntryState[enemyskills.EnemySkillEntry]
	EnemyState      common.EntryState[enemies.EnemyEntry]
	SpawnTableState common.EntryState[spawntables.SpawnTableEntry]
	TreasureState   common.EntryState[treasures.TreasureEntry]
	LootTableState  common.EntryState[loottables.LootTableEntry]
}

type Counts struct {
	Items       int
	Grimoire    int
	Skills      int
	EnemySkills int
	Enemies     int
	SpawnTables int
	Treasures   int
	LootTables  int
}

type ValidationIssue struct {
	Entity  string
	ID      string
	Field   string
	Message string
}

type ValidationReport struct {
	OK     bool
	Counts Counts
	Issues []ValidationIssue
}

func DefaultDependencies(cfg config.Config) Dependencies {
	return Dependencies{
		ItemRepo:           store.NewItemStateRepository(cfg.ItemStatePath),
		GrimoireRepo:       store.NewGrimoireStateRepository(cfg.GrimoireStatePath),
		SkillRepo:          store.NewEntryStateRepository[skills.SkillEntry](cfg.SkillStatePath),
		EnemySkillRepo:     store.NewEntryStateRepository[enemyskills.EnemySkillEntry](cfg.EnemySkillStatePath),
		EnemyRepo:          store.NewEntryStateRepository[enemies.EnemyEntry](cfg.EnemyStatePath),
		SpawnTableRepo:     store.NewEntryStateRepository[spawntables.SpawnTableEntry](cfg.SpawnTableStatePath),
		TreasureRepo:       store.NewEntryStateRepository[treasures.TreasureEntry](cfg.TreasureStatePath),
		LootTableRepo:      store.NewEntryStateRepository[loottables.LootTableEntry](cfg.LootTablesStatePath),
		CounterRepo:        store.NewCounterRepository(cfg.IDCounterStatePath),
		ExportSettingsPath: cfg.ExportSettingsPath,
		Now:                time.Now,
	}
}

func NewService(cfg config.Config, deps Dependencies) Service {
	defaults := DefaultDependencies(cfg)
	if deps.ItemRepo == nil {
		deps.ItemRepo = defaults.ItemRepo
	}
	if deps.GrimoireRepo == nil {
		deps.GrimoireRepo = defaults.GrimoireRepo
	}
	if deps.SkillRepo == nil {
		deps.SkillRepo = defaults.SkillRepo
	}
	if deps.EnemySkillRepo == nil {
		deps.EnemySkillRepo = defaults.EnemySkillRepo
	}
	if deps.EnemyRepo == nil {
		deps.EnemyRepo = defaults.EnemyRepo
	}
	if deps.SpawnTableRepo == nil {
		deps.SpawnTableRepo = defaults.SpawnTableRepo
	}
	if deps.TreasureRepo == nil {
		deps.TreasureRepo = defaults.TreasureRepo
	}
	if deps.LootTableRepo == nil {
		deps.LootTableRepo = defaults.LootTableRepo
	}
	if deps.CounterRepo == nil {
		deps.CounterRepo = defaults.CounterRepo
	}
	if deps.ExportSettingsPath == "" {
		deps.ExportSettingsPath = defaults.ExportSettingsPath
	}
	if deps.Now == nil {
		deps.Now = defaults.Now
	}
	return Service{cfg: cfg, deps: deps}
}

func (s Service) LoadStates() (StateBundle, error) {
	itemState, err := s.deps.ItemRepo.LoadItemState()
	if err != nil {
		return StateBundle{}, fmt.Errorf("load items: %w", err)
	}
	grimoireState, err := s.deps.GrimoireRepo.LoadGrimoireState()
	if err != nil {
		return StateBundle{}, fmt.Errorf("load grimoire: %w", err)
	}
	skillState, err := s.deps.SkillRepo.LoadState()
	if err != nil {
		return StateBundle{}, fmt.Errorf("load skills: %w", err)
	}
	enemySkillState, err := s.deps.EnemySkillRepo.LoadState()
	if err != nil {
		return StateBundle{}, fmt.Errorf("load enemy skills: %w", err)
	}
	enemyState, err := s.deps.EnemyRepo.LoadState()
	if err != nil {
		return StateBundle{}, fmt.Errorf("load enemies: %w", err)
	}
	spawnTableState, err := s.deps.SpawnTableRepo.LoadState()
	if err != nil {
		return StateBundle{}, fmt.Errorf("load spawn tables: %w", err)
	}
	treasureState, err := s.deps.TreasureRepo.LoadState()
	if err != nil {
		return StateBundle{}, fmt.Errorf("load treasures: %w", err)
	}
	lootTableState, err := s.deps.LootTableRepo.LoadState()
	if err != nil {
		return StateBundle{}, fmt.Errorf("load loottables: %w", err)
	}
	return StateBundle{
		ItemState:       itemState,
		GrimoireState:   grimoireState,
		SkillState:      skillState,
		EnemySkillState: enemySkillState,
		EnemyState:      enemyState,
		SpawnTableState: spawnTableState,
		TreasureState:   treasureState,
		LootTableState:  lootTableState,
	}, nil
}

func (s Service) ValidateAll() (ValidationReport, error) {
	states, err := s.LoadStates()
	if err != nil {
		return ValidationReport{}, err
	}
	return ValidateBundle(states, s.deps.ExportSettingsPath, s.cfg.MinecraftLootTableRoot, s.deps.Now()), nil
}

func (s Service) AllocateCastID() (int, error) {
	state, err := s.deps.CounterRepo.LoadCounterState()
	if err != nil {
		return 0, err
	}
	next, castID := idseq.NextCastID(state)
	if err := s.deps.CounterRepo.SaveCounterState(next); err != nil {
		return 0, err
	}
	return castID, nil
}

func (s Service) ExportDatapack() export.SaveDataResponse {
	states, err := s.LoadStates()
	if err != nil {
		return export.SaveDataResponse{
			OK:      false,
			Code:    "LOAD_FAILED",
			Message: "Failed to load savedata.",
			Details: err.Error(),
		}
	}

	if err := export.ValidateSettings(s.deps.ExportSettingsPath); err != nil {
		return export.SaveDataResponse{
			OK:      false,
			Code:    "INVALID_CONFIG",
			Message: "Invalid export settings.",
			Details: err.Error(),
		}
	}

	report := ValidateBundle(states, "", s.cfg.MinecraftLootTableRoot, s.deps.Now())
	if !report.OK {
		return export.SaveDataResponse{
			OK:      false,
			Code:    "VALIDATION_FAILED",
			Message: "Savedata validation failed.",
			Details: report.String(),
		}
	}
	return export.ExportDatapack(export.ExportParams{
		ItemState:              states.ItemState,
		GrimoireState:          states.GrimoireState,
		Skills:                 states.SkillState.Entries,
		EnemySkills:            states.EnemySkillState.Entries,
		Enemies:                states.EnemyState.Entries,
		SpawnTables:            states.SpawnTableState.Entries,
		Treasures:              states.TreasureState.Entries,
		LootTables:             states.LootTableState.Entries,
		ExportSettingsPath:     s.deps.ExportSettingsPath,
		MinecraftLootTableRoot: s.cfg.MinecraftLootTableRoot,
	})
}
