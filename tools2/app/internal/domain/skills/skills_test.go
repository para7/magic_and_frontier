package skills

import (
	"testing"
	"time"
)

func TestValidateSaveSuccessCases(t *testing.T) {
	itemIDs := map[string]struct{}{"00000000-0000-4000-8000-000000000002": {}}
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name  string
		input SaveInput
	}{
		{
			name: "valid and trimmed",
			input: SaveInput{
				ID:     "00000000-0000-4000-8000-000000000001",
				Name:   " Slash ",
				Script: " say x ",
				ItemID: "00000000-0000-4000-8000-000000000002",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ValidateSave(tt.input, itemIDs, now)
			if !res.OK || res.Entry == nil {
				t.Fatalf("expected success, got %+v", res)
			}
			if res.Entry.Name != "Slash" {
				t.Fatalf("name = %q", res.Entry.Name)
			}
			if res.Entry.Script != "say x" {
				t.Fatalf("script = %q", res.Entry.Script)
			}
		})
	}
}

func TestValidateSaveValidationErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		input     SaveInput
		itemIDs   map[string]struct{}
		wantField string
	}{
		{
			name: "missing referenced item",
			input: SaveInput{
				ID:     "00000000-0000-4000-8000-000000000001",
				Name:   "Slash",
				Script: "say x",
				ItemID: "00000000-0000-4000-8000-000000000099",
			},
			itemIDs:   map[string]struct{}{},
			wantField: "itemId",
		},
		{
			name: "name whitespace only",
			input: SaveInput{
				ID:     "00000000-0000-4000-8000-000000000001",
				Name:   " \n ",
				Script: "say x",
				ItemID: "00000000-0000-4000-8000-000000000002",
			},
			itemIDs:   map[string]struct{}{"00000000-0000-4000-8000-000000000002": {}},
			wantField: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ValidateSave(tt.input, tt.itemIDs, now)
			if res.OK {
				t.Fatalf("expected validation error")
			}
			if res.FieldErrors[tt.wantField] == "" {
				t.Fatalf("expected %s field error, got %#v", tt.wantField, res.FieldErrors)
			}
		})
	}
}
