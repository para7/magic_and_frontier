package application

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"tools2/app/internal/config"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/export"
	"tools2/app/internal/idseq"
	"tools2/app/internal/mcsource"
	"tools2/app/internal/store"
)

type Dependencies struct {
	ItemRepo           store.ItemStateRepository
	GrimoireRepo       store.GrimoireStateRepository
	SkillRepo          store.EntryStateRepository[skills.SkillEntry]
	EnemySkillRepo     store.EntryStateRepository[enemyskills.EnemySkillEntry]
	EnemyRepo          store.EntryStateRepository[enemies.EnemyEntry]
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
	TreasureState   common.EntryState[treasures.TreasureEntry]
	LootTableState  common.EntryState[loottables.LootTableEntry]
}

type Counts struct {
	Items       int
	Grimoire    int
	Skills      int
	EnemySkills int
	Enemies     int
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

func (s Service) AllocateID(kind idseq.Kind) (string, error) {
	state, err := s.deps.CounterRepo.LoadCounterState()
	if err != nil {
		return "", err
	}
	next, id := idseq.NextID(state, kind)
	if err := s.deps.CounterRepo.SaveCounterState(next); err != nil {
		return "", err
	}
	return id, nil
}

func (s Service) AllocateGrimoireIdentity() (string, int, error) {
	state, err := s.deps.CounterRepo.LoadCounterState()
	if err != nil {
		return "", 0, err
	}
	next, id := idseq.NextID(state, idseq.KindGrimoire)
	next, castID := idseq.NextCastID(next)
	if err := s.deps.CounterRepo.SaveCounterState(next); err != nil {
		return "", 0, err
	}
	return id, castID, nil
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
		ItemState:          states.ItemState,
		GrimoireState:      states.GrimoireState,
		Skills:             states.SkillState.Entries,
		EnemySkills:        states.EnemySkillState.Entries,
		Enemies:            states.EnemyState.Entries,
		Treasures:          states.TreasureState.Entries,
		LootTables:         states.LootTableState.Entries,
		ExportSettingsPath: s.deps.ExportSettingsPath,
		MinecraftLootTableRoot: s.cfg.MinecraftLootTableRoot,
	})
}

