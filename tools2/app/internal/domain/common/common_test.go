package common

import (
	"strings"
	"testing"
)

type dummy struct{ ID string }
type validatorChild struct {
	CountMin *float64 `json:"countMin" validate:"omitempty,gte=1,lte=64"`
}
type validatorInput struct {
	Name  string           `json:"name" validate:"trimmed_required"`
	Items []validatorChild `json:"items" validate:"dive"`
}

type validatorRulesInput struct {
	Required string `json:"required" validate:"trimmed_required"`
	Min      string `json:"min" validate:"trimmed_min=2"`
	Max      string `json:"max" validate:"trimmed_max=3"`
	OneOf    string `json:"oneOf" validate:"trimmed_oneof=alpha beta"`
	ID       string `json:"id" validate:"uuid_any"`
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

func TestValidateStructCustomValidators(t *testing.T) {
	tests := []struct {
		name          string
		input         validatorRulesInput
		wantField     string
		wantTag       string
		wantViolation bool
	}{
		{
			name: "all valid after trimming",
			input: validatorRulesInput{
				Required: "  ok  ",
				Min:      "  ab  ",
				Max:      " abc ",
				OneOf:    "  beta  ",
				ID:       " 00000000-0000-4000-8000-000000000001 ",
			},
		},
		{
			name: "trimmed required rejects whitespace only",
			input: validatorRulesInput{
				Required: " \n\t ",
				Min:      "ab",
				Max:      "abc",
				OneOf:    "alpha",
				ID:       "00000000-0000-4000-8000-000000000001",
			},
			wantField:     "required",
			wantTag:       "trimmed_required",
			wantViolation: true,
		},
		{
			name: "trimmed min rejects short normalized text",
			input: validatorRulesInput{
				Required: "ok",
				Min:      " a ",
				Max:      "abc",
				OneOf:    "alpha",
				ID:       "00000000-0000-4000-8000-000000000001",
			},
			wantField:     "min",
			wantTag:       "trimmed_min",
			wantViolation: true,
		},
		{
			name: "trimmed max rejects long normalized text",
			input: validatorRulesInput{
				Required: "ok",
				Min:      "ab",
				Max:      " abcd ",
				OneOf:    "alpha",
				ID:       "00000000-0000-4000-8000-000000000001",
			},
			wantField:     "max",
			wantTag:       "trimmed_max",
			wantViolation: true,
		},
		{
			name: "trimmed oneof rejects unknown value",
			input: validatorRulesInput{
				Required: "ok",
				Min:      "ab",
				Max:      "abc",
				OneOf:    "gamma",
				ID:       "00000000-0000-4000-8000-000000000001",
			},
			wantField:     "oneOf",
			wantTag:       "trimmed_oneof",
			wantViolation: true,
		},
		{
			name: "uuid_any rejects malformed uuid",
			input: validatorRulesInput{
				Required: "ok",
				Min:      "ab",
				Max:      "abc",
				OneOf:    "alpha",
				ID:       "not-a-uuid",
			},
			wantField:     "id",
			wantTag:       "uuid_any",
			wantViolation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := ValidateStruct(tt.input)
			if !tt.wantViolation {
				if len(violations) != 0 {
					t.Fatalf("expected no violations, got %#v", violations)
				}
				return
			}
			if len(violations) == 0 {
				t.Fatalf("expected violation for %s", tt.wantField)
			}
			got := violations[0]
			if got.Field != tt.wantField || got.Tag != tt.wantTag {
				t.Fatalf("violation = %#v, want field=%q tag=%q", got, tt.wantField, tt.wantTag)
			}
		})
	}
}

func TestViolationsToFieldErrorsKeepsFirstMessage(t *testing.T) {
	errs := ViolationsToFieldErrors([]FieldViolation{
		{Field: "name", Tag: "trimmed_required"},
		{Field: "name", Tag: "trimmed_max"},
		{Field: "id", Tag: "uuid_any"},
	}, DefaultValidationMessage)

	if errs["name"] != "Required." {
		t.Fatalf("name error = %q", errs["name"])
	}
	if errs["id"] != "Must be a UUID." {
		t.Fatalf("id error = %q", errs["id"])
	}
	if got := len(errs); got != 2 {
		t.Fatalf("len(errs) = %d", got)
	}
}

func TestDefaultValidationMessage(t *testing.T) {
	msg := DefaultValidationMessage(FieldViolation{Tag: "gte", Param: "1"})
	if !strings.Contains(msg, "gte 1") {
		t.Fatalf("message = %q", msg)
	}
}
