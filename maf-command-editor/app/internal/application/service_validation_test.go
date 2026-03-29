package application

import (
	"testing"
)

func TestValidateCheckedInSavedata(t *testing.T) {
	cfg := repoSavedataConfig(t)
	svc := NewService(cfg, Dependencies{Now: fixedNow})

	report, err := svc.ValidateAll()
	if err != nil {
		t.Fatalf("ValidateAll() error = %v", err)
	}
	if !report.OK {
		t.Fatalf("checked-in savedata validation failed:\n%s", report.String())
	}
}
