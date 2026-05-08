package active

import (
	"encoding/json"
	"testing"

	model "maf_command_editor/app/domain/model"
)

func boolPtr(b bool) *bool { return &b }

func validActive() Active {
	return Active{
		ID:          "active_1",
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

func (testDBMaster) HasItem(string) bool   { return true }
func (testDBMaster) HasActive(string) bool { return true }
func (testDBMaster) GetActive(string) (model.ActiveSnapshot, bool) {
	return model.ActiveSnapshot{ID: "active_1", LootEnable: true}, true
}
func (testDBMaster) HasPassive(string) bool { return true }
func (testDBMaster) GetPassive(string) (model.PassiveSnapshot, bool) {
	v := true
	return model.PassiveSnapshot{ID: "passive_1", GenerateGrimoire: &v}, true
}
func (testDBMaster) HasBow(string) bool                { return true }
func (testDBMaster) HasEnemySkill(string) bool         { return true }
func (testDBMaster) HasEnemy(string) bool              { return true }
func (testDBMaster) HasSpawnTable(string) bool         { return true }
func (testDBMaster) HasMinecraftLootTable(string) bool { return true }

func TestActiveValidateStructAllValid(t *testing.T) {
	entity := &ActiveEntity{}

	errs := entity.ValidateStruct(validActive())
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %#v", errs)
	}
}

func TestActiveValidateStructPerFieldOKNG(t *testing.T) {
	entity := &ActiveEntity{}

	tests := []struct {
		name         string
		patch        func(*Active)
		wantErrField string
	}{
		{
			name: "id ok",
			patch: func(g *Active) {
				g.ID = "active_ok"
			},
		},
		{
			name: "id ok hyphen",
			patch: func(g *Active) {
				g.ID = "active-ok"
			},
		},
		{
			name: "id ng whitespace only",
			patch: func(g *Active) {
				g.ID = " \n "
			},
			wantErrField: "id",
		},
		{
			name: "id ng uppercase",
			patch: func(g *Active) {
				g.ID = "Active_1"
			},
			wantErrField: "id",
		},
		{
			name: "id ng with space",
			patch: func(g *Active) {
				g.ID = "fire bolt"
			},
			wantErrField: "id",
		},
		{
			name: "id ng colon",
			patch: func(g *Active) {
				g.ID = "foo:bar"
			},
			wantErrField: "id",
		},
		{
			name: "id ng slash",
			patch: func(g *Active) {
				g.ID = "foo/bar"
			},
			wantErrField: "id",
		},
		{
			name: "id ng dot",
			patch: func(g *Active) {
				g.ID = "foo.bar"
			},
			wantErrField: "id",
		},
		{
			name: "castTime ok lower bound",
			patch: func(g *Active) {
				g.CastTime = 0
			},
		},
		{
			name: "castTime ok upper bound",
			patch: func(g *Active) {
				g.CastTime = 12000
			},
		},
		{
			name: "castTime ng below lower bound",
			patch: func(g *Active) {
				g.CastTime = -1
			},
			wantErrField: "castTime",
		},
		{
			name: "castTime ng above upper bound",
			patch: func(g *Active) {
				g.CastTime = 12001
			},
			wantErrField: "castTime",
		},
		{
			name: "mpCost ok lower bound",
			patch: func(g *Active) {
				g.MPCost = 0
			},
		},
		{
			name: "mpCost ok upper bound",
			patch: func(g *Active) {
				g.MPCost = 1000000
			},
		},
		{
			name: "mpCost ng below lower bound",
			patch: func(g *Active) {
				g.MPCost = -1
			},
			wantErrField: "mpCost",
		},
		{
			name: "mpCost ng above upper bound",
			patch: func(g *Active) {
				g.MPCost = 1000001
			},
			wantErrField: "mpCost",
		},
		{
			name: "script ok",
			patch: func(g *Active) {
				g.Script = []string{"function maf:ok"}
			},
		},
		{
			name: "script ng empty",
			patch: func(g *Active) {
				g.Script = []string{}
			},
			wantErrField: "script",
		},
		{
			name: "title ok",
			patch: func(g *Active) {
				g.Title = "A"
			},
		},
		{
			name: "title ng whitespace only",
			patch: func(g *Active) {
				g.Title = "  "
			},
			wantErrField: "title",
		},
		{
			name: "description ok empty",
			patch: func(g *Active) {
				g.Description = ""
			},
		},
		{
			name: "description ok text",
			patch: func(g *Active) {
				g.Description = "some description"
			},
		},
		{
			name: "loot_enable ok true",
			patch: func(g *Active) {
				g.LootEnable = boolPtr(true)
			},
		},
		{
			name: "loot_enable ok false",
			patch: func(g *Active) {
				g.LootEnable = boolPtr(false)
			},
		},
		{
			name: "loot_enable ok nil defaults true",
			patch: func(g *Active) {
				g.LootEnable = nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := validActive()
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

func TestActiveValidateAllDetectsDuplicateID(t *testing.T) {
	entity := &ActiveEntity{
		data: []Active{
			validActive(),
			validActive(),
		},
	}

	allErrs := entity.ValidateAll(testDBMaster{})
	for _, recordErrs := range allErrs {
		for _, err := range recordErrs {
			if err.ID == "active_1" && err.Field == "id" && err.Tag == "unique" {
				return
			}
		}
	}
	t.Fatalf("expected duplicate id error, got %#v", allErrs)
}

func TestActiveUnmarshalJSONDefaultsLootEnableToTrue(t *testing.T) {
	var g Active
	if err := json.Unmarshal([]byte(`{"id":"g1"}`), &g); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if g.LootEnable == nil || !*g.LootEnable {
		t.Fatalf("expected loot_enable to default true, got %#v", g.LootEnable)
	}
}

func TestActiveUnmarshalJSONKeepsExplicitLootEnableFalse(t *testing.T) {
	var g Active
	if err := json.Unmarshal([]byte(`{"id":"g1","loot_enable":false}`), &g); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if g.LootEnable == nil || *g.LootEnable {
		t.Fatalf("expected loot_enable to stay false, got %#v", g.LootEnable)
	}
}
