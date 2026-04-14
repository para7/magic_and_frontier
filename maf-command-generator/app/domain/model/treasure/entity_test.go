package treasure

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func ptrFloat(v float64) *float64 { return &v }

func validTreasure() Treasure {
	return Treasure{
		ID:        "treasure_1",
		TablePath: "minecraft:chests/simple_dungeon",
		LootPools: []model.DropRef{
			{Kind: "minecraft_item", RefID: "minecraft:diamond", Weight: 1, CountMin: ptrFloat(1), CountMax: ptrFloat(1)},
		},
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

func TestTreasureValidateStructAllValid(t *testing.T) {
	entity := &TreasureEntity{}
	errs := entity.ValidateStruct(validTreasure())
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %#v", errs)
	}
}

func TestTreasureValidateStructPerField(t *testing.T) {
	entity := &TreasureEntity{}

	tests := []struct {
		name         string
		patch        func(*Treasure)
		wantErrField string
	}{
		{name: "id ok", patch: func(tr *Treasure) { tr.ID = "ok" }},
		{name: "id ng empty", patch: func(tr *Treasure) { tr.ID = " " }, wantErrField: "id"},
		{name: "tablePath ok", patch: func(tr *Treasure) { tr.TablePath = "minecraft:chests/stronghold_corridor" }},
		{name: "tablePath ng empty", patch: func(tr *Treasure) { tr.TablePath = " " }, wantErrField: "tablePath"},
		{name: "lootPools ng empty", patch: func(tr *Treasure) { tr.LootPools = nil }, wantErrField: "lootPools"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := validTreasure()
			tt.patch(&tr)
			errs := entity.ValidateStruct(tr)
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

func TestTreasureValidateAllDetectsDuplicateID(t *testing.T) {
	entity := &TreasureEntity{
		data: []Treasure{
			validTreasure(),
			validTreasure(),
		},
	}

	allErrs := entity.ValidateAll(testDBMaster{})
	for _, recordErrs := range allErrs {
		for _, err := range recordErrs {
			if err.ID == "treasure_1" && err.Field == "id" && err.Tag == "unique" {
				return
			}
		}
	}
	t.Fatalf("expected duplicate id error, got %#v", allErrs)
}

type passiveNotGeneratableDBMaster struct{ testDBMaster }

func (passiveNotGeneratableDBMaster) GetPassive(string) (model.PassiveSnapshot, bool) {
	v := false
	return model.PassiveSnapshot{ID: "passive_1", GenerateGrimoire: &v}, true
}

func TestTreasureValidateRelationRejectsPassiveWithGenerateGrimoireFalse(t *testing.T) {
	entity := &TreasureEntity{}
	slot := 1
	tr := validTreasure()
	tr.LootPools = []model.DropRef{{Kind: "passive", RefID: "passive_1", Slot: &slot, Weight: 1}}
	errs := entity.ValidateRelation(tr, passiveNotGeneratableDBMaster{})
	if !hasFieldError(errs, "lootPools[0].refId") {
		t.Fatalf("expected passive generate_grimoire relation error, got %#v", errs)
	}
}
