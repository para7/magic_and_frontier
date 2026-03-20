package web

import (
	"strconv"
	"strings"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/treasures"
)

func toIDSet[T any](entries []T, idOf func(T) string) map[string]struct{} {
	out := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		id := strings.TrimSpace(idOf(entry))
		if id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}

func compactLines(value string) []string {
	raw := strings.Split(value, "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func splitCSV(value string, max int) []string {
	parts := strings.Split(value, ",")
	if len(parts) > max {
		return nil
	}
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		out = append(out, strings.TrimSpace(part))
	}
	return out
}

func parseRequiredInt(errs map[string]string, key, value string) int {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0
	}
	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		errs[key] = "Must be a number."
		return 0
	}
	return parsed
}

func parseRequiredFloat(errs map[string]string, key, value string) float64 {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0
	}
	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		errs[key] = "Must be a number."
		return 0
	}
	return parsed
}

func parseOptionalFloat(errs map[string]string, key, value string) *float64 {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		errs[key] = "Must be a number."
		return nil
	}
	return &parsed
}

func mergeFieldErrors(primary, secondary map[string]string) map[string]string {
	if len(primary) == 0 && len(secondary) == 0 {
		return map[string]string{}
	}
	out := map[string]string{}
	for key, value := range secondary {
		out[key] = value
	}
	for key, value := range primary {
		out[key] = value
	}
	return out
}

func mapFieldErrors(errs common.FieldErrors, mapField func(string) string) map[string]string {
	out := map[string]string{}
	for key, value := range errs {
		mapped := mapField(key)
		if mapped == "" {
			mapped = key
		}
		if _, exists := out[mapped]; !exists {
			out[mapped] = value
		}
	}
	return out
}

func mapGrimoireField(key string) string {
	return key
}

func mapTreasureField(key string) string {
	if strings.HasPrefix(key, "lootPools.") {
		return "lootPoolsText"
	}
	return key
}

func mapLootTableField(key string) string {
	if strings.HasPrefix(key, "lootPools.") {
		return "lootPoolsText"
	}
	return key
}

func mapEnemyField(key string) string {
	if strings.HasPrefix(key, "enemySkillIds.") {
		return "enemySkillIds"
	}
	if strings.HasPrefix(key, "equipment.") {
		return "equipmentText"
	}
	if strings.HasPrefix(key, "drops.") {
		return "dropsText"
	}
	return key
}

func parseTreasurePools(errs map[string]string, value string) []treasures.DropRef {
	lines := compactLines(value)
	out := make([]treasures.DropRef, 0, len(lines))
	for _, line := range lines {
		parts := splitCSV(line, 5)
		if len(parts) < 3 {
			errs["lootPoolsText"] = "Each loot line must be `kind,refId,weight,countMin,countMax`."
			return nil
		}
		weight, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			errs["lootPoolsText"] = "Weight must be numeric."
			return nil
		}
		out = append(out, treasures.DropRef{
			Kind:     parts[0],
			RefID:    parts[1],
			Weight:   weight,
			CountMin: parseOptionalFloat(errs, "lootPoolsText", valueOrIndex(parts, 3)),
			CountMax: parseOptionalFloat(errs, "lootPoolsText", valueOrIndex(parts, 4)),
		})
		if errs["lootPoolsText"] != "" {
			errs["lootPoolsText"] = "Count values must be numeric when provided."
			return nil
		}
	}
	return out
}

func parseEquipment(errs map[string]string, value string) enemies.Equipment {
	lines := compactLines(value)
	equipment := enemies.Equipment{}
	for _, line := range lines {
		parts := splitCSV(line, 5)
		if len(parts) < 4 {
			errs["equipmentText"] = "Each equipment line must be `slot,kind,refId,count,dropChance`."
			return enemies.Equipment{}
		}
		count := parseRequiredInt(errs, "equipmentText", parts[3])
		if errs["equipmentText"] != "" {
			errs["equipmentText"] = "Equipment count must be numeric."
			return enemies.Equipment{}
		}
		slot := &enemies.EquipmentSlot{
			Kind:       parts[1],
			RefID:      parts[2],
			Count:      count,
			DropChance: parseOptionalFloat(errs, "equipmentText", valueOrIndex(parts, 4)),
		}
		if errs["equipmentText"] != "" {
			errs["equipmentText"] = "Equipment dropChance must be numeric when provided."
			return enemies.Equipment{}
		}
		switch parts[0] {
		case "mainhand":
			equipment.Mainhand = slot
		case "offhand":
			equipment.Offhand = slot
		case "head":
			equipment.Head = slot
		case "chest":
			equipment.Chest = slot
		case "legs":
			equipment.Legs = slot
		case "feet":
			equipment.Feet = slot
		default:
			errs["equipmentText"] = "Equipment slot must be one of mainhand,offhand,head,chest,legs,feet."
			return enemies.Equipment{}
		}
	}
	return equipment
}

func parseEnemyDrops(errs map[string]string, value string) []enemies.DropRef {
	lines := compactLines(value)
	out := make([]enemies.DropRef, 0, len(lines))
	for _, line := range lines {
		parts := splitCSV(line, 5)
		if len(parts) < 3 {
			errs["dropsText"] = "Each drop line must be `kind,refId,weight,countMin,countMax`."
			return nil
		}
		weight, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			errs["dropsText"] = "Drop weight must be numeric."
			return nil
		}
		out = append(out, enemies.DropRef{
			Kind:     parts[0],
			RefID:    parts[1],
			Weight:   weight,
			CountMin: parseOptionalFloat(errs, "dropsText", valueOrIndex(parts, 3)),
			CountMax: parseOptionalFloat(errs, "dropsText", valueOrIndex(parts, 4)),
		})
		if errs["dropsText"] != "" {
			errs["dropsText"] = "Drop count values must be numeric when provided."
			return nil
		}
	}
	return out
}

func valueOrIndex(parts []string, index int) string {
	if index >= 0 && index < len(parts) {
		return parts[index]
	}
	return ""
}

func findEntry[T any](entries []T, id string, idOf func(T) string) (T, bool) {
	var zero T
	for _, entry := range entries {
		if idOf(entry) == id {
			return entry, true
		}
	}
	return zero, false
}

func duplicateCastID(entries []grimoire.GrimoireEntry, id string, castID int) string {
	for _, entry := range entries {
		if entry.ID == id {
			continue
		}
		if entry.CastID == castID {
			return entry.ID
		}
	}
	return ""
}

func duplicateTreasureTablePath(entries []treasures.TreasureEntry, entryID, tablePath string) string {
	for _, entry := range entries {
		if entry.ID != entryID && entry.TablePath == tablePath {
			return entry.ID
		}
	}
	return ""
}
