package treasures

import (
	"testing"
	"time"
)

func TestValidateSaveHappyPath(t *testing.T) {
	itemIDs := map[string]struct{}{"00000000-0000-4000-8000-000000000010": {}}
	res := ValidateSave(SaveInput{
		ID:        "00000000-0000-4000-8000-000000000001",
		Name:      "Chest",
		LootPools: []DropRef{{Kind: "item", RefID: "00000000-0000-4000-8000-000000000010", Weight: 10}},
	}, itemIDs, map[string]struct{}{}, time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC))
	if !res.OK || res.Entry == nil {
		t.Fatalf("expected success")
	}
}
