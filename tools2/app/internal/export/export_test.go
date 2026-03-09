package export

import (
	"encoding/json"
	"strings"
	"testing"

	"tools2/app/internal/domain/grimoire"
)

func TestToSpellLootTable_OmitsLoreFunction(t *testing.T) {
	entry := grimoire.GrimoireEntry{
		ID:          "33333333-3333-4333-8333-333333333333",
		CastID:      100,
		Script:      "maf:spell/firebolt",
		Title:       "Firebolt",
		Description: "Basic sample projectile spell.\nSecond line.",
	}
	variant := grimoire.Variant{Cast: 10, Cost: 5}

	lootTable := toSpellLootTable(entry, variant)
	data, err := json.Marshal(lootTable)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, "set_lore") {
		t.Fatalf("loot table should contain set_lore: %s", text)
	}
	if !strings.Contains(text, `"text":"Basic sample projectile spell."`) || !strings.Contains(text, `"text":"Second line."`) {
		t.Fatalf("loot table should contain description lines: %s", text)
	}
}

func TestToSpellGiveCommand_OmitsLoreComponent(t *testing.T) {
	entry := grimoire.GrimoireEntry{
		ID:          "33333333-3333-4333-8333-333333333333",
		CastID:      100,
		Script:      "maf:spell/firebolt",
		Title:       "Firebolt",
		Description: "Basic sample projectile spell.\nSecond line.",
	}
	variant := grimoire.Variant{Cast: 10, Cost: 5}

	command := toSpellGiveCommand(entry, variant)
	if !strings.Contains(command, "lore:[") {
		t.Fatalf("give command should contain lore component: %s", command)
	}
	if !strings.Contains(command, `{"text":"Basic sample projectile spell."}`) || !strings.Contains(command, `{"text":"Second line."}`) {
		t.Fatalf("give command should contain description lines: %s", command)
	}
}
