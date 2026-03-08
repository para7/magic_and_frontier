package skills

import (
	"testing"
	"time"
)

func TestValidateSaveHappyPath(t *testing.T) {
	itemIDs := map[string]struct{}{"00000000-0000-4000-8000-000000000002": {}}
	res := ValidateSave(SaveInput{
		ID:     "00000000-0000-4000-8000-000000000001",
		Name:   "Slash",
		Script: "say x",
		ItemID: "00000000-0000-4000-8000-000000000002",
	}, itemIDs, time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC))
	if !res.OK || res.Entry == nil {
		t.Fatalf("expected success")
	}
}
