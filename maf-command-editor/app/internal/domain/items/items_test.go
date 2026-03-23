package items

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestValidateSaveSuccessCases(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:     "stone-item",
		ItemID: " minecraft:stone ",
	}, map[string]struct{}{}, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success, got %+v", result)
	}
	if result.Entry.ItemID != "minecraft:stone" {
		t.Fatalf("itemId = %q", result.Entry.ItemID)
	}
	if !strings.Contains(result.Entry.NBT, `Count:1b`) {
		t.Fatalf("nbt mismatch: %s", result.Entry.NBT)
	}
}

func TestValidateSaveValidationErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		input     SaveInput
		skillIDs  map[string]struct{}
		wantField string
	}{
		{name: "empty id", input: SaveInput{ID: "  ", ItemID: "minecraft:stone"}, skillIDs: map[string]struct{}{}, wantField: "id"},
		{name: "invalid enchantment line", input: SaveInput{ID: "items_1", ItemID: "minecraft:stone", Enchantments: "minecraft:sharpness"}, skillIDs: map[string]struct{}{}, wantField: "enchantments"},
		{name: "missing skill", input: SaveInput{ID: "items_1", ItemID: "minecraft:stone", SkillID: "skill_2"}, skillIDs: map[string]struct{}{}, wantField: "skillId"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSave(tt.input, tt.skillIDs, now)
			if result.OK {
				t.Fatalf("expected validation error")
			}
			if result.FieldErrors[tt.wantField] == "" {
				t.Fatalf("expected %s field error, got %#v", tt.wantField, result.FieldErrors)
			}
		})
	}
}

func TestStateJSONShape(t *testing.T) {
	state := ItemState{Items: []ItemEntry{{ID: "items_1"}}}
	raw, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"items"`) {
		t.Fatalf("json shape mismatch: %s", raw)
	}
}
