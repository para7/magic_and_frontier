package bow

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func validBow() BowPassive {
	return BowPassive{
		ID:        "bow_1",
		Name:      "Bow 1",
		Role:      "test",
		ScriptHit: []string{"say hit"},
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
func (testDBMaster) HasPassive(string) bool            { return false }
func (testDBMaster) HasBow(string) bool                { return true }
func (testDBMaster) HasEnemySkill(string) bool         { return true }
func (testDBMaster) HasEnemy(string) bool              { return true }
func (testDBMaster) HasSpawnTable(string) bool         { return true }
func (testDBMaster) HasTreasure(string) bool           { return true }
func (testDBMaster) HasLootTable(string) bool          { return true }
func (testDBMaster) HasMinecraftLootTable(string) bool { return true }

// passiveIDConflictDBMaster はbowとpassiveのID衝突をテストするためのモック
type passiveIDConflictDBMaster struct{ testDBMaster }

func (passiveIDConflictDBMaster) HasPassive(string) bool { return true }

func TestBowValidateStructAllValid(t *testing.T) {
	entity := &BowEntity{}
	errs := entity.ValidateStruct(validBow())
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %#v", errs)
	}
}

func TestBowValidateStructPerField(t *testing.T) {
	entity := &BowEntity{}

	tests := []struct {
		name         string
		patch        func(*BowPassive)
		wantErrField string
	}{
		{name: "id ok", patch: func(it *BowPassive) { it.ID = "ok" }},
		{name: "id ng empty", patch: func(it *BowPassive) { it.ID = "  " }, wantErrField: "id"},
		{name: "life_sub ok", patch: func(it *BowPassive) { it.LifeSub = ptrInt(100) }},
		{name: "life_sub ng under", patch: func(it *BowPassive) { it.LifeSub = ptrInt(-1) }, wantErrField: "life_sub"},
		{name: "life_sub ng over", patch: func(it *BowPassive) { it.LifeSub = ptrInt(1201) }, wantErrField: "life_sub"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := validBow()
			tt.patch(&it)
			errs := entity.ValidateStruct(it)
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

func TestBowValidateRelationPassiveIDConflict(t *testing.T) {
	entity := &BowEntity{}
	b := validBow()

	errs := entity.ValidateRelation(b, passiveIDConflictDBMaster{})
	if !hasFieldError(errs, "id") {
		t.Fatalf("expected id conflict error, got %#v", errs)
	}
}

func TestBowValidateRelationNoConflict(t *testing.T) {
	entity := &BowEntity{}
	b := validBow()

	errs := entity.ValidateRelation(b, testDBMaster{})
	if len(errs) != 0 {
		t.Fatalf("expected no relation errors, got %#v", errs)
	}
}

func TestBowValidateAllDetectsDuplicateID(t *testing.T) {
	entity := &BowEntity{
		data: []BowPassive{
			validBow(),
			validBow(),
		},
	}

	allErrs := entity.ValidateAll(testDBMaster{})
	for _, recordErrs := range allErrs {
		for _, err := range recordErrs {
			if err.ID == "bow_1" && err.Field == "id" && err.Tag == "unique" {
				return
			}
		}
	}
	t.Fatalf("expected duplicate id error, got %#v", allErrs)
}

func ptrInt(v int) *int {
	return &v
}
