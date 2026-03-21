package spawntables

import (
	"fmt"
	"time"

	"tools2/app/internal/domain/common"
)

func ValidateSave(input SaveInput, enemyIDs map[string]struct{}, now time.Time) common.SaveResult[SpawnTableEntry] {
	errs := common.ViolationsToFieldErrors(common.ValidateStruct(input), common.DefaultValidationMessage)
	id := common.RequirePrefixedSequenceID(errs, "id", input.ID, "spawntable_")
	sourceMobType := common.NormalizeText(input.SourceMobType)
	if sourceMobType != "" && !common.IsNamespacedResourceID(sourceMobType) {
		errs.Add("sourceMobType", "Must be a minecraft entity id.")
	}

	replacements := make([]ReplacementEntry, 0, len(input.Replacements))
	seenEnemy := map[string]struct{}{}
	for i, replacement := range input.Replacements {
		enemyID := common.NormalizeText(replacement.EnemyID)
		if !common.IsPrefixedSequenceID(enemyID, "enemy_") {
			errs.Add(fmt.Sprintf("replacements.%d.enemyId", i), "Invalid ID format.")
			continue
		}
		if _, ok := enemyIDs[enemyID]; !ok {
			errs.Add(fmt.Sprintf("replacements.%d.enemyId", i), "Referenced enemy does not exist.")
			continue
		}
		if _, exists := seenEnemy[enemyID]; exists {
			errs.Add(fmt.Sprintf("replacements.%d.enemyId", i), "Duplicate enemy id is not allowed.")
			continue
		}
		if replacement.Weight <= 0 {
			errs.Add(fmt.Sprintf("replacements.%d.weight", i), "Must be greater than 0.")
			continue
		}
		seenEnemy[enemyID] = struct{}{}
		replacements = append(replacements, ReplacementEntry{EnemyID: enemyID, Weight: replacement.Weight})
	}

	if input.MinX > input.MaxX {
		errs.Add("minX", "Must be <= maxX.")
	}
	if input.MinY > input.MaxY {
		errs.Add("minY", "Must be <= maxY.")
	}
	if input.MinZ > input.MaxZ {
		errs.Add("minZ", "Must be <= maxZ.")
	}

	totalWeight := input.BaseMobWeight
	for _, replacement := range replacements {
		totalWeight += replacement.Weight
	}
	if totalWeight <= 0 {
		errs.Add("baseMobWeight", "baseMobWeight + replacement weights must be greater than 0.")
	}

	if errs.Any() {
		return common.SaveValidationError[SpawnTableEntry](errs, "Validation failed. Fix the highlighted fields.")
	}

	entry := SpawnTableEntry{
		ID:            id,
		SourceMobType: sourceMobType,
		Dimension:     common.NormalizeText(input.Dimension),
		MinX:          input.MinX,
		MaxX:          input.MaxX,
		MinY:          input.MinY,
		MaxY:          input.MaxY,
		MinZ:          input.MinZ,
		MaxZ:          input.MaxZ,
		BaseMobWeight: input.BaseMobWeight,
		Replacements:  replacements,
		UpdatedAt:     now.UTC().Format(time.RFC3339),
	}
	return common.SaveSuccess(entry, common.SaveModeCreated)
}

func RangesOverlap(left, right SpawnTableEntry) bool {
	return intervalOverlap(left.MinX, left.MaxX, right.MinX, right.MaxX) &&
		intervalOverlap(left.MinY, left.MaxY, right.MinY, right.MaxY) &&
		intervalOverlap(left.MinZ, left.MaxZ, right.MinZ, right.MaxZ)
}

func FirstOverlap(entries []SpawnTableEntry, candidate SpawnTableEntry) (string, bool) {
	for _, entry := range entries {
		if entry.ID == candidate.ID {
			continue
		}
		if entry.SourceMobType != candidate.SourceMobType || entry.Dimension != candidate.Dimension {
			continue
		}
		if RangesOverlap(entry, candidate) {
			return entry.ID, true
		}
	}
	return "", false
}

func AllOverlaps(entries []SpawnTableEntry) [][2]string {
	pairs := make([][2]string, 0)
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			left := entries[i]
			right := entries[j]
			if left.SourceMobType != right.SourceMobType || left.Dimension != right.Dimension {
				continue
			}
			if RangesOverlap(left, right) {
				pairs = append(pairs, [2]string{left.ID, right.ID})
			}
		}
	}
	return pairs
}

func intervalOverlap(leftMin, leftMax, rightMin, rightMax int) bool {
	return leftMin <= rightMax && rightMin <= leftMax
}
