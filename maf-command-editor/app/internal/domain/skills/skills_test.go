package skills

import (
	"testing"
	"time"
)

func TestValidateSaveSuccess(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:          "slash-primary",
		Name:        " Slash ",
		SkillType:   " bow ",
		Description: " desc ",
		Script:      "say slash",
	}, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success, got %+v", result)
	}
	if result.Entry.Name != "Slash" || result.Entry.Description != "desc" || result.Entry.SkillType != "bow" {
		t.Fatalf("entry = %#v", result.Entry)
	}
}

func TestValidateSaveDefaultsSkillType(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{
		ID:     "slash-default",
		Name:   "Slash",
		Script: "say slash",
	}, now)
	if !result.OK || result.Entry == nil {
		t.Fatalf("expected success, got %+v", result)
	}
	if result.Entry.SkillType != "sword" {
		t.Fatalf("entry = %#v", result.Entry)
	}
}

func TestValidateSaveErrors(t *testing.T) {
	now := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	result := ValidateSave(SaveInput{ID: "   ", SkillType: "gun", Script: " "}, now)
	if result.OK {
		t.Fatalf("expected validation error")
	}
	if result.FieldErrors["id"] == "" || result.FieldErrors["script"] == "" || result.FieldErrors["skilltype"] == "" {
		t.Fatalf("fieldErrors = %#v", result.FieldErrors)
	}
}
