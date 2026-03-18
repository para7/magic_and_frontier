package loottables

import (
	"testing"
	"time"

	"tools2/app/internal/domain/treasures"
)

func TestValidateSaveSuccess(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:        "loottable_1",
		TablePath: "maf:treasure/test",
		LootPools: []treasures.DropRef{{Kind: "minecraft_item", RefID: "minecraft:apple", Weight: 1}},
	}, map[string]struct{}{}, map[string]struct{}{}, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success, got %+v", result)
	}
	if result.Entry.TablePath != "maf:treasure/test" {
		t.Fatalf("entry = %#v", result.Entry)
	}
}

func TestValidateSaveErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:        "bad",
		TablePath: "bad path",
		LootPools: []treasures.DropRef{{Kind: "item", RefID: "items_9", Weight: 1}},
	}, map[string]struct{}{}, map[string]struct{}{}, now)
	if result.OK {
		t.Fatalf("expected validation error")
	}
	if result.FieldErrors["id"] == "" || result.FieldErrors["tablePath"] == "" || result.FieldErrors["lootPools.0.refId"] == "" {
		t.Fatalf("fieldErrors = %#v", result.FieldErrors)
	}
}

func TestValidateSaveRejectsTraversalTablePath(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:        "loottable_1",
		TablePath: "maf:loot/../escape",
		LootPools: []treasures.DropRef{{Kind: "minecraft_item", RefID: "minecraft:apple", Weight: 1}},
	}, map[string]struct{}{}, map[string]struct{}{}, now)
	if result.OK {
		t.Fatalf("expected validation error")
	}
	if result.FieldErrors["tablePath"] == "" {
		t.Fatalf("fieldErrors = %#v", result.FieldErrors)
	}
}
