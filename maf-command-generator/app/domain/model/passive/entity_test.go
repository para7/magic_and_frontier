package passive

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func validPassive() Passive {
	return Passive{
		ID:          "passive_1",
		Name:        "剣の心得",
		Condition:   "on_sword_hit",
		Slots:       []int{1, 2},
		Description: "剣攻撃時に発動するパッシブ",
		Script:      []string{"function maf:skill/test"},
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
func (testDBMaster) HasBow(string) bool                { return false }
func (testDBMaster) HasEnemySkill(string) bool         { return true }
func (testDBMaster) HasEnemy(string) bool              { return true }
func (testDBMaster) HasSpawnTable(string) bool         { return true }
func (testDBMaster) HasTreasure(string) bool           { return true }
func (testDBMaster) HasLootTable(string) bool          { return true }
func (testDBMaster) HasMinecraftLootTable(string) bool { return true }

// bowIDConflictDBMaster は condition=bow のpassiveとbowのID衝突をテストするためのモック
type bowIDConflictDBMaster struct{ testDBMaster }

func (bowIDConflictDBMaster) HasBow(string) bool { return true }

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
		{name: "id ok underscore", patch: func(p *Passive) { p.ID = "passive_ok" }},
		{name: "id ok hyphen", patch: func(p *Passive) { p.ID = "passive-ok" }},
		{name: "id ng empty", patch: func(p *Passive) { p.ID = " " }, wantErrField: "id"},
		{name: "id ng space", patch: func(p *Passive) { p.ID = "fire bolt" }, wantErrField: "id"},
		{name: "id ng uppercase", patch: func(p *Passive) { p.ID = "PassiveOk" }, wantErrField: "id"},
		{name: "id ng colon", patch: func(p *Passive) { p.ID = "foo:bar" }, wantErrField: "id"},
		{name: "id ng slash", patch: func(p *Passive) { p.ID = "foo/bar" }, wantErrField: "id"},
		{name: "id ng dot", patch: func(p *Passive) { p.ID = "foo.bar" }, wantErrField: "id"},
		{name: "name ok empty", patch: func(p *Passive) { p.Name = "" }},
		{name: "name ok max", patch: func(p *Passive) { p.Name = string(make([]rune, 80)) }},
		{name: "name ng over max", patch: func(p *Passive) { p.Name = string(make([]rune, 81)) }, wantErrField: "name"},
		{name: "condition ok always", patch: func(p *Passive) { p.Condition = "always" }},
		{name: "condition ok sword", patch: func(p *Passive) { p.Condition = "on_sword_hit" }},
		{name: "condition ok bow", patch: func(p *Passive) { p.Condition = "bow" }},
		{name: "condition ng empty", patch: func(p *Passive) { p.Condition = " " }, wantErrField: "condition"},
		{name: "condition ng unknown", patch: func(p *Passive) { p.Condition = "unknown" }, wantErrField: "condition"},
		{name: "slots ok", patch: func(p *Passive) { p.Slots = []int{1, 3} }},
		{name: "slots ng empty", patch: func(p *Passive) { p.Slots = nil }, wantErrField: "slots"},
		{name: "slots ng under", patch: func(p *Passive) { p.Slots = []int{0} }, wantErrField: "slots[0]"},
		{name: "slots ng over", patch: func(p *Passive) { p.Slots = []int{4} }, wantErrField: "slots[0]"},
		{name: "slots ng duplicate", patch: func(p *Passive) { p.Slots = []int{2, 2} }, wantErrField: "slots"},
		{name: "script ok", patch: func(p *Passive) { p.Script = []string{"function maf:test"} }},
		{name: "script ng empty", patch: func(p *Passive) { p.Script = []string{} }, wantErrField: "script"},
		{name: "bow ok nil", patch: func(p *Passive) { p.Bow = nil }},
		{name: "bow ok with life_sub", patch: func(p *Passive) { p.Condition = "bow"; p.Bow = &BowConfig{LifeSub: ptrInt(100)} }},
		{name: "bow ng life_sub under", patch: func(p *Passive) { p.Condition = "bow"; p.Bow = &BowConfig{LifeSub: ptrInt(-1)} }, wantErrField: "life_sub"},
		{name: "bow ng life_sub over", patch: func(p *Passive) { p.Condition = "bow"; p.Bow = &BowConfig{LifeSub: ptrInt(1201)} }, wantErrField: "life_sub"},
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

func TestPassiveValidateRelationBowRequiresBowCondition(t *testing.T) {
	entity := &PassiveEntity{}
	p := validPassive()
	p.Bow = &BowConfig{LifeSub: ptrInt(100)}

	errs := entity.ValidateRelation(p, testDBMaster{})
	if !hasFieldError(errs, "bow") {
		t.Fatalf("expected relation error for bow, got %#v", errs)
	}
}

func TestPassiveValidateRelationBowAllowedWithBowCondition(t *testing.T) {
	entity := &PassiveEntity{}
	p := validPassive()
	p.Condition = "bow"
	p.Bow = &BowConfig{LifeSub: ptrInt(100)}

	errs := entity.ValidateRelation(p, testDBMaster{})
	if len(errs) != 0 {
		t.Fatalf("expected no relation errors, got %#v", errs)
	}
}

func TestPassiveValidateRelationBowIDConflict(t *testing.T) {
	entity := &PassiveEntity{}
	p := validPassive()
	p.Condition = "bow"

	errs := entity.ValidateRelation(p, bowIDConflictDBMaster{})
	if !hasFieldError(errs, "id") {
		t.Fatalf("expected id conflict error, got %#v", errs)
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

func ptrInt(v int) *int {
	return &v
}
