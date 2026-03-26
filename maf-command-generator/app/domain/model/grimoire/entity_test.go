package grimoire

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func validGrimoire() Grimoire {
	return Grimoire{
		ID:          "grimoire_1",
		CastID:      1,
		CastTime:    20,
		MPCost:      5,
		Script:      "function maf:test",
		Title:       "Firebolt",
		Description: "desc",
		UpdatedAt:   "2026-03-27T00:00:00Z",
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
			name: "id ng whitespace only",
			patch: func(g *Grimoire) {
				g.ID = " \n "
			},
			wantErrField: "id",
		},
		{
			name: "castid ok",
			patch: func(g *Grimoire) {
				g.CastID = 2
			},
		},
		{
			name: "castid ng zero",
			patch: func(g *Grimoire) {
				g.CastID = 0
			},
			wantErrField: "castid",
		},
		{
			name: "castid ng negative",
			patch: func(g *Grimoire) {
				g.CastID = -1
			},
			wantErrField: "castid",
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
				g.Script = "function maf:ok"
			},
		},
		{
			name: "script ng whitespace only",
			patch: func(g *Grimoire) {
				g.Script = " \n \t "
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
		{
			name: "updatedAt ok empty",
			patch: func(g *Grimoire) {
				g.UpdatedAt = ""
			},
		},
		{
			name: "updatedAt ok iso8601",
			patch: func(g *Grimoire) {
				g.UpdatedAt = "2026-03-27T12:34:56Z"
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
