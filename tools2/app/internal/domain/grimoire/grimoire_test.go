package grimoire

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestValidateSaveHappyPath(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID: "00000000-0000-4000-8000-000000000001", CastID: 1, Script: "function maf:test", Title: "T",
	}, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success")
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
