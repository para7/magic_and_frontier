package enemyskills

import (
	"testing"
	"time"
)

func TestValidateSaveHappyPath(t *testing.T) {
	cd := 10.0
	res := ValidateSave(SaveInput{
		ID:       "00000000-0000-4000-8000-000000000001",
		Name:     "Roar",
		Script:   "say roar",
		Cooldown: &cd,
		Trigger:  "on_spawn",
	}, time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC))
	if !res.OK || res.Entry == nil {
		t.Fatalf("expected success")
	}
}
