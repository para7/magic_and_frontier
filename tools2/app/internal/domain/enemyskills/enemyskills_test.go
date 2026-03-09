package enemyskills

import (
	"testing"
	"time"
)

func TestValidateSaveSuccessCases(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	zero := 0.0
	ten := 10.0
	tests := []struct {
		name        string
		input       SaveInput
		wantTrigger Trigger
	}{
		{
			name: "cooldown omitted",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", Name: "Roar", Script: "say roar",
			},
		},
		{
			name: "zero cooldown and trimmed trigger",
			input: SaveInput{
				ID:       "00000000-0000-4000-8000-000000000001",
				Name:     "Roar",
				Script:   " say roar ",
				Cooldown: &zero,
				Trigger:  " on_spawn ",
			},
			wantTrigger: TriggerOnSpawn,
		},
		{
			name: "positive cooldown",
			input: SaveInput{
				ID:       "00000000-0000-4000-8000-000000000001",
				Name:     "Roar",
				Script:   "say roar",
				Cooldown: &ten,
				Trigger:  "on_hit",
			},
			wantTrigger: TriggerOnHit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ValidateSave(tt.input, now)
			if !res.OK || res.Entry == nil {
				t.Fatalf("expected success, got %+v", res)
			}
			if res.Entry.Script != "say roar" {
				t.Fatalf("script = %q", res.Entry.Script)
			}
			if tt.input.Trigger == "" {
				if res.Entry.Trigger != nil {
					t.Fatalf("expected nil trigger, got %#v", res.Entry.Trigger)
				}
				return
			}
			if res.Entry.Trigger == nil || *res.Entry.Trigger != tt.wantTrigger {
				t.Fatalf("trigger = %#v", res.Entry.Trigger)
			}
		})
	}
}

func TestValidateSaveValidationErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	negative := -1.0
	tests := []struct {
		name      string
		input     SaveInput
		wantField string
	}{
		{
			name: "negative cooldown",
			input: SaveInput{
				ID:       "00000000-0000-4000-8000-000000000001",
				Name:     "Roar",
				Script:   "say roar",
				Cooldown: &negative,
			},
			wantField: "cooldown",
		},
		{
			name: "unknown trigger",
			input: SaveInput{
				ID:      "00000000-0000-4000-8000-000000000001",
				Name:    "Roar",
				Script:  "say roar",
				Trigger: "unknown",
			},
			wantField: "trigger",
		},
		{
			name: "name whitespace only",
			input: SaveInput{
				ID: "00000000-0000-4000-8000-000000000001", Name: " \t ", Script: "say roar",
			},
			wantField: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ValidateSave(tt.input, now)
			if res.OK {
				t.Fatalf("expected validation error")
			}
			if res.FieldErrors[tt.wantField] == "" {
				t.Fatalf("expected %s field error, got %#v", tt.wantField, res.FieldErrors)
			}
		})
	}
}
