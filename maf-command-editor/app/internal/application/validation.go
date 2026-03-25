package application

import (
	"sort"
	"strings"
	"time"

	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/export"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/mcsource"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
)

func ValidateBundle(states StateBundle, exportSettingsPath string, minecraftLootTableRoot string, now time.Time) ValidationReport {
	report := ValidationReport{
		OK: true,
		Counts: Counts{
			Items:       len(states.ItemState.Entries),
			Grimoire:    len(states.GrimoireState.Entries),
			Skills:      len(states.SkillState.Entries),
			EnemySkills: len(states.EnemySkillState.Entries),
			Enemies:     len(states.EnemyState.Entries),
			SpawnTables: len(states.SpawnTableState.Entries),
			Treasures:   len(states.TreasureState.Entries),
			LootTables:  len(states.LootTableState.Entries),
		},
	}
	if exportSettingsPath != "" {
		if err := export.ValidateSettings(exportSettingsPath); err != nil {
			report.Issues = append(report.Issues, ValidationIssue{
				Entity:  "export_settings",
				Field:   "path",
				Message: err.Error(),
			})
		}
	}

	itemIDs := entryIDs(states.ItemState.Entries, func(entry items.ItemEntry) string { return entry.ID })
	grimoireIDs := entryIDs(states.GrimoireState.Entries, func(entry grimoire.GrimoireEntry) string { return entry.ID })
	skillIDs := entryIDs(states.SkillState.Entries, func(entry skills.SkillEntry) string { return entry.ID })
	enemySkillIDs := entryIDs(states.EnemySkillState.Entries, func(entry enemyskills.EnemySkillEntry) string { return entry.ID })
	enemyIDs := entryIDs(states.EnemyState.Entries, func(entry enemies.EnemyEntry) string { return entry.ID })
	castIDs := map[int]string{}
	treasureTablePaths := map[string]string{}
	validTreasureTablePaths := map[string]struct{}{}
	if sources, err := mcsource.ListLootTables(minecraftLootTableRoot); err != nil {
		report.Issues = append(report.Issues, ValidationIssue{
			Entity:  "minecraft_loot_table_root",
			Field:   "path",
			Message: err.Error(),
		})
	} else {
		for _, source := range sources {
			validTreasureTablePaths[source.TablePath] = struct{}{}
		}
	}

	for _, entry := range states.ItemState.Entries {
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
		appendSaveIssues(&report, "enemy_skill", entry.ID, enemyskills.ValidateSave(enemySkillToInput(entry), now))
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
	for _, entry := range states.SpawnTableState.Entries {
		appendSaveIssues(&report, "spawn_table", entry.ID, spawntables.ValidateSave(spawnTableToInput(entry), enemyIDs, now))
	}
	for _, pair := range spawntables.AllOverlaps(states.SpawnTableState.Entries) {
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