func ValidateBundle(states StateBundle, exportSettingsPath string, minecraftLootTableRoot string, now time.Time) ValidationReport {
	report := ValidationReport{
		OK: true,
		Counts: Counts{
			Items:       len(states.ItemState.Items),
			Grimoire:    len(states.GrimoireState.Entries),
			Skills:      len(states.SkillState.Entries),
			EnemySkills: len(states.EnemySkillState.Entries),
			Enemies:     len(states.EnemyState.Entries),
			Treasures:   len(states.TreasureState.Entries),
			LootTables:  len(states.LootTableState.Entries),
		},
	}
	if exportSettingsPath != "" {
		if err := export.ValidateSettings(exportSettingsPath); err != nil {
			report.Issues = append(report.Issues, ValidationIssue{
				Entity:  "export-settings",
				Field:   "path",
				Message: err.Error(),
			})
		}
	}

	itemIDs := entryIDs(states.ItemState.Items, func(entry items.ItemEntry) string { return entry.ID })
	grimoireIDs := entryIDs(states.GrimoireState.Entries, func(entry grimoire.GrimoireEntry) string { return entry.ID })
	skillIDs := entryIDs(states.SkillState.Entries, func(entry skills.SkillEntry) string { return entry.ID })
	enemySkillIDs := entryIDs(states.EnemySkillState.Entries, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
	castIDs := map[int]string{}
	treasureTablePaths := map[string]string{}
	validTreasureTablePaths := map[string]struct{}{}
	if sources, err := mcsource.ListLootTables(minecraftLootTableRoot); err != nil {
		report.Issues = append(report.Issues, ValidationIssue{
			Entity:  "minecraft-loot-table-root",
			Field:   "path",
			Message: err.Error(),
		})
	} else {
		for _, source := range sources {
			validTreasureTablePaths[source.TablePath] = struct{}{}
		}
	}

	for _, entry := range states.ItemState.Items {
		appendSaveIssues(&report, "item", entry.ID, items.ValidateSave(itemToInput(entry), skillIDs, now))
	}
	for _, entry := range states.GrimoireState.Entries {
		appendSaveIssues(&report, "grimoire", entry.ID, grimoire.ValidateSave(grimoireToInput(entry), now))
		if prevID, exists := castIDs[entry.CastID]; exists && prevID != entry.ID {
			report.Issues = append(report.Issues, ValidationIssue{
				Entity:  "grimoire",
				ID:      entry.ID,
				Field:   "castid",
				Message: "Cast ID is already used by " + prevID + ".",
			})
		} else {
			castIDs[entry.CastID] = entry.ID
		}
	}
	for _, entry := range states.SkillState.Entries {
		appendSaveIssues(&report, "skill", entry.ID, skills.ValidateSave(skillToInput(entry), now))
	}
	for _, entry := range states.EnemySkillState.Entries {
		appendSaveIssues(&report, "enemy-skill", entry.ID, enemyskills.ValidateSave(enemySkillToInput(entry), now))
	}
	for _, entry := range states.TreasureState.Entries {
		appendSaveIssues(&report, "treasure", entry.ID, treasures.ValidateSave(treasureToInput(entry), itemIDs, grimoireIDs, validTreasureTablePaths, now))
		if prevID, exists := treasureTablePaths[strings.TrimSpace(entry.TablePath)]; exists && prevID != entry.ID {
			report.Issues = append(report.Issues, ValidationIssue{
				Entity:  "treasure",
				ID:      entry.ID,
				Field:   "tablePath",
				Message: "Loot table path is already used by " + prevID + ".",
			})
		} else {
			treasureTablePaths[strings.TrimSpace(entry.TablePath)] = entry.ID
		}
	}
	for _, entry := range states.LootTableState.Entries {
		appendSaveIssues(&report, "loottable", entry.ID, loottables.ValidateSave(loottableToInput(entry), itemIDs, grimoireIDs, now))
	}
	for _, entry := range states.EnemyState.Entries {
		appendSaveIssues(&report, "enemy", entry.ID, enemies.ValidateSave(enemyToInput(entry), enemySkillIDs, itemIDs, grimoireIDs, now))
	}

	report.OK = len(report.Issues) == 0
	sort.Slice(report.Issues, func(i, j int) bool {
		left := report.Issues[i]
		right := report.Issues[j]
		if left.Entity != right.Entity {
			return left.Entity < right.Entity
		}
		if left.ID != right.ID {
			return left.ID < right.ID
		}
		if left.Field != right.Field {
			return left.Field < right.Field
		}
		return left.Message < right.Message
	})
	return report
}

func (r ValidationReport) String() string {
	if r.OK {
		return "ok"
	}
	lines := make([]string, 0, len(r.Issues))
	for _, issue := range r.Issues {
		label := issue.Entity
		if issue.ID != "" {
			label += "[" + issue.ID + "]"
		}
		if issue.Field != "" {
			label += "." + issue.Field
		}
		lines = append(lines, label+": "+issue.Message)
	}
	return strings.Join(lines, "\n")
}

func appendSaveIssues[T any](report *ValidationReport, entity, id string, result common.SaveResult[T]) {
	if result.OK {
		return
	}
	if len(result.FieldErrors) == 0 {
		report.Issues = append(report.Issues, ValidationIssue{
			Entity:  entity,
			ID:      id,
			Message: result.FormError,
		})
		return
	}
	for field, message := range result.FieldErrors {
		report.Issues = append(report.Issues, ValidationIssue{
			Entity:  entity,
			ID:      id,
			Field:   field,
			Message: message,
		})
	}
}

func itemToInput(entry items.ItemEntry) items.SaveInput {
	return items.SaveInput{
		ID:                  entry.ID,
		ItemID:              entry.ItemID,
		Count:               entry.Count,
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
		LootPools: append([]treasures.DropRef{}, entry.LootPools...),
	}
}

func enemyToInput(entry enemies.EnemyEntry) enemies.SaveInput {
	return enemies.SaveInput{
		ID:            entry.ID,
		MobType:       entry.MobType,
		Name:          entry.Name,
		HP:            entry.HP,
		Attack:        entry.Attack,
		Defense:       entry.Defense,
		MoveSpeed:     entry.MoveSpeed,
		Equipment:     entry.Equipment,
		EnemySkillIDs: append([]string{}, entry.EnemySkillIDs...),
		DropMode:      entry.DropMode,
		Drops:         append([]enemies.DropRef{}, entry.Drops...),
	}
}

func entryIDs[T any](entries []T, idOf func(T) string) map[string]struct{} {
	out := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		id := strings.TrimSpace(idOf(entry))
		if id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}
