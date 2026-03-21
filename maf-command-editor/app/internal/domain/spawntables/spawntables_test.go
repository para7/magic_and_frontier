package spawntables

import (
	"testing"
	"time"
)

func TestValidateSaveSuccess(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:            "spawntable_1",
		SourceMobType: "minecraft:zombie",
		Dimension:     "minecraft:overworld",
		MinX:          0,
		MaxX:          100,
		MinY:          -64,
		MaxY:          320,
		MinZ:          0,
		MaxZ:          100,
		BaseMobWeight: 8000,
		Replacements: []ReplacementEntry{
			{EnemyID: "enemy_1", Weight: 2000},
		},
	}, map[string]struct{}{"enemy_1": {}}, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success, got %+v", result)
	}
	if result.Entry.ID != "spawntable_1" {
		t.Fatalf("entry = %#v", result.Entry)
	}
}

func TestValidateSaveErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:            "bad",
		SourceMobType: "zombie",
		Dimension:     "minecraft:overworld",
		MinX:          100,
		MaxX:          0,
		MinY:          0,
		MaxY:          10,
		MinZ:          0,
		MaxZ:          10,
		BaseMobWeight: 0,
		Replacements: []ReplacementEntry{
			{EnemyID: "enemy_999", Weight: 1},
		},
	}, map[string]struct{}{}, now)
	if result.OK {
		t.Fatalf("expected validation error")
	}
	if result.FieldErrors["id"] == "" || result.FieldErrors["sourceMobType"] == "" || result.FieldErrors["minX"] == "" || result.FieldErrors["replacements.0.enemyId"] == "" {
		t.Fatalf("fieldErrors = %#v", result.FieldErrors)
	}
}

func TestValidateSaveRejectsNonPositiveReplacementWeight(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	enemyIDs := map[string]struct{}{"enemy_1": {}}
	cases := []struct {
		name   string
		weight int
	}{
		{name: "zero", weight: 0},
		{name: "negative", weight: -10},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateSave(SaveInput{
				ID:            "spawntable_1",
				SourceMobType: "minecraft:zombie",
				Dimension:     "minecraft:overworld",
				MinX:          0,
				MaxX:          10,
				MinY:          0,
				MaxY:          10,
				MinZ:          0,
				MaxZ:          10,
				BaseMobWeight: 100,
				Replacements: []ReplacementEntry{
					{EnemyID: "enemy_1", Weight: tc.weight},
				},
			}, enemyIDs, now)
			if result.OK {
				t.Fatalf("expected validation error")
			}
			if result.FieldErrors["replacements.0.weight"] == "" {
				t.Fatalf("fieldErrors = %#v", result.FieldErrors)
			}
		})
	}
}

func TestRangesOverlap(t *testing.T) {
	left := SpawnTableEntry{MinX: 0, MaxX: 10, MinY: 0, MaxY: 10, MinZ: 0, MaxZ: 10}
	right := SpawnTableEntry{MinX: 10, MaxX: 20, MinY: 5, MaxY: 15, MinZ: 5, MaxZ: 15}
	if !RangesOverlap(left, right) {
		t.Fatalf("expected overlap")
	}
}

func TestFirstOverlap(t *testing.T) {
	entries := []SpawnTableEntry{
		{ID: "spawntable_1", SourceMobType: "minecraft:zombie", Dimension: "minecraft:overworld", MinX: 0, MaxX: 10, MinY: 0, MaxY: 10, MinZ: 0, MaxZ: 10},
	}
	candidate := SpawnTableEntry{ID: "spawntable_2", SourceMobType: "minecraft:zombie", Dimension: "minecraft:overworld", MinX: 5, MaxX: 15, MinY: 5, MaxY: 15, MinZ: 5, MaxZ: 15}
	conflictID, ok := FirstOverlap(entries, candidate)
	if !ok || conflictID != "spawntable_1" {
		t.Fatalf("conflictID=%q ok=%v", conflictID, ok)
	}
}
