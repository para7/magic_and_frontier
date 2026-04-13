package item

import (
	"strings"
	"testing"
)

func TestBuildItemComponentsFormat(t *testing.T) {
	value, errMsg := BuildItemComponents(Item{
		ID: "item_1",
		Minecraft: MinecraftItem{
			ItemID: "minecraft:diamond_sword",
			Components: map[string]string{
				"minecraft:custom_name":       `'{"text":"Blade"}'`,
				"minecraft:lore":              `['{"text":"line1"}','{"text":"line2"}']`,
				"minecraft:enchantments":      `{"minecraft:sharpness":5}`,
				"minecraft:unbreakable":       `{}`,
				"minecraft:custom_model_data": `{floats:[42f]}`,
			},
		},
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

func TestBuildItemComponentsSortsComponentKeys(t *testing.T) {
	value, errMsg := BuildItemComponents(Item{
		ID: "item_1",
		Minecraft: MinecraftItem{
			ItemID: "minecraft:stone",
			Components: map[string]string{
				"minecraft:z": `{}`,
				"minecraft:a": `{}`,
			},
		},
	})
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	if strings.Index(value, `"minecraft:a":{}`) > strings.Index(value, `"minecraft:z":{}`) {
		t.Fatalf("component keys should be sorted for stable output: %s", value)
	}
}

func TestBuildItemComponentsRejectsInvalidComponentKey(t *testing.T) {
	_, errMsg := BuildItemComponents(Item{
		ID: "item_1",
		Minecraft: MinecraftItem{
			ItemID: "minecraft:book",
			Components: map[string]string{
				"display": `{}`,
			},
		},
	})
	if errMsg == "" {
		t.Fatal("expected invalid component key error, got none")
	}
}

func TestNormalizeComponentsTrimsAndSortsKeys(t *testing.T) {
	entries, errMsg := NormalizeComponents(map[string]string{
		" minecraft:z ": " {} ",
		"minecraft:a":   " {levels:{\"minecraft:sharpness\":5}} ",
	})
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}
	if len(entries) != 2 {
		t.Fatalf("entries length = %d, want 2", len(entries))
	}
	if entries[0].Key != "minecraft:a" || entries[0].Value != `{levels:{"minecraft:sharpness":5}}` {
		t.Fatalf("unexpected first entry: %#v", entries[0])
	}
	if entries[1].Key != "minecraft:z" || entries[1].Value != "{}" {
		t.Fatalf("unexpected second entry: %#v", entries[1])
	}
}
