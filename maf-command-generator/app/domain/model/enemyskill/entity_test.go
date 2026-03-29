package enemyskill

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func validEnemySkill() EnemySkill {
	return EnemySkill{
		ID:          "eskill_1",
		Name:        "炎の一撃",
		Description: "炎のダメージを与える",
		Script:      "function maf:enemy_skill/test",
		UpdatedAt:   "2026-03-29T00:00:00Z",
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

func TestEnemySkillValidateStructAllValid(t *testing.T) {
	entity := &EnemySkillEntity{}
	errs := entity.ValidateStruct(validEnemySkill())
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %#v", errs)
	}
}

func TestEnemySkillValidateStructPerField(t *testing.T) {
	entity := &EnemySkillEntity{}

	tests := []struct {
		name         string
		patch        func(*EnemySkill)
		wantErrField string
	}{
		{name: "id ok", patch: func(e *EnemySkill) { e.ID = "ok" }},
		{name: "id ng empty", patch: func(e *EnemySkill) { e.ID = "  " }, wantErrField: "id"},
		{name: "name ok empty", patch: func(e *EnemySkill) { e.Name = "" }},
		{name: "name ng over max", patch: func(e *EnemySkill) { e.Name = string(make([]rune, 81)) }, wantErrField: "name"},
		{name: "description ok empty", patch: func(e *EnemySkill) { e.Description = "" }},
		{name: "description ng over max", patch: func(e *EnemySkill) { e.Description = string(make([]rune, 401)) }, wantErrField: "description"},
		{name: "script ok", patch: func(e *EnemySkill) { e.Script = "function maf:test" }},
		{name: "script ng empty", patch: func(e *EnemySkill) { e.Script = "  " }, wantErrField: "script"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validEnemySkill()
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

func TestEnemySkillValidateAllDetectsDuplicateID(t *testing.T) {
	entity := &EnemySkillEntity{
		data: []EnemySkill{
			validEnemySkill(),
			validEnemySkill(),
		},
	}

	allErrs := entity.ValidateAll(testDBMaster{})
	for _, recordErrs := range allErrs {
		for _, err := range recordErrs {
			if err.ID == "eskill_1" && err.Field == "id" && err.Tag == "unique" {
				return
			}
		}
	}
	t.Fatalf("expected duplicate id error, got %#v", allErrs)
}
