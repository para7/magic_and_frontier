package application

import (
	"fmt"
	"time"

	"maf-command-editor/app/internal/config"
	"maf-command-editor/app/internal/domain/common"
	"maf-command-editor/app/internal/domain/entity/enemies"
	"maf-command-editor/app/internal/domain/entity/enemyskills"
	"maf-command-editor/app/internal/domain/entity/grimoire"
	"maf-command-editor/app/internal/domain/entity/items"
	"maf-command-editor/app/internal/domain/entity/loottables"
	"maf-command-editor/app/internal/domain/entity/skills"
	"maf-command-editor/app/internal/domain/entity/spawntables"
	"maf-command-editor/app/internal/domain/entity/treasures"
	"maf-command-editor/app/internal/domain/export"
	"maf-command-editor/app/internal/domain/idseq"
	dmaster "maf-command-editor/app/internal/domain/master"
)

type Dependencies struct {
	ItemRepo           common.StateRepository[items.ItemEntry]
	GrimoireRepo       common.StateRepository[grimoire.GrimoireEntry]
	SkillRepo          common.StateRepository[skills.SkillEntry]
	EnemySkillRepo     common.StateRepository[enemyskills.EnemySkillEntry]
	EnemyRepo          common.StateRepository[enemies.EnemyEntry]
	SpawnTableRepo     common.StateRepository[spawntables.SpawnTableEntry]
	TreasureRepo       common.StateRepository[treasures.TreasureEntry]
	LootTableRepo      common.StateRepository[loottables.LootTableEntry]
	CounterRepo        idseq.Repository
	Master             dmaster.DBMaster
	ExportSettingsPath string
	Now                func() time.Time
}

type Service struct {
	cfg       config.Config
	deps      Dependencies
	masterErr error
}

type StateBundle struct {
	ItemState       common.EntryState[items.ItemEntry]
	GrimoireState   common.EntryState[grimoire.GrimoireEntry]
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
		ItemRepo:           common.StateRepository[items.ItemEntry]{Path: cfg.ItemStatePath},
		GrimoireRepo:       common.StateRepository[grimoire.GrimoireEntry]{Path: cfg.GrimoireStatePath},
		SkillRepo:          common.StateRepository[skills.SkillEntry]{Path: cfg.SkillStatePath},
		EnemySkillRepo:     common.StateRepository[enemyskills.EnemySkillEntry]{Path: cfg.EnemySkillStatePath},
		EnemyRepo:          common.StateRepository[enemies.EnemyEntry]{Path: cfg.EnemyStatePath},
		SpawnTableRepo:     common.StateRepository[spawntables.SpawnTableEntry]{Path: cfg.SpawnTableStatePath},
		TreasureRepo:       common.StateRepository[treasures.TreasureEntry]{Path: cfg.TreasureStatePath},
		LootTableRepo:      common.StateRepository[loottables.LootTableEntry]{Path: cfg.LootTablesStatePath},
		CounterRepo:        idseq.Repository{Path: cfg.IDCounterStatePath},
		ExportSettingsPath: cfg.ExportSettingsPath,
		Now:                time.Now,
	}
}

func NewService(cfg config.Config, deps Dependencies) Service {
	defaults := DefaultDependencies(cfg)
	if deps.ItemRepo.Path == "" {
		deps.ItemRepo = defaults.ItemRepo
	}
	if deps.GrimoireRepo.Path == "" {
		deps.GrimoireRepo = defaults.GrimoireRepo
	}
	if deps.SkillRepo.Path == "" {
		deps.SkillRepo = defaults.SkillRepo
	}
	if deps.EnemySkillRepo.Path == "" {
		deps.EnemySkillRepo = defaults.EnemySkillRepo
	}
	if deps.EnemyRepo.Path == "" {
		deps.EnemyRepo = defaults.EnemyRepo
	}
	if deps.SpawnTableRepo.Path == "" {
		deps.SpawnTableRepo = defaults.SpawnTableRepo
	}
	if deps.TreasureRepo.Path == "" {
		deps.TreasureRepo = defaults.TreasureRepo
	}
	if deps.LootTableRepo.Path == "" {
		deps.LootTableRepo = defaults.LootTableRepo
	}
	if deps.CounterRepo.Path == "" {
		deps.CounterRepo = defaults.CounterRepo
	}
	if deps.ExportSettingsPath == "" {
		deps.ExportSettingsPath = defaults.ExportSettingsPath
	}
	if deps.Now == nil {
		deps.Now = defaults.Now
	}
	var masterErr error
	if deps.Master == nil {
		var loadedMaster dmaster.DBMaster
		loadedMaster, masterErr = dmaster.NewJSONMaster(dmaster.Dependencies{
			ItemRepo:               deps.ItemRepo,
			GrimoireRepo:           deps.GrimoireRepo,
			SkillRepo:              deps.SkillRepo,
			EnemySkillRepo:         deps.EnemySkillRepo,
			EnemyRepo:              deps.EnemyRepo,
			SpawnTableRepo:         deps.SpawnTableRepo,
			TreasureRepo:           deps.TreasureRepo,
			LootTableRepo:          deps.LootTableRepo,
			MinecraftLootTableRoot: cfg.MinecraftLootTableRoot,
			Now:                    deps.Now,
		})
		if masterErr == nil {
			deps.Master = loadedMaster
		}
	}
	return Service{cfg: cfg, deps: deps, masterErr: masterErr}
}

