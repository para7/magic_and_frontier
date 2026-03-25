package application

import (
	"strings"
	"testing"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/items"
)

func TestServiceAllocateCastID(t *testing.T) {
	cfg := testConfig(t)
	svc := NewService(cfg, Dependencies{Now: fixedNow})
	castID, err := svc.AllocateCastID()
	if err != nil {
		t.Fatal(err)
	}
	if castID != 1 {
		t.Fatalf("castID = %d", castID)
	}
	nextCastID, err := svc.AllocateCastID()
	if err != nil {
		t.Fatal(err)
	}
	if nextCastID != 2 {
		t.Fatalf("nextCastID = %d", nextCastID)
	}
}

func TestServiceExportDatapackRejectsInvalidSavedata(t *testing.T) {
	cfg := testConfig(t)
	writeJSONFile(t, cfg.ItemStatePath, common.EntryState[items.ItemEntry]{Entries: []items.ItemEntry{
		{ID: "items_1", ItemID: "minecraft:apple", SkillID: "skill_999"},
	}})
	svc := NewService(cfg, Dependencies{Now: fixedNow})
	result := svc.ExportDatapack()
	if result.OK {
		t.Fatalf("expected export failure")
	}
	if result.Code != "VALIDATION_FAILED" {
		t.Fatalf("code = %q", result.Code)
	}
	if !strings.Contains(result.Details, "Referenced skill does not exist.") {
		t.Fatalf("details = %s", result.Details)
	}
}
