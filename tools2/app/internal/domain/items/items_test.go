package items

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestValidateSaveHappyPath(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID: "00000000-0000-4000-8000-000000000001", ItemID: "minecraft:stone", Count: 1,
	}, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success")
	}
	if result.Entry.UpdatedAt != now.Format(time.RFC3339) {
		t.Fatalf("updatedAt = %s", result.Entry.UpdatedAt)
	}
	if !strings.Contains(result.Entry.NBT, `id:"minecraft:stone"`) {
		t.Fatalf("nbt should include id, got: %s", result.Entry.NBT)
	}
	if !strings.Contains(result.Entry.NBT, "Count:1b") {
		t.Fatalf("nbt should include count, got: %s", result.Entry.NBT)
	}
}

func TestValidateSaveRejectsInvalidEnchantmentLine(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:           "00000000-0000-4000-8000-000000000001",
		ItemID:       "minecraft:stone",
		Count:        1,
		Enchantments: "minecraft:sharpness",
	}, now)
	if result.OK {
		t.Fatalf("expected validation error")
	}
	if result.FieldErrors["enchantments"] == "" {
		t.Fatalf("expected enchantments field error")
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
