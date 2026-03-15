package enemies

import (
	"testing"
	"time"
)

func TestValidateSaveSuccess(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	dropChance := 0.5
	result := ValidateSave(SaveInput{
		ID:      "enemy_1",
		MobType: "minecraft:zombie",
		Name:    " Zombie ",
		HP:      20,
		Equipment: Equipment{
			Mainhand: &EquipmentSlot{Kind: "minecraft_item", RefID: "minecraft:iron_sword", Count: 1, DropChance: &dropChance},
		},
		EnemySkillIDs: []string{"enemyskill_1", "enemyskill_1"},
		DropMode:      "replace",
		Drops:         []DropRef{{Kind: "minecraft_item", RefID: "minecraft:apple", Weight: 1}},
	}, map[string]struct{}{"enemyskill_1": {}}, map[string]struct{}{}, map[string]struct{}{}, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success, got %+v", result)
	}
	if len(result.Entry.EnemySkillIDs) != 1 {
		t.Fatalf("enemySkillIds = %#v", result.Entry.EnemySkillIDs)
	}
	if result.Entry.Equipment.Mainhand == nil {
		t.Fatalf("expected equipment")
	}
}

func TestValidateSaveErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:            "bad",
		MobType:       " ",
		HP:            0,
		EnemySkillIDs: []string{"bad"},
		DropMode:      "bad",
		Drops:         []DropRef{{Kind: "item", RefID: "items_2", Weight: 1}},
	}, map[string]struct{}{}, map[string]struct{}{}, map[string]struct{}{}, now)
	if result.OK {
		t.Fatalf("expected validation error")
	}
	if result.FieldErrors["id"] == "" || result.FieldErrors["mobType"] == "" || result.FieldErrors["enemySkillIds.0"] == "" || result.FieldErrors["drops.0.refId"] == "" {
		t.Fatalf("fieldErrors = %#v", result.FieldErrors)
	}
}
