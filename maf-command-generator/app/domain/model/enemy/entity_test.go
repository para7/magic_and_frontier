package enemy

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func validEnemy() Enemy {
	return Enemy{
		ID:       "enemy_1",
		MobType:  "minecraft:zombie",
		HP:       20,
		DropMode: "append",
	}
}

func hasFieldError(errs []model.ValidationError, field string) bool {
	for _, err := range errs {
		if err.Field == field {
			return true
		}
	}
	return false
}

type testDBMaster struct{}

func (testDBMaster) HasItem(string) bool               { return true }
func (testDBMaster) HasGrimoire(string) bool           { return true }
func (testDBMaster) HasPassive(string) bool            { return true }
func (testDBMaster) HasEnemySkill(string) bool         { return true }
func (testDBMaster) HasEnemy(string) bool              { return true }
func (testDBMaster) HasSpawnTable(string) bool         { return true }
func (testDBMaster) HasTreasure(string) bool           { return true }
func (testDBMaster) HasLootTable(string) bool          { return true }
func (testDBMaster) HasMinecraftLootTable(string) bool { return true }

func TestEnemyValidateStructAllValid(t *testing.T) {
	entity := &EnemyEntity{}
	errs := entity.ValidateStruct(validEnemy())
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %#v", errs)
	}
}

func TestEnemyValidateStructPerField(t *testing.T) {
	entity := &EnemyEntity{}

	tests := []struct {
		name         string
		patch        func(*Enemy)
		wantErrField string
	}{
		{name: "id ok", patch: func(e *Enemy) { e.ID = "ok" }},
		{name: "id ng empty", patch: func(e *Enemy) { e.ID = "  " }, wantErrField: "id"},
		{name: "mobType ok", patch: func(e *Enemy) { e.MobType = "minecraft:skeleton" }},
		{name: "mobType ng empty", patch: func(e *Enemy) { e.MobType = " " }, wantErrField: "mobType"},
		{name: "hp ok", patch: func(e *Enemy) { e.HP = 1 }},
		{name: "hp ok max", patch: func(e *Enemy) { e.HP = 100000 }},
		{name: "hp ng below min", patch: func(e *Enemy) { e.HP = 0 }, wantErrField: "hp"},
		{name: "hp ng above max", patch: func(e *Enemy) { e.HP = 100001 }, wantErrField: "hp"},
		{name: "dropMode ok append", patch: func(e *Enemy) { e.DropMode = "append" }},
		{name: "dropMode ok replace", patch: func(e *Enemy) { e.DropMode = "replace" }},
		{name: "dropMode ng invalid", patch: func(e *Enemy) { e.DropMode = "add" }, wantErrField: "dropMode"},
		{name: "dropMode ng empty", patch: func(e *Enemy) { e.DropMode = "" }, wantErrField: "dropMode"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validEnemy()
			tt.patch(&e)
			errs := entity.ValidateStruct(e)
			if tt.wantErrField == "" {
				if len(errs) != 0 {
					t.Fatalf("expected no errors, got %#v", errs)
				}
				return
			}
			if !hasFieldError(errs, tt.wantErrField) {
				t.Fatalf("expected error for field %q, got %#v", tt.wantErrField, errs)
			}
		})
	}
}

func TestEnemyValidateAllDetectsDuplicateID(t *testing.T) {
	entity := &EnemyEntity{
		data: []Enemy{
			validEnemy(),
			validEnemy(),
		},
	}

	allErrs := entity.ValidateAll(testDBMaster{})
	for _, recordErrs := range allErrs {
		for _, err := range recordErrs {
			if err.ID == "enemy_1" && err.Field == "id" && err.Tag == "unique" {
				return
			}
		}
	}
	t.Fatalf("expected duplicate id error, got %#v", allErrs)
}
