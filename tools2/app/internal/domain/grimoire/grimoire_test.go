package grimoire

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestValidateSaveSuccessCases(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	input := SaveInput{
		ID:          "grimoire_1",
		CastID:      1,
		CastTime:    20,
		MPCost:      5,
		Script:      " function maf:test ",
		Title:       " T ",
		Description: " desc ",
	}

	result := ValidateSave(input, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success, got %+v", result)
	}
	if result.Entry.Title != "T" {
		t.Fatalf("title = %q", result.Entry.Title)
	}
	if result.Entry.CastTime != 20 || result.Entry.MPCost != 5 {
		t.Fatalf("entry = %#v", result.Entry)
	}
}

func TestValidateSaveValidationErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		input     SaveInput
		wantField string
	}{
		{name: "invalid id", input: SaveInput{ID: "bad", CastID: 1, Script: "function maf:test", Title: "T"}, wantField: "id"},
		{name: "cast id below minimum", input: SaveInput{ID: "grimoire_1", CastID: 0, Script: "function maf:test", Title: "T"}, wantField: "castid"},
		{name: "title whitespace only", input: SaveInput{ID: "grimoire_1", CastID: 1, Script: "function maf:test", Title: "   "}, wantField: "title"},
		{name: "script whitespace only", input: SaveInput{ID: "grimoire_1", CastID: 1, Script: " \n ", Title: "T"}, wantField: "script"},
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
	state := GrimoireState{Entries: []GrimoireEntry{{ID: "grimoire_1"}}}
	raw, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"entries"`) {
		t.Fatalf("json shape mismatch: %s", raw)
	}
}
