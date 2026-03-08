package enemies

import (
	"testing"
	"time"
)

func TestValidateSaveHappyPath(t *testing.T) {
	enemySkillIDs := map[string]struct{}{"00000000-0000-4000-8000-000000000011": {}}
	min := 1.0
	max := 2.0
	res := ValidateSave(SaveInput{
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
			CountMin: &min,
			CountMax: &max,
		}},
	}, enemySkillIDs, map[string]struct{}{"00000000-0000-4000-8000-000000000022": {}}, map[string]struct{}{}, time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC))
	if !res.OK || res.Entry == nil {
		t.Fatalf("expected success: %+v", res)
	}
}

func TestValidateSaveRejectsDropTableCountRange(t *testing.T) {
	enemySkillIDs := map[string]struct{}{"00000000-0000-4000-8000-000000000011": {}}
	min := 10.0
	max := 2.0
	res := ValidateSave(SaveInput{
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
			CountMin: &min,
			CountMax: &max,
		}},
	}, enemySkillIDs, map[string]struct{}{"00000000-0000-4000-8000-000000000022": {}}, map[string]struct{}{}, time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC))
	if res.OK {
		t.Fatalf("expected validation error")
	}
	if res.FieldErrors["dropTable.0.countMin"] == "" {
		t.Fatalf("expected dropTable.0.countMin validation error")
	}
}
