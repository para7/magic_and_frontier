package skills

import (
	"testing"
	"time"
)

func TestValidateSaveSuccess(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:          "skill_1",
		Name:        " Slash ",
		Description: " desc ",
		Script:      "say slash",
	}, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success, got %+v", result)
	}
	if result.Entry.Name != "Slash" || result.Entry.Description != "desc" {
		t.Fatalf("entry = %#v", result.Entry)
	}
}

func TestValidateSaveErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{ID: "bad", Script: " "}, now)
	if result.OK {
		t.Fatalf("expected validation error")
	}
	if result.FieldErrors["id"] == "" || result.FieldErrors["script"] == "" {
		t.Fatalf("fieldErrors = %#v", result.FieldErrors)
	}
}
