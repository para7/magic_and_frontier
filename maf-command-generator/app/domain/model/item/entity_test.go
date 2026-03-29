package item

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func validItem() Item {
	return Item{
		ID:        "item_1",
		ItemID:    "minecraft:diamond_sword",
		UpdatedAt: "2026-03-29T00:00:00Z",
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

func TestItemValidateStructAllValid(t *testing.T) {
	entity := &ItemEntity{}
	errs := entity.ValidateStruct(validItem())
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %#v", errs)
	}
}

func TestItemValidateStructPerField(t *testing.T) {
	entity := &ItemEntity{}

	tests := []struct {
		name         string
		patch        func(*Item)
		wantErrField string
	}{
		{name: "id ok", patch: func(it *Item) { it.ID = "ok" }},
		{name: "id ng empty", patch: func(it *Item) { it.ID = "  " }, wantErrField: "id"},
		{name: "itemId ok", patch: func(it *Item) { it.ItemID = "minecraft:stone" }},
		{name: "itemId ng empty", patch: func(it *Item) { it.ItemID = " " }, wantErrField: "itemId"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := validItem()
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

func TestItemComponentsEnchantmentError(t *testing.T) {
	_, errMsg := BuildItemComponents(Item{
		ID:           "item_1",
		ItemID:       "minecraft:sword",
		Enchantments: "bad_format",
	})
	if errMsg == "" {
		t.Fatal("expected enchantment format error, got none")
	}
}

func TestItemValidateAllDetectsDuplicateID(t *testing.T) {
	entity := &ItemEntity{
		data: []Item{
			validItem(),
			validItem(),
		},
	}

	allErrs := entity.ValidateAll(testDBMaster{})
	for _, recordErrs := range allErrs {
		for _, err := range recordErrs {
			if err.ID == "item_1" && err.Field == "id" && err.Tag == "unique" {
				return
			}
		}
	}
	t.Fatalf("expected duplicate id error, got %#v", allErrs)
}
