package export_convert

import (
	"fmt"
	"strings"
	"testing"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	passiveModel "maf_command_editor/app/domain/model/passive"
)

func TestGrimoireBookAndLootShareModel(t *testing.T) {
	entry := grimoireModel.Grimoire{
		ID:          "fire_1",
		CastTime:    40,
		CoolTime:    20,
		MPCost:      13,
		Title:       `Fire "Bolt"`,
		Description: `Deal "big" damage`,
	}

	book := GrimoireToBook(entry)
	wantItemName := fmt.Sprintf("minecraft:item_name={text:%s}", JsonString(entry.Title))
	if !strings.Contains(book, wantItemName) {
		t.Fatalf("book should contain escaped item_name; got: %s", book)
	}
	wantLore := fmt.Sprintf(
		"minecraft:lore=[%s,%s]",
		textComponentSNBT(entry.Description),
		textComponentSNBT(fmt.Sprintf("消費MP:%d 詠唱時間:%d", entry.MPCost, entry.CastTime)),
	)
	if !strings.Contains(book, wantLore) {
		t.Fatalf("book should contain updated grimoire lore; got: %s", book)
	}
	wantCustomData := spellCustomData(entry)
	if !strings.Contains(book, "minecraft:custom_data="+wantCustomData) {
		t.Fatalf("book should contain spell custom_data; got: %s", book)
	}

	lootEntry := toSpellLootEntry(entry, nil, nil)
	components := lootComponentsByFunction(t, lootEntry)
	itemNameComponent := mapByKey(t, components, "minecraft:item_name")
	if itemNameComponent["text"] != entry.Title {
		t.Fatalf("loot item_name mismatch: %#v", itemNameComponent)
	}
	lore := loreLinesByKey(t, components, "minecraft:lore")
	if len(lore) != 2 {
		t.Fatalf("loot lore line count = %d, want 2", len(lore))
	}
	if lore[0] != entry.Description {
		t.Fatalf("loot lore[0] mismatch: got %q want %q", lore[0], entry.Description)
	}
	wantLine2 := fmt.Sprintf("消費MP:%d 詠唱時間:%d", entry.MPCost, entry.CastTime)
	if lore[1] != wantLine2 {
		t.Fatalf("loot lore[1] mismatch: got %q want %q", lore[1], wantLine2)
	}
	customData := customDataTagByFunction(t, lootEntry)
	if customData != wantCustomData {
		t.Fatalf("loot custom_data mismatch: got %q want %q", customData, wantCustomData)
	}
}

func TestPassiveBookAndLootShareModel(t *testing.T) {
	entry := passiveModel.Passive{
		ID:          "passive_1",
		Name:        "Quickstep",
		Role:        "素早く動ける",
		Condition:   "always",
		Slots:       []int{1, 2},
		Description: "",
	}
	slot := 2

	book := PassiveToBook(entry, slot)
	wantItemName := passiveItemName(entry)
	if !strings.Contains(book, fmt.Sprintf("minecraft:item_name={text:%s}", JsonString(wantItemName))) {
		t.Fatalf("book should contain passive item name; got: %s", book)
	}
	wantCustomData := passiveSpellCustomData(entry, slot)
	if !strings.Contains(book, "minecraft:custom_data="+wantCustomData) {
		t.Fatalf("book should contain passive custom_data; got: %s", book)
	}

	lootEntry := toPassiveLootEntry(entry, slot, nil, nil)
	components := lootComponentsByFunction(t, lootEntry)
	itemNameComponent := mapByKey(t, components, "minecraft:item_name")
	if itemNameComponent["text"] != wantItemName {
		t.Fatalf("loot item_name mismatch: %#v", itemNameComponent)
	}
	lore := loreLinesByKey(t, components, "minecraft:lore")
	if len(lore) != 2 {
		t.Fatalf("loot lore line count = %d, want 2", len(lore))
	}
	if lore[0] != entry.Role {
		t.Fatalf("loot lore[0] mismatch: got %q want %q", lore[0], entry.Role)
	}
	wantSlotLine := fmt.Sprintf("パッシブスキル / スロット%d", slot)
	if lore[1] != wantSlotLine {
		t.Fatalf("loot lore[1] mismatch: got %q want %q", lore[1], wantSlotLine)
	}
	customData := customDataTagByFunction(t, lootEntry)
	if customData != wantCustomData {
		t.Fatalf("loot custom_data mismatch: got %q want %q", customData, wantCustomData)
	}
}

func TestGrimoireToBookEscapesSpecialCharacters(t *testing.T) {
	entry := grimoireModel.Grimoire{
		ID:          "g1",
		CastTime:    10,
		CoolTime:    0,
		MPCost:      5,
		Title:       `Quote " and \ slash`,
		Description: `desc "line" \ path`,
	}

	book := GrimoireToBook(entry)
	if !strings.Contains(book, fmt.Sprintf("id:%s", JsonString(entry.ID))) {
		t.Fatalf("id should be embedded in custom_data: %s", book)
	}
	if !strings.Contains(book, fmt.Sprintf("title:%s", JsonString(entry.Title))) {
		t.Fatalf("title should be JSON-escaped in custom_data: %s", book)
	}
	if !strings.Contains(book, fmt.Sprintf("description:%s", JsonString(entry.Description))) {
		t.Fatalf("description should be JSON-escaped in custom_data: %s", book)
	}
}

func lootComponentsByFunction(t *testing.T, lootEntry map[string]any) map[string]any {
	t.Helper()
	componentsFunction := lootFunctionByName(t, lootEntry, "minecraft:set_components")
	return mapByKey(t, componentsFunction, "components")
}

func customDataTagByFunction(t *testing.T, lootEntry map[string]any) string {
	t.Helper()
	customDataFunction := lootFunctionByName(t, lootEntry, "minecraft:set_custom_data")
	raw, ok := customDataFunction["tag"].(string)
	if !ok {
		t.Fatalf("custom_data tag must be string: %#v", customDataFunction["tag"])
	}
	return raw
}

func lootFunctionByName(t *testing.T, lootEntry map[string]any, functionID string) map[string]any {
	t.Helper()
	rawFunctions, ok := lootEntry["functions"].([]any)
	if !ok {
		t.Fatalf("loot entry functions missing: %#v", lootEntry["functions"])
	}
	for _, rawFn := range rawFunctions {
		fn, ok := rawFn.(map[string]any)
		if !ok {
			continue
		}
		if fn["function"] == functionID {
			return fn
		}
	}
	t.Fatalf("function %q not found in %#v", functionID, rawFunctions)
	return nil
}

func mapByKey(t *testing.T, input map[string]any, key string) map[string]any {
	t.Helper()
	raw, ok := input[key].(map[string]any)
	if !ok {
		t.Fatalf("%q should be map[string]any: %#v", key, input[key])
	}
	return raw
}

func loreLinesByKey(t *testing.T, input map[string]any, key string) []string {
	t.Helper()
	raw, ok := input[key].([]any)
	if !ok {
		t.Fatalf("%q should be []any: %#v", key, input[key])
	}
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		lineMap, ok := line.(map[string]any)
		if !ok {
			t.Fatalf("lore line should be map[string]any: %#v", line)
		}
		text, ok := lineMap["text"].(string)
		if !ok {
			t.Fatalf("lore line text should be string: %#v", lineMap["text"])
		}
		out = append(out, text)
	}
	return out
}
