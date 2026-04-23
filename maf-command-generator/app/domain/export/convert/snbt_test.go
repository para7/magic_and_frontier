package export_convert

import (
	"encoding/json"
	"testing"
)

func TestMapToSNBTSortsKeysAndOmitsNil(t *testing.T) {
	got := MapToSNBT(map[string]any{
		"z": nil,
		"b": json.Number("2"),
		"a": true,
		"c": "x",
	})
	want := `{a:1b,b:2,c:"x"}`
	if got != want {
		t.Fatalf("unexpected snbt\nwant: %s\ngot:  %s", want, got)
	}
}

func TestMapToSNBTNestedValues(t *testing.T) {
	got := MapToSNBT(map[string]any{
		"root": map[string]any{
			"n": json.Number("0.22"),
			"l": []any{
				json.Number("2"),
				false,
				"a\"b",
				map[string]any{"k": json.Number("1")},
				nil,
			},
		},
	})
	want := `{root:{l:[2,0b,"a\"b",{k:1}],n:0.22}}`
	if got != want {
		t.Fatalf("unexpected nested snbt\nwant: %s\ngot:  %s", want, got)
	}
}

func TestMapToSNBTQuotesUnsafeKeys(t *testing.T) {
	got := MapToSNBT(map[string]any{
		"minecraft:item_name": "ok",
	})
	want := `{"minecraft:item_name":"ok"}`
	if got != want {
		t.Fatalf("unexpected snbt key escape\nwant: %s\ngot:  %s", want, got)
	}
}
