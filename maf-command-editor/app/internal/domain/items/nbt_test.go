package items

import (
	"strings"
	"testing"
)

func TestSplitSNBTEntries(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantKeys []string
	}{
		{name: "empty", input: "", wantKeys: nil},
		{name: "whitespace only", input: "   ", wantKeys: nil},
		{name: "single entry", input: "Foo:bar", wantKeys: []string{"Foo"}},
		{name: "multiple flat entries", input: "A:1,B:2,C:3", wantKeys: []string{"A", "B", "C"}},
		{name: "nested braces", input: "A:{x:1,y:2},B:3", wantKeys: []string{"A", "B"}},
		{name: "nested brackets", input: "A:[1,2,3],B:test", wantKeys: []string{"A", "B"}},
		{name: "double quoted string with comma", input: `A:"hello,world",B:2`, wantKeys: []string{"A", "B"}},
		{name: "single quoted string with comma", input: `Name:'{"text":"hi,there"}',X:1`, wantKeys: []string{"Name", "X"}},
		{name: "deep nesting", input: "A:{b:{c:[1,{d:2}]}},B:ok", wantKeys: []string{"A", "B"}},
		{name: "escaped quote inside double quote", input: `A:"say \"hi\"",B:1`, wantKeys: []string{"A", "B"}},
		{name: "whitespace around keys", input: " A : 1 , B : 2 ", wantKeys: []string{"A", "B"}},
		{name: "outer braces already stripped", input: `display:{Name:'{"text":"test"}'}`, wantKeys: []string{"display"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitSNBTEntries(tt.input)
			if len(got) != len(tt.wantKeys) {
				t.Fatalf("len=%d want %d; entries=%v", len(got), len(tt.wantKeys), got)
			}
			for i, key := range tt.wantKeys {
				if got[i].Key != key {
					t.Errorf("entry[%d].Key = %q, want %q", i, got[i].Key, key)
				}
			}
		})
	}
}

func TestBuildItemNBTCustomMerge(t *testing.T) {
	tests := []struct {
		name      string
		input     SaveInput
		wantIn    []string
		wantNotIn []string
	}{
		{
			name: "custom key with no conflict passes through",
			input: SaveInput{
				ID:        "items_1",
				ItemID:    "minecraft:stone",
				CustomNBT: `{SomeCustomKey:42}`,
			},
			wantIn: []string{"SomeCustomKey:42"},
		},
		{
			name: "form enchantments wins over custom Enchantments",
			input: SaveInput{
				ID:           "items_1",
				ItemID:       "minecraft:sword",
				Enchantments: "minecraft:sharpness 5",
				CustomNBT:    `{Enchantments:[{id:"minecraft:fire_aspect",lvl:2s}]}`,
			},
			wantIn:    []string{`"minecraft:sharpness"`},
			wantNotIn: []string{`"minecraft:fire_aspect"`},
		},
		{
			name: "form display wins over custom display",
			input: SaveInput{
				ID:         "items_1",
				ItemID:     "minecraft:stone",
				CustomName: "FormName",
				CustomNBT:  `{display:{Name:'{"text":"CustomName"}'}}`,
			},
			wantIn:    []string{"FormName"},
			wantNotIn: []string{"CustomName"},
		},
		{
			name: "custom Enchantments passes when form enchantments empty",
			input: SaveInput{
				ID:        "items_1",
				ItemID:    "minecraft:stone",
				CustomNBT: `{Enchantments:[{id:"minecraft:mending",lvl:1s}]}`,
			},
			wantIn: []string{`"minecraft:mending"`},
		},
		{
			name: "partial conflict: conflicting key dropped, other key kept",
			input: SaveInput{
				ID:              "items_1",
				ItemID:          "minecraft:stone",
				CustomModelData: "42",
				CustomNBT:       `{CustomModelData:99,ExtraKey:123}`,
			},
			wantIn:    []string{"CustomModelData:42", "ExtraKey:123"},
			wantNotIn: []string{"CustomModelData:99"},
		},
		{
			name: "no custom NBT unchanged",
			input: SaveInput{
				ID:         "items_1",
				ItemID:     "minecraft:stone",
				CustomName: "TestItem",
			},
			wantIn: []string{"Count:1b", "TestItem"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nbt, errMsg := buildItemNBT(tt.input)
			if errMsg != "" {
				t.Fatalf("unexpected error: %s", errMsg)
			}
			for _, want := range tt.wantIn {
				if !strings.Contains(nbt, want) {
					t.Errorf("expected %q in nbt, got: %s", want, nbt)
				}
			}
			for _, notWant := range tt.wantNotIn {
				if strings.Contains(nbt, notWant) {
					t.Errorf("expected %q NOT in nbt, got: %s", notWant, nbt)
				}
			}
		})
	}
}
