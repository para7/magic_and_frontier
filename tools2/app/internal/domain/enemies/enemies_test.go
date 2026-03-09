package enemies

import (
	"testing"
	"time"
)

func TestValidateSaveSuccessCases(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	enemySkillIDs := map[string]struct{}{
		"00000000-0000-4000-8000-000000000011": {},
		"00000000-0000-4000-8000-000000000012": {},
	}
	itemIDs := map[string]struct{}{"00000000-0000-4000-8000-000000000022": {}}
	min := 1.0
	max := 2.0
	tests := []struct {
		name  string
		input SaveInput
	}{
		{
			name: "minimum valid ranges and duplicate skill ids normalize",
			input: SaveInput{
				ID:            "00000000-0000-4000-8000-000000000001",
				Name:          " Zombie ",
				HP:            1,
				DropTableID:   " drop-1 ",
				EnemySkillIDs: []string{"00000000-0000-4000-8000-000000000011", " 00000000-0000-4000-8000-000000000011 ", "00000000-0000-4000-8000-000000000012"},
				SpawnRule:     SpawnRule{Origin: Vec3{X: 0, Y: 64, Z: 0}, Distance: Distance{Min: 0, Max: 10}},
				DropTable: []DropRef{{
					Kind:     " item ",
					RefID:    "00000000-0000-4000-8000-000000000022",
					Weight:   10,
					CountMin: &min,
					CountMax: &max,
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ValidateSave(tt.input, enemySkillIDs, itemIDs, map[string]struct{}{}, now)
			if !res.OK || res.Entry == nil {
				t.Fatalf("expected success: %+v", res)
			}
			if res.Entry.Name != "Zombie" {
				t.Fatalf("name = %q", res.Entry.Name)
			}
			if res.Entry.DropTableID != "drop-1" {
				t.Fatalf("dropTableId = %q", res.Entry.DropTableID)
			}
			if len(res.Entry.EnemySkillIDs) != 2 {
				t.Fatalf("enemySkillIds = %#v", res.Entry.EnemySkillIDs)
			}
			if len(res.Entry.DropTable) != 1 || res.Entry.DropTable[0].Kind != "item" {
				t.Fatalf("dropTable = %#v", res.Entry.DropTable)
			}
		})
	}
}

func TestValidateSaveValidationErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	enemySkillIDs := map[string]struct{}{"00000000-0000-4000-8000-000000000011": {}}
	itemIDs := map[string]struct{}{"00000000-0000-4000-8000-000000000022": {}}
	highMin := 10.0
	lowMax := 2.0
	axisHigh := 10.0
	axisLow := 2.0
	tests := []struct {
		name      string
		input     SaveInput
		wantField string
	}{
		{
			name: "missing enemy skill reference",
			input: SaveInput{
				ID:            "00000000-0000-4000-8000-000000000001",
				Name:          "Zombie",
				HP:            20,
				DropTableID:   "drop-1",
				EnemySkillIDs: []string{"00000000-0000-4000-8000-000000000099"},
				SpawnRule:     SpawnRule{Origin: Vec3{X: 0, Y: 64, Z: 0}, Distance: Distance{Min: 1, Max: 10}},
			},
			wantField: "enemySkillIds.0",
		},
		{
			name: "missing item drop reference",
			input: SaveInput{
				ID:            "00000000-0000-4000-8000-000000000001",
				Name:          "Zombie",
				HP:            20,
				DropTableID:   "drop-1",
				EnemySkillIDs: []string{"00000000-0000-4000-8000-000000000011"},
				SpawnRule:     SpawnRule{Origin: Vec3{X: 0, Y: 64, Z: 0}, Distance: Distance{Min: 1, Max: 10}},
				DropTable:     []DropRef{{Kind: "item", RefID: "00000000-0000-4000-8000-000000000099", Weight: 10}},
			},
			wantField: "dropTable.0.refId",
		},
		{
			name: "drop table count reversed",
			input: SaveInput{
				ID:            "00000000-0000-4000-8000-000000000001",
				Name:          "Zombie",
				HP:            20,
				DropTableID:   "drop-1",
				EnemySkillIDs: []string{"00000000-0000-4000-8000-000000000011"},
				SpawnRule:     SpawnRule{Origin: Vec3{X: 0, Y: 64, Z: 0}, Distance: Distance{Min: 1, Max: 10}},
				DropTable: []DropRef{{
					Kind:     "item",
					RefID:    "00000000-0000-4000-8000-000000000022",
					Weight:   10,
					CountMin: &highMin,
					CountMax: &lowMax,
				}},
			},
			wantField: "dropTable.0.countMin",
		},
		{
			name: "distance reversed",
			input: SaveInput{
				ID:            "00000000-0000-4000-8000-000000000001",
				Name:          "Zombie",
				HP:            20,
				DropTableID:   "drop-1",
				EnemySkillIDs: []string{"00000000-0000-4000-8000-000000000011"},
				SpawnRule:     SpawnRule{Origin: Vec3{X: 0, Y: 64, Z: 0}, Distance: Distance{Min: 10, Max: 1}},
			},
			wantField: "spawnRule.distance.min",
		},
		{
			name: "axis bounds reversed",
			input: SaveInput{
				ID:            "00000000-0000-4000-8000-000000000001",
				Name:          "Zombie",
				HP:            20,
				DropTableID:   "drop-1",
				EnemySkillIDs: []string{"00000000-0000-4000-8000-000000000011"},
				SpawnRule: SpawnRule{
					Origin:     Vec3{X: 0, Y: 64, Z: 0},
					Distance:   Distance{Min: 1, Max: 10},
					AxisBounds: &AxisBounds{XMin: &axisHigh, XMax: &axisLow},
				},
			},
			wantField: "spawnRule.axisBounds.xMin",
		},
		{
			name: "hp below minimum",
			input: SaveInput{
				ID:            "00000000-0000-4000-8000-000000000001",
				Name:          "Zombie",
				HP:            0,
				DropTableID:   "drop-1",
				EnemySkillIDs: []string{"00000000-0000-4000-8000-000000000011"},
				SpawnRule:     SpawnRule{Origin: Vec3{X: 0, Y: 64, Z: 0}, Distance: Distance{Min: 1, Max: 10}},
			},
			wantField: "hp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ValidateSave(tt.input, enemySkillIDs, itemIDs, map[string]struct{}{}, now)
			if res.OK {
				t.Fatalf("expected validation error")
			}
			if res.FieldErrors[tt.wantField] == "" {
				t.Fatalf("expected %s validation error, got %#v", tt.wantField, res.FieldErrors)
			}
		})
	}
}
