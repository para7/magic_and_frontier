package item

import (
	"strings"
	"testing"
)

func TestBuildItemComponentsFormat(t *testing.T) {
	value, errMsg := BuildItemComponents(Item{
		ID:              "item_1",
		ItemID:          "minecraft:diamond_sword",
		CustomName:      "Blade",
		Lore:            "line1\nline2",
		Enchantments:    "minecraft:sharpness 5",
		Unbreakable:     true,
		CustomModelData: "42",
	})
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	wantIn := []string{
		`count:1`,
		`components:{`,
		`"minecraft:custom_name":`,
		`"minecraft:lore":[`,
		`"minecraft:enchantments":{"minecraft:sharpness":5}`,
		`"minecraft:unbreakable":{}`,
		`"minecraft:custom_model_data":{floats:[42f]}`,
	}
	for _, want := range wantIn {
		if !strings.Contains(value, want) {
			t.Fatalf("expected %q in result, got: %s", want, value)
		}
	}

	wantNotIn := []string{`Count:1b`, `tag:{`, `display:{`, `Enchantments:[`}
	for _, notWant := range wantNotIn {
		if strings.Contains(value, notWant) {
			t.Fatalf("expected %q not to appear, got: %s", notWant, value)
		}
	}
}

func TestBuildItemComponentsCustomMerge(t *testing.T) {
	value, errMsg := BuildItemComponents(Item{
		ID:              "item_1",
		ItemID:          "minecraft:stone",
		CustomModelData: "42",
		CustomNBT:       `{CustomModelData:99,"minecraft:custom_data":{x:1}}`,
	})
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}
	if strings.Contains(value, "CustomModelData:99") {
		t.Fatalf("legacy conflicting key should be dropped: %s", value)
	}
	if !strings.Contains(value, `"minecraft:custom_model_data":{floats:[42f]}`) {
		t.Fatalf("form value should win for custom model data: %s", value)
	}
	if !strings.Contains(value, `"minecraft:custom_data":{x:1}`) {
		t.Fatalf("non-conflicting custom key should be merged: %s", value)
	}
}

func TestSplitSNBTEntriesWithNamespacedKeys(t *testing.T) {
	entries := splitSNBTEntries(`"minecraft:custom_data":{a:1},"minecraft:enchantments":{"minecraft:sharpness":5},Legacy:1`)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d: %#v", len(entries), entries)
	}
	if entries[0].Key != `"minecraft:custom_data"` {
		t.Fatalf("unexpected key[0]: %q", entries[0].Key)
	}
	if entries[1].Key != `"minecraft:enchantments"` {
		t.Fatalf("unexpected key[1]: %q", entries[1].Key)
	}
	if entries[2].Key != `Legacy` {
		t.Fatalf("unexpected key[2]: %q", entries[2].Key)
	}
}

func TestBuildItemComponentsHideFlagsMapping(t *testing.T) {
	value, errMsg := BuildItemComponents(Item{
		ID:        "item_1",
		ItemID:    "minecraft:book",
		HideFlags: "3",
	})
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}
	if !strings.Contains(value, `"minecraft:tooltip_display":{hidden_components:[`) {
		t.Fatalf("tooltip_display should be emitted for hideFlags: %s", value)
	}
	if !strings.Contains(value, `"minecraft:enchantments"`) ||
		!strings.Contains(value, `"minecraft:stored_enchantments"`) ||
		!strings.Contains(value, `"minecraft:attribute_modifiers"`) {
		t.Fatalf("hideFlags=3 should map to enchantments and attribute modifiers: %s", value)
	}
}
