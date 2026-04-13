package loottable

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func ptrFloat(v float64) *float64 { return &v }

func validLootTable() LootTable {
	return LootTable{
		ID:   "lt_1",
		Memo: "テスト用ルートテーブル",
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

func TestLootTableValidateStructAllValid(t *testing.T) {
	entity := &LootTableEntity{}
	errs := entity.ValidateStruct(validLootTable())
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %#v", errs)
	}
}

func TestLootTableValidateStructPerField(t *testing.T) {
	entity := &LootTableEntity{}

	tests := []struct {
		name         string
		patch        func(*LootTable)
		wantErrField string
	}{
		{name: "id ok", patch: func(lt *LootTable) { lt.ID = "ok" }},
		{name: "id ng empty", patch: func(lt *LootTable) { lt.ID = "  " }, wantErrField: "id"},
		{name: "memo ok empty", patch: func(lt *LootTable) { lt.Memo = "" }},
		{name: "memo ng over max", patch: func(lt *LootTable) { lt.Memo = string(make([]rune, 401)) }, wantErrField: "memo"},
		{name: "lootPools ng empty", patch: func(lt *LootTable) { lt.LootPools = nil }, wantErrField: "lootPools"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lt := validLootTable()
			tt.patch(&lt)
			errs := entity.ValidateStruct(lt)
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

func TestLootTableValidateAllDetectsDuplicateID(t *testing.T) {
	entity := &LootTableEntity{
		data: []LootTable{
			validLootTable(),
			validLootTable(),
		},
	}

	allErrs := entity.ValidateAll(testDBMaster{})
	for _, recordErrs := range allErrs {
		for _, err := range recordErrs {
			if err.ID == "lt_1" && err.Field == "id" && err.Tag == "unique" {
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

func TestLootTableValidateRelationRejectsPassiveWithGenerateGrimoireFalse(t *testing.T) {
	entity := &LootTableEntity{}
	slot := 1
	lt := validLootTable()
	lt.LootPools = []model.DropRef{{Kind: "passive", RefID: "passive_1", Slot: &slot, Weight: 1}}
	errs := entity.ValidateRelation(lt, passiveNotGeneratableDBMaster{})
	if !hasFieldError(errs, "lootPools[0].refId") {
		t.Fatalf("expected passive generate_grimoire relation error, got %#v", errs)
	}
}
