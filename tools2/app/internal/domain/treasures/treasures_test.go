package treasures

import (
	"testing"
	"time"
)

func TestValidateSaveSuccessCases(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	min := 1.0
	max := 2.0
	tests := []struct {
		name        string
		input       SaveInput
		itemIDs     map[string]struct{}
		grimoireIDs map[string]struct{}
	}{
		{
			name: "item loot pool with trimming",
			input: SaveInput{
				ID:   "00000000-0000-4000-8000-000000000001",
				Name: " Chest ",
				LootPools: []DropRef{{
					Kind: " item ", RefID: "00000000-0000-4000-8000-000000000010", Weight: 10, CountMin: &min, CountMax: &max,
				}},
			},
			itemIDs:     map[string]struct{}{"00000000-0000-4000-8000-000000000010": {}},
			grimoireIDs: map[string]struct{}{},
		},
		{
			name: "grimoire loot pool",
			input: SaveInput{
				ID:   "00000000-0000-4000-8000-000000000001",
				Name: "Chest",
				LootPools: []DropRef{{
					Kind: "grimoire", RefID: "00000000-0000-4000-8000-000000000020", Weight: 5,
				}},
			},
			itemIDs:     map[string]struct{}{},
			grimoireIDs: map[string]struct{}{"00000000-0000-4000-8000-000000000020": {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ValidateSave(tt.input, tt.itemIDs, tt.grimoireIDs, now)
			if !res.OK || res.Entry == nil {
				t.Fatalf("expected success, got %+v", res)
			}
			if res.Entry.Name != "Chest" {
				t.Fatalf("name = %q", res.Entry.Name)
			}
			if len(res.Entry.LootPools) != 1 {
				t.Fatalf("lootPools = %#v", res.Entry.LootPools)
			}
		})
	}
}

func TestValidateSaveValidationErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	min := 10.0
	max := 2.0
	tests := []struct {
		name        string
		input       SaveInput
		itemIDs     map[string]struct{}
		grimoireIDs map[string]struct{}
		wantField   string
	}{
		{
			name: "loot pools required",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", Name: "Chest",
			},
			itemIDs:     map[string]struct{}{},
			grimoireIDs: map[string]struct{}{},
			wantField:   "lootPools",
		},
		{
			name: "missing item reference",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", Name: "Chest",
				LootPools: []DropRef{{Kind: "item", RefID: "00000000-0000-4000-8000-000000000010", Weight: 10}},
			},
			itemIDs:     map[string]struct{}{},
			grimoireIDs: map[string]struct{}{},
			wantField:   "lootPools.0.refId",
		},
		{
			name: "count range reversed",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", Name: "Chest",
				LootPools: []DropRef{{Kind: "item", RefID: "00000000-0000-4000-8000-000000000010", Weight: 10, CountMin: &min, CountMax: &max}},
			},
			itemIDs:     map[string]struct{}{"00000000-0000-4000-8000-000000000010": {}},
			grimoireIDs: map[string]struct{}{},
			wantField:   "lootPools.0.countMin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ValidateSave(tt.input, tt.itemIDs, tt.grimoireIDs, now)
			if res.OK {
				t.Fatalf("expected validation error")
			}
			if res.FieldErrors[tt.wantField] == "" {
				t.Fatalf("expected %s field error, got %#v", tt.wantField, res.FieldErrors)
			}
		})
	}
}
