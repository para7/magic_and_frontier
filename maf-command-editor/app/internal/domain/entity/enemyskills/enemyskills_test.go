package enemyskills

import (
	"testing"
	"time"
)

func TestValidateSaveSuccess(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:          "roar-main",
		Name:        " Roar ",
		Description: " desc ",
		Script:      "say roar",
	}, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success, got %+v", result)
	}
	if result.Entry.Name != "Roar" || result.Entry.Description != "desc" {
		t.Fatalf("entry = %#v", result.Entry)
	}
}

func TestValidateSaveErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{ID: "", Script: " "}, now)
	if result.OK {
		t.Fatalf("expected validation error")
	}
	if result.FieldErrors["id"] == "" || result.FieldErrors["script"] == "" {
		t.Fatalf("fieldErrors = %#v", result.FieldErrors)
	}
}
