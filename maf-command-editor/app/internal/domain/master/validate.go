package master

import (
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

func (m *JSONMaster) ValidateSavedAll() ValidationReport {
	m.mu.RLock()
	itemState := m.itemState
	grimoireState := m.grimoireState
	skillState := m.skillState
	enemySkillState := m.enemySkillState
	enemyState := m.enemyState
	spawnTableState := m.spawnTableState
	treasureState := m.treasureState
	lootTableState := m.lootTableState
	treasureSourcePaths := make(map[string]struct{}, len(m.treasureSourcePaths))
	for key := range m.treasureSourcePaths {
		treasureSourcePaths[key] = struct{}{}
	}
	treasureSourceErr := m.treasureSourceErr
	m.mu.RUnlock()

	report := ValidationReport{
		OK: true,
		Counts: Counts{
			Items:       len(itemState.Items),
			Grimoire:    len(grimoireState.Entries),
			Skills:      len(skillState.Entries),
			EnemySkills: len(enemySkillState.Entries),
			Enemies:     len(enemyState.Entries),
			SpawnTables: len(spawnTableState.Entries),
			Treasures:   len(treasureState.Entries),
			LootTables:  len(lootTableState.Entries),
		},
	}

	itemIDs := idSet(itemState.Items, func(entry items.ItemEntry) string { return entry.ID })
	grimoireIDs := idSet(grimoireState.Entries, func(entry grimoire.GrimoireEntry) string { return entry.ID })
	skillIDs := idSet(skillState.Entries, func(entry skills.SkillEntry) string { return entry.ID })
	enemySkillIDs := idSet(enemySkillState.Entries, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
	enemyIDs := idSet(enemyState.Entries, func(entry enemies.EnemyEntry) string { return entry.ID })
	castIDs := map[int]string{}
	treasureTablePaths := map[string]string{}

	if treasureSourceErr != nil {
		report.Issues = append(report.Issues, ValidationIssue{
			Entity:  "minecraft_loot_table_root",
			Field:   "path",
			Message: treasureSourceErr.Error(),
		})
	}

	for _, entry := range itemState.Items {
		appendSaveIssues(&report, "item", entry.ID, items.ValidateSave(itemToInput(entry), skillIDs, m.nowUTC()))
	}
	for _, entry := range grimoireState.Entries {
		appendSaveIssues(&report, "grimoire", entry.ID, grimoire.ValidateSave(grimoireToInput(entry), m.nowUTC()))
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
	for _, entry := range skillState.Entries {
		appendSaveIssues(&report, "skill", entry.ID, skills.ValidateSave(skillToInput(entry), m.nowUTC()))
	}
	for _, entry := range enemySkillState.Entries {
		appendSaveIssues(&report, "enemy_skill", entry.ID, enemyskills.ValidateSave(enemySkillToInput(entry), m.nowUTC()))
	}
	for _, entry := range treasureState.Entries {
		appendSaveIssues(&report, "treasure", entry.ID, treasures.ValidateSave(treasureToInput(entry), itemIDs, grimoireIDs, treasureSourcePaths, m.nowUTC()))
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
	for _, entry := range lootTableState.Entries {
		appendSaveIssues(&report, "loottable", entry.ID, loottables.ValidateSave(loottableToInput(entry), itemIDs, grimoireIDs, m.nowUTC()))
	}
	for _, entry := range enemyState.Entries {
		appendSaveIssues(&report, "enemy", entry.ID, enemies.ValidateSave(enemyToInput(entry), enemySkillIDs, itemIDs, grimoireIDs, m.nowUTC()))
	}
	for _, entry := range spawnTableState.Entries {
		appendSaveIssues(&report, "spawn_table", entry.ID, spawntables.ValidateSave(spawnTableToInput(entry), enemyIDs, m.nowUTC()))
	}
	for _, pair := range spawntables.AllOverlaps(spawnTableState.Entries) {
		report.Issues = append(report.Issues, ValidationIssue{
			Entity:  "spawn_table",
			ID:      pair[0],
			Field:   "range",
			Message: "Range overlaps with " + pair[1] + ".",
		})
		report.Issues = append(report.Issues, ValidationIssue{
			Entity:  "spawn_table",
			ID:      pair[1],
			Field:   "range",
			Message: "Range overlaps with " + pair[0] + ".",
		})
	}

	report.OK = len(report.Issues) == 0
	return report.Sorted()
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
