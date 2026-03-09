package grimoire

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
		wantTitle string
	}{
		{
			name: "minimum cast id and trimmed fields",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", CastID: 1, Script: " function maf:test ", Title: " T ",
			},
			wantTitle: "T",
		},
		{
			name: "valid description preserved after trim",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", CastID: 2, Script: " say x ", Title: " Book ", Description: " desc ",
			},
			wantTitle: "Book",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSave(tt.input, now)
			if !result.OK || result.Entry == nil {
				t.Fatalf("expected success, got %+v", result)
			}
			if result.Entry.Title != tt.wantTitle {
				t.Fatalf("title = %q", result.Entry.Title)
			}
			if result.Entry.Script != strings.TrimSpace(tt.input.Script) {
				t.Fatalf("script = %q", result.Entry.Script)
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
			name: "cast id below minimum",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", CastID: 0, Script: "function maf:test", Title: "T",
			},
			wantField: "castid",
		},
		{
			name: "title whitespace only",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", CastID: 1, Script: "function maf:test", Title: "   ",
			},
			wantField: "title",
		},
		{
			name: "script whitespace only",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", CastID: 1, Script: " \n ", Title: "T",
			},
			wantField: "script",
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
	state := GrimoireState{Entries: []GrimoireEntry{{ID: "x"}}}
	raw, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"entries"`) {
		t.Fatalf("json shape mismatch: %s", raw)
	}
}
