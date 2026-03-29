package passive

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func validPassive() Passive {
	return Passive{
		ID:          "passive_1",
		Name:        "斬撃スキル",
		SkillType:   "sword",
		Description: "剣スキル",
		Script:      "function maf:skill/test",
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

func TestPassiveValidateStructAllValid(t *testing.T) {
	entity := &PassiveEntity{}
	errs := entity.ValidateStruct(validPassive())
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %#v", errs)
	}
}

func TestPassiveValidateStructPerField(t *testing.T) {
	entity := &PassiveEntity{}

	tests := []struct {
		name         string
		patch        func(*Passive)
		wantErrField string
	}{
		{name: "id ok", patch: func(p *Passive) { p.ID = "ok" }},
		{name: "id ng empty", patch: func(p *Passive) { p.ID = " " }, wantErrField: "id"},
		{name: "name ok empty", patch: func(p *Passive) { p.Name = "" }},
		{name: "name ok max", patch: func(p *Passive) { p.Name = string(make([]rune, 80)) }},
		{name: "name ng over max", patch: func(p *Passive) { p.Name = string(make([]rune, 81)) }, wantErrField: "name"},
		{name: "skilltype ok sword", patch: func(p *Passive) { p.SkillType = "sword" }},
		{name: "skilltype ok bow", patch: func(p *Passive) { p.SkillType = "bow" }},
		{name: "skilltype ok axe", patch: func(p *Passive) { p.SkillType = "axe" }},
		{name: "skilltype ng invalid", patch: func(p *Passive) { p.SkillType = "staff" }, wantErrField: "skilltype"},
		{name: "skilltype ng empty", patch: func(p *Passive) { p.SkillType = "" }, wantErrField: "skilltype"},
		{name: "script ok", patch: func(p *Passive) { p.Script = "function maf:test" }},
		{name: "script ng empty", patch: func(p *Passive) { p.Script = "  " }, wantErrField: "script"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := validPassive()
			tt.patch(&p)
			errs := entity.ValidateStruct(p)
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

func TestPassiveValidateAllDetectsDuplicateID(t *testing.T) {
	entity := &PassiveEntity{
		data: []Passive{
			validPassive(),
			validPassive(),
		},
	}

	allErrs := entity.ValidateAll(testDBMaster{})
	for _, recordErrs := range allErrs {
		for _, err := range recordErrs {
			if err.ID == "passive_1" && err.Field == "id" && err.Tag == "unique" {
				return
			}
		}
	}
	t.Fatalf("expected duplicate id error, got %#v", allErrs)
}
