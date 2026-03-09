package items

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestValidateSaveSuccessCases(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		input     SaveInput
		wantItem  string
		wantCount string
	}{
		{
			name: "minimum count",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", ItemID: " minecraft:stone ", Count: 1,
			},
			wantItem:  "minecraft:stone",
			wantCount: "Count:1b",
		},
		{
			name: "maximum count",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", ItemID: " minecraft:diamond ", Count: 64,
			},
			wantItem:  "minecraft:diamond",
			wantCount: "Count:64b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSave(tt.input, now)
			if !result.OK || result.Entry == nil {
				t.Fatalf("expected success, got %+v", result)
			}
			if result.Entry.UpdatedAt != now.Format(time.RFC3339) {
				t.Fatalf("updatedAt = %s", result.Entry.UpdatedAt)
			}
			if result.Entry.ItemID != tt.wantItem {
				t.Fatalf("itemId = %q", result.Entry.ItemID)
			}
			if !strings.Contains(result.Entry.NBT, `id:"`+tt.wantItem+`"`) {
				t.Fatalf("nbt should include normalized id, got: %s", result.Entry.NBT)
			}
			if !strings.Contains(result.Entry.NBT, tt.wantCount) {
				t.Fatalf("nbt should include count, got: %s", result.Entry.NBT)
			}
		})
	}
}

func TestValidateSaveValidationErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		input     SaveInput
		wantField string
	}{
		{
			name: "invalid enchantment line",
			input: SaveInput{
				ID:           "00000000-0000-4000-8000-000000000001",
				ItemID:       "minecraft:stone",
				Count:        1,
				Enchantments: "minecraft:sharpness",
			},
			wantField: "enchantments",
		},
		{
			name: "count below minimum",
			input: SaveInput{
				ID:     "00000000-0000-4000-8000-000000000001",
				ItemID: "minecraft:stone",
				Count:  0,
			},
			wantField: "count",
		},
		{
			name: "item id whitespace only",
			input: SaveInput{
				ID:     "00000000-0000-4000-8000-000000000001",
				ItemID: " \n ",
				Count:  1,
			},
			wantField: "itemId",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSave(tt.input, now)
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
	state := ItemState{Items: []ItemEntry{{ID: "x"}}}
	raw, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"items"`) {
		t.Fatalf("json shape mismatch: %s", raw)
	}
}
