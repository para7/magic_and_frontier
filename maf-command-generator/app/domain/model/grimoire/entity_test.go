package grimoire

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func validGrimoire() Grimoire {
	return Grimoire{
		ID:          "grimoire_1",
		CastTime:    20,
		MPCost:      5,
		Script:      []string{"function maf:test"},
		Title:       "Firebolt",
		Description: "desc",
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

func (testDBMaster) HasItem(string) bool     { return true }
func (testDBMaster) HasGrimoire(string) bool { return true }
func (testDBMaster) HasPassive(string) bool  { return true }
func (testDBMaster) GetPassive(string) (model.PassiveSnapshot, bool) {
	v := true
	return model.PassiveSnapshot{ID: "passive_1", GenerateGrimoire: &v}, true
}
func (testDBMaster) HasBow(string) bool                { return true }
func (testDBMaster) HasEnemySkill(string) bool         { return true }
func (testDBMaster) HasEnemy(string) bool              { return true }
func (testDBMaster) HasSpawnTable(string) bool         { return true }
func (testDBMaster) HasTreasure(string) bool           { return true }
func (testDBMaster) HasLootTable(string) bool          { return true }
func (testDBMaster) HasMinecraftLootTable(string) bool { return true }

func TestGrimoireValidateStructAllValid(t *testing.T) {
	entity := &GrimoireEntity{}

	errs := entity.ValidateStruct(validGrimoire())
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %#v", errs)
	}
}

func TestGrimoireValidateStructPerFieldOKNG(t *testing.T) {
	entity := &GrimoireEntity{}

	tests := []struct {
		name         string
		patch        func(*Grimoire)
		wantErrField string
	}{
		{
			name: "id ok",
			patch: func(g *Grimoire) {
				g.ID = "grimoire_ok"
			},
		},
		{
			name: "id ok hyphen",
			patch: func(g *Grimoire) {
				g.ID = "grimoire-ok"
			},
		},
		{
			name: "id ng whitespace only",
			patch: func(g *Grimoire) {
				g.ID = " \n "
			},
			wantErrField: "id",
		},
		{
			name: "id ng uppercase",
			patch: func(g *Grimoire) {
				g.ID = "Grimoire_1"
			},
			wantErrField: "id",
		},
		{
			name: "id ng with space",
			patch: func(g *Grimoire) {
				g.ID = "fire bolt"
			},
			wantErrField: "id",
		},
		{
			name: "id ng colon",
			patch: func(g *Grimoire) {
				g.ID = "foo:bar"
			},
			wantErrField: "id",
		},
		{
			name: "id ng slash",
			patch: func(g *Grimoire) {
				g.ID = "foo/bar"
			},
			wantErrField: "id",
		},
		{
			name: "id ng dot",
			patch: func(g *Grimoire) {
				g.ID = "foo.bar"
			},
			wantErrField: "id",
		},
		{
			name: "castTime ok lower bound",
			patch: func(g *Grimoire) {
				g.CastTime = 0
			},
		},
		{
			name: "castTime ok upper bound",
			patch: func(g *Grimoire) {
				g.CastTime = 12000
			},
		},
		{
			name: "castTime ng below lower bound",
			patch: func(g *Grimoire) {
				g.CastTime = -1
			},
			wantErrField: "castTime",
		},
		{
			name: "castTime ng above upper bound",
			patch: func(g *Grimoire) {
				g.CastTime = 12001
			},
			wantErrField: "castTime",
		},
		{
			name: "mpCost ok lower bound",
			patch: func(g *Grimoire) {
				g.MPCost = 0
			},
		},
		{
			name: "mpCost ok upper bound",
			patch: func(g *Grimoire) {
				g.MPCost = 1000000
			},
		},
		{
			name: "mpCost ng below lower bound",
			patch: func(g *Grimoire) {
				g.MPCost = -1
			},
			wantErrField: "mpCost",
		},
		{
			name: "mpCost ng above upper bound",
			patch: func(g *Grimoire) {
				g.MPCost = 1000001
			},
			wantErrField: "mpCost",
		},
		{
			name: "script ok",
			patch: func(g *Grimoire) {
				g.Script = []string{"function maf:ok"}
			},
		},
		{
			name: "script ng empty",
			patch: func(g *Grimoire) {
				g.Script = []string{}
			},
			wantErrField: "script",
		},
		{
			name: "title ok",
			patch: func(g *Grimoire) {
				g.Title = "A"
			},
		},
		{
			name: "title ng whitespace only",
			patch: func(g *Grimoire) {
				g.Title = "  "
			},
			wantErrField: "title",
		},
		{
			name: "description ok empty",
			patch: func(g *Grimoire) {
				g.Description = ""
			},
		},
		{
			name: "description ok text",
			patch: func(g *Grimoire) {
				g.Description = "some description"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := validGrimoire()
			tt.patch(&g)

			errs := entity.ValidateStruct(g)
			if tt.wantErrField == "" {
				if len(errs) != 0 {
					t.Fatalf("expected no validation errors, got %#v", errs)
				}
				return
			}

			if len(errs) == 0 {
				t.Fatalf("expected validation error for field %q, got none", tt.wantErrField)
			}
			if !hasFieldError(errs, tt.wantErrField) {
				t.Fatalf("expected validation error for field %q, got %#v", tt.wantErrField, errs)
			}
		})
	}
}

func TestGrimoireValidateAllDetectsDuplicateID(t *testing.T) {
	entity := &GrimoireEntity{
		data: []Grimoire{
			validGrimoire(),
			validGrimoire(),
		},
	}

	allErrs := entity.ValidateAll(testDBMaster{})
	for _, recordErrs := range allErrs {
		for _, err := range recordErrs {
			if err.ID == "grimoire_1" && err.Field == "id" && err.Tag == "unique" {
				return
			}
		}
	}
	t.Fatalf("expected duplicate id error, got %#v", allErrs)
}
