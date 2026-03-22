package export

import (
	"encoding/json"
	"strings"
	"testing"

	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/treasures"
)

func TestToSpellLootTableIncludesLoreAndNewCustomData(t *testing.T) {
	entry := grimoire.GrimoireEntry{
		ID:          "grimoire_12",
		CastID:      100,
		CastTime:    10,
		MPCost:      5,
		Script:      "function maf:grimoire/grimoire_12",
		Title:       "Firebolt",
		Description: "Basic sample projectile spell.\nSecond line.",
	}

	lootTable := toSpellLootTable(entry)
	data, err := json.Marshal(lootTable)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, "set_lore") {
		t.Fatalf("loot table should contain set_lore: %s", text)
	}
	if !strings.Contains(text, `"mode":"append"`) {
		t.Fatalf("loot table should contain lore mode: %s", text)
	}
	if !strings.Contains(text, `"target":"item_name"`) {
		t.Fatalf("loot table should contain item_name target: %s", text)
	}
	if !strings.Contains(text, `grimoire_id:\"grimoire_12\"`) {
		t.Fatalf("loot table should contain grimoire id custom data: %s", text)
	}
	if !strings.Contains(text, `castid:100`) || !strings.Contains(text, `cost:5`) || !strings.Contains(text, `cast:10`) {
		t.Fatalf("loot table should contain spell metadata: %s", text)
	}
}

func TestItemCustomDataIncludesOptionalSkillTags(t *testing.T) {
	entry := items.ItemEntry{
		ID:     "items_4",
		ItemID: "minecraft:stone",
		NBT:    "{id:\"minecraft:stone\",Count:1b}",
	}
	if strings.Contains(itemCustomData(entry), "maf_skill") {
		t.Fatalf("unexpected skill tags for item without skill")
	}

	entry.SkillID = "skill_7"
	customData := itemCustomData(entry)
	if !strings.Contains(customData, "maf_skill:1b") || !strings.Contains(customData, `"skill_7"`) {
		t.Fatalf("custom data should include skill tags: %s", customData)
	}
}

func TestBuildDropLootTableSupportsMinecraftItemAndGrimoire(t *testing.T) {
	drops := []treasures.DropRef{
		{Kind: "minecraft_item", RefID: "minecraft:apple", Weight: 1},
		{Kind: "grimoire", RefID: "grimoire_2", Weight: 2},
	}
	grimoires := map[string]grimoire.GrimoireEntry{
		"grimoire_2": {
			ID:       "grimoire_2",
			CastID:   3,
			CastTime: 8,
			MPCost:   11,
			Title:    "Heal",
		},
	}

	lootTable, err := buildDropLootTable(drops, map[string]items.ItemEntry{}, grimoires, "treasure(treasure_1)")
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(lootTable)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, `"name":"minecraft:apple"`) || !strings.Contains(text, `"name":"minecraft:written_book"`) {
		t.Fatalf("unexpected loot table: %s", text)
	}
}

func TestGrimoireDebugGiveCommandIncludesCustomDataAndLore(t *testing.T) {
	entry := grimoire.GrimoireEntry{
		ID:          "grimoire_12",
		CastID:      100,
		CastTime:    10,
		MPCost:      5,
		Title:       `Mage's "Fire"`,
		Description: "Line 1\nLine 2",
	}

	command := grimoireDebugGiveCommand(entry)
	if !strings.HasPrefix(command, "give @s minecraft:written_book[") || !strings.HasSuffix(command, "] 1") {
		t.Fatalf("unexpected command shape: %s", command)
	}
	if !strings.Contains(command, "item_name='{\"text\":\"Mage\\'s \\\"Fire\\\"\"}'") {
		t.Fatalf("item_name should be escaped safely: %s", command)
	}
	if !strings.Contains(command, "lore=['{\"text\":\"Line 1\"}','{\"text\":\"Line 2\"}']") {
		t.Fatalf("lore should include split lines: %s", command)
	}
	if !strings.Contains(command, `custom_data={maf:{grimoire_id:"grimoire_12",spell:{castid:100,cost:5,cast:10,title:"Mage's \"Fire\"",description:"Line 1\nLine 2"}}}`) {
		t.Fatalf("custom_data should include spell payload: %s", command)
	}
}

func TestGrimoireDebugGiveCommandOmitsLoreWhenDescriptionEmpty(t *testing.T) {
	entry := grimoire.GrimoireEntry{
		ID:          "grimoire_1",
		CastID:      1,
		CastTime:    20,
		MPCost:      5,
		Title:       "Firebolt",
		Description: "",
	}

	command := grimoireDebugGiveCommand(entry)
	if strings.Contains(command, "lore=[") {
		t.Fatalf("lore should be omitted when description is empty: %s", command)
	}
}