func (s Service) Master() (dmaster.DBMaster, error) {
	if s.masterErr != nil {
		return nil, s.masterErr
	}
	if s.deps.Master == nil {
		return nil, fmt.Errorf("master is not initialized")
	}
	return s.deps.Master, nil
}

func (s Service) LoadStates() (StateBundle, error) {
	if s.masterErr != nil {
		return StateBundle{}, s.masterErr
	}
	if s.deps.Master != nil {
		return StateBundle{
			ItemState:       common.EntryState[items.ItemEntry]{Entries: s.deps.Master.Items().ListAll()},
			GrimoireState:   common.EntryState[grimoire.GrimoireEntry]{Entries: s.deps.Master.Grimoires().ListAll()},
			SkillState:      common.EntryState[skills.SkillEntry]{Entries: s.deps.Master.Skills().ListAll()},
			EnemySkillState: common.EntryState[enemyskills.EnemySkillEntry]{Entries: s.deps.Master.EnemySkills().ListAll()},
			EnemyState:      common.EntryState[enemies.EnemyEntry]{Entries: s.deps.Master.Enemies().ListAll()},
			SpawnTableState: common.EntryState[spawntables.SpawnTableEntry]{Entries: s.deps.Master.SpawnTables().ListAll()},
			TreasureState:   common.EntryState[treasures.TreasureEntry]{Entries: s.deps.Master.Treasures().ListAll()},
			LootTableState:  common.EntryState[loottables.LootTableEntry]{Entries: s.deps.Master.LootTables().ListAll()},
		}, nil
	}
	itemState, err := s.deps.ItemRepo.LoadState()
	if err != nil {
		return StateBundle{}, fmt.Errorf("load items: %w", err)
	}
	grimoireState, err := s.deps.GrimoireRepo.LoadState()
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
	if s.masterErr != nil {
		return ValidationReport{}, s.masterErr
	}
	if s.deps.Master != nil {
		report := fromMasterReport(s.deps.Master.ValidateSavedAll())
		if s.deps.ExportSettingsPath != "" {
			if err := export.ValidateSettings(s.deps.ExportSettingsPath); err != nil {
				report.Issues = append(report.Issues, ValidationIssue{
					Entity:  "export_settings",
					Field:   "path",
					Message: err.Error(),
				})
			}
		}
		report.OK = len(report.Issues) == 0
		return report, nil
	}
	states, err := s.LoadStates()
	if err != nil {
		return ValidationReport{}, err
	}
	return ValidateBundle(states, s.deps.ExportSettingsPath, s.cfg.MinecraftLootTableRoot, s.deps.Now()), nil
}

func (s Service) AllocateCastID() (int, error) {
	state, err := s.deps.CounterRepo.Load()
	if err != nil {
		return 0, err
	}
	next, castID := idseq.NextCastID(state)
	if err := s.deps.CounterRepo.Save(next); err != nil {
		return 0, err
	}
	return castID, nil
}

func (s Service) ExportDatapack() export.SaveDataResponse {
	if s.masterErr != nil {
		return export.SaveDataResponse{
			OK:      false,
			Code:    "LOAD_FAILED",
			Message: "Failed to load savedata.",
			Details: s.masterErr.Error(),
		}
	}
	if s.deps.Master != nil {
		return export.ExportDatapackFromMaster(s.deps.Master, export.MasterExportParams{
			ExportSettingsPath:     s.deps.ExportSettingsPath,
			MinecraftLootTableRoot: s.cfg.MinecraftLootTableRoot,
		})
	}
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
		Items:                  states.ItemState.Entries,
		Grimoires:              states.GrimoireState.Entries,
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

func fromMasterReport(report dmaster.ValidationReport) ValidationReport {
	issues := make([]ValidationIssue, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, ValidationIssue{
			Entity:  issue.Entity,
			ID:      issue.ID,
			Field:   issue.Field,
			Message: issue.Message,
		})
	}
	return ValidationReport{
		OK: report.OK,
		Counts: Counts{
			Items:       report.Counts.Items,
			Grimoire:    report.Counts.Grimoire,
			Skills:      report.Counts.Skills,
			EnemySkills: report.Counts.EnemySkills,
			Enemies:     report.Counts.Enemies,
			SpawnTables: report.Counts.SpawnTables,
			Treasures:   report.Counts.Treasures,
			LootTables:  report.Counts.LootTables,
		},
		Issues: issues,
	}
}
