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
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/export"
	"tools2/app/internal/store"
)

type Dependencies struct {
	ItemRepo           store.ItemStateRepository
	GrimoireRepo       store.GrimoireStateRepository
	SkillRepo          store.EntryStateRepository[skills.SkillEntry]
	EnemySkillRepo     store.EntryStateRepository[enemyskills.EnemySkillEntry]
	EnemyRepo          store.EntryStateRepository[enemies.EnemyEntry]
	TreasureRepo       store.EntryStateRepository[treasures.TreasureEntry]
	ExportSettingsPath string
	Now                func() time.Time
}

type Service struct {
	deps Dependencies
}

type StateBundle struct {
	ItemState       items.ItemState
	GrimoireState   grimoire.GrimoireState
	SkillState      common.EntryState[skills.SkillEntry]
	EnemySkillState common.EntryState[enemyskills.EnemySkillEntry]
	EnemyState      common.EntryState[enemies.EnemyEntry]
	TreasureState   common.EntryState[treasures.TreasureEntry]
}

type Counts struct {
	Items       int
	Grimoire    int
	Skills      int
	EnemySkills int
	Enemies     int
	Treasures   int
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
	if deps.ExportSettingsPath == "" {
		deps.ExportSettingsPath = defaults.ExportSettingsPath
	}
	if deps.Now == nil {
		deps.Now = defaults.Now
	}
	return Service{deps: deps}
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
	return StateBundle{
		ItemState:       itemState,
		GrimoireState:   grimoireState,
		SkillState:      skillState,
		EnemySkillState: enemySkillState,
		EnemyState:      enemyState,
		TreasureState:   treasureState,
	}, nil
}

func (s Service) ValidateAll() (ValidationReport, error) {
	states, err := s.LoadStates()
	if err != nil {
		return ValidationReport{}, err
	}
	return ValidateBundle(states, s.deps.ExportSettingsPath, s.deps.Now()), nil
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

	report := ValidateBundle(states, "", s.deps.Now())
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
		ExportSettingsPath: s.deps.ExportSettingsPath,
	})
}

func ValidateBundle(states StateBundle, exportSettingsPath string, now time.Time) ValidationReport {
	report := ValidationReport{
		OK: true,
		Counts: Counts{
			Items:       len(states.ItemState.Items),
			Grimoire:    len(states.GrimoireState.Entries),
			Skills:      len(states.SkillState.Entries),
			EnemySkills: len(states.EnemySkillState.Entries),
			Enemies:     len(states.EnemyState.Entries),
			Treasures:   len(states.TreasureState.Entries),
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
	enemySkillIDs := entryIDs(states.EnemySkillState.Entries, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
	treasureIDs := entryIDs(states.TreasureState.Entries, func(entry treasures.TreasureEntry) string { return entry.ID })

	for _, entry := range states.ItemState.Items {
		appendSaveIssues(&report, "item", entry.ID, items.ValidateSave(itemToInput(entry), now))
	}
	for _, entry := range states.GrimoireState.Entries {
		appendSaveIssues(&report, "grimoire", entry.ID, grimoire.ValidateSave(grimoireToInput(entry), now))
	}
	for _, entry := range states.SkillState.Entries {
		appendSaveIssues(&report, "skill", entry.ID, skills.ValidateSave(skillToInput(entry), itemIDs, now))
	}
	for _, entry := range states.EnemySkillState.Entries {
		appendSaveIssues(&report, "enemy-skill", entry.ID, enemyskills.ValidateSave(enemySkillToInput(entry), now))
	}
	for _, entry := range states.TreasureState.Entries {
		appendSaveIssues(&report, "treasure", entry.ID, treasures.ValidateSave(treasureToInput(entry), itemIDs, grimoireIDs, now))
	}
	for _, entry := range states.EnemyState.Entries {
		appendSaveIssues(&report, "enemy", entry.ID, enemies.ValidateSave(enemyToInput(entry), enemySkillIDs, itemIDs, grimoireIDs, now))
		if strings.TrimSpace(entry.DropTableID) != "" {
			if _, ok := treasureIDs[strings.TrimSpace(entry.DropTableID)]; !ok {
				report.Issues = append(report.Issues, ValidationIssue{
					Entity:  "enemy",
					ID:      entry.ID,
					Field:   "dropTableId",
					Message: "Referenced treasure does not exist.",
				})
			}
		}
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
		Script:      entry.Script,
		Title:       entry.Title,
		Description: entry.Description,
		Variants:    append([]grimoire.Variant{}, entry.Variants...),
	}
}

func skillToInput(entry skills.SkillEntry) skills.SaveInput {
	return skills.SaveInput{
		ID:     entry.ID,
		Name:   entry.Name,
		Script: entry.Script,
		ItemID: entry.ItemID,
	}
}

func enemySkillToInput(entry enemyskills.EnemySkillEntry) enemyskills.SaveInput {
	input := enemyskills.SaveInput{
		ID:       entry.ID,
		Name:     entry.Name,
		Script:   entry.Script,
		Cooldown: entry.Cooldown,
	}
	if entry.Trigger != nil {
		input.Trigger = string(*entry.Trigger)
	}
	return input
}

func treasureToInput(entry treasures.TreasureEntry) treasures.SaveInput {
	return treasures.SaveInput{
		ID:        entry.ID,
		Name:      entry.Name,
		LootPools: append([]treasures.DropRef{}, entry.LootPools...),
	}
}

func enemyToInput(entry enemies.EnemyEntry) enemies.SaveInput {
	return enemies.SaveInput{
		ID:            entry.ID,
		Name:          entry.Name,
		HP:            entry.HP,
		Attack:        entry.Attack,
		Defense:       entry.Defense,
		MoveSpeed:     entry.MoveSpeed,
		DropTableID:   entry.DropTableID,
		EnemySkillIDs: append([]string{}, entry.EnemySkillIDs...),
		SpawnRule:     entry.SpawnRule,
		DropTable:     append([]enemies.DropRef{}, entry.DropTable...),
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
