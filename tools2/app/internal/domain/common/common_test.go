package common

import "testing"

type dummy struct{ ID string }
type validatorChild struct {
	CountMin *float64 `json:"countMin" validate:"omitempty,gte=1,lte=64"`
}
type validatorInput struct {
	Name  string           `json:"name" validate:"trimmed_required"`
	Items []validatorChild `json:"items" validate:"dive"`
}

func TestNormalizeText(t *testing.T) {
	got := NormalizeText("  a\r\nb\r  ")
	if got != "a\nb" {
		t.Fatalf("NormalizeText = %q", got)
	}
}

func TestEntryOps(t *testing.T) {
	state := EmptyEntryState[dummy]()
	next, mode := UpsertEntries(state, dummy{ID: "1"}, func(d dummy) string { return d.ID })
	if mode != SaveModeCreated || len(next.Entries) != 1 {
		t.Fatalf("create upsert failed")
	}
	next2, mode2 := UpsertEntries(next, dummy{ID: "1"}, func(d dummy) string { return d.ID })
	if mode2 != SaveModeUpdated || len(next2.Entries) != 1 {
		t.Fatalf("update upsert failed")
	}
	next3, ok := DeleteEntries(next2, "1", func(d dummy) string { return d.ID })
	if !ok || len(next3.Entries) != 0 {
		t.Fatalf("delete failed")
	}
}

func TestValidateStructUsesJSONFieldPaths(t *testing.T) {
	zero := 0.0
	violations := ValidateStruct(validatorInput{
		Name: "   ",
		Items: []validatorChild{{
			CountMin: &zero,
		}},
	})
	errs := ViolationsToFieldErrors(violations, DefaultValidationMessage)

	if errs["name"] != "Required." {
		t.Fatalf("name error = %q", errs["name"])
	}
	if errs["items.0.countMin"] == "" {
		t.Fatalf("expected nested json path error, got: %#v", errs)
	}
}
