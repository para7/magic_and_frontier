package custom_validator

import (
	"testing"

	model "maf_command_editor/app/domain/model"
)

func TestMafSlugIDValidation(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "letters and digits", value: "prominence01"},
		{name: "underscore", value: "near_poison"},
		{name: "hyphen", value: "enemy-boss"},
		{name: "empty", value: "", wantErr: true},
		{name: "space", value: "fire bolt", wantErr: true},
		{name: "uppercase", value: "FireBolt", wantErr: true},
		{name: "colon", value: "foo:bar", wantErr: true},
		{name: "slash", value: "foo/bar", wantErr: true},
		{name: "dot", value: "foo.bar", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Var(tt.value, "maf_slug_id")
			if tt.wantErr && err == nil {
				t.Fatalf("expected validation error for %q", tt.value)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no validation error for %q, got %v", tt.value, err)
			}
		})
	}
}

func TestFormatValidationErrorForSlugID(t *testing.T) {
	err := model.ValidationError{
		Entity: "grimoire",
		ID:     "Fire Bolt",
		Field:  "id",
		Tag:    "maf_slug_id",
	}

	got := FormatValidationError(err)
	want := "grimoire【Fire Bolt】id: 半角小文字英数字、_、- のみ使用できます"
	if got != want {
		t.Fatalf("FormatValidationError() = %q, want %q", got, want)
	}
}
