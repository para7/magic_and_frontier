package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"tools2/app/internal/domain/enemies"
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

func TestGenerateItemOutputsUsesConfiguredLootDir(t *testing.T) {
	settings := ExportSettings{
		OutputRoot: t.TempDir(),
		Namespace:  "maf",
		Paths: ExportPaths{
			ItemFunctionDir: "data/maf/function/generated/item",
			ItemLootDir:     "data/maf/loot_table/generated/item",
		},
	}

	_, err := generateItemOutputs(settings, []items.ItemEntry{{
		ID:     "items_1",
		ItemID: "minecraft:apple",
		Count:  2,
	}})
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(settings.OutputRoot, settings.Paths.ItemFunctionDir, "items_1.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "loot give @s loot maf:generated/item/items_1") {
		t.Fatalf("item function should reference generated loot table path: %s", string(data))
	}
}

func TestToEnemyFunctionLinesUsesMobTypeAndDeathLootTable(t *testing.T) {
	attack := 6.0
	entry := enemies.EnemyEntry{
		ID:        "enemy_3",
		MobType:   "minecraft:skeleton",
		Name:      "Guardian",
		HP:        40,
		Attack:    &attack,
		DropMode:  "replace",
		Drops:     []enemies.DropRef{{Kind: "minecraft_item", RefID: "minecraft:bone", Weight: 1}},
		Equipment: enemies.Equipment{},
		EnemySkillIDs: []string{
			"enemyskill_1",
		},
	}
	settings := ExportSettings{
		Namespace: "maf",
		Paths: ExportPaths{
			EnemyLootDir: "data/maf/loot_table/generated/enemy",
		},
	}

	lines := toEnemyFunctionLines(settings, entry, map[string]items.ItemEntry{})
	text := strings.Join(lines, "\n")
	if !strings.Contains(text, "summon minecraft:skeleton") {
		t.Fatalf("enemy summon should use mob type: %s", text)
	}
	if !strings.Contains(text, `DeathLootTable:"maf:generated/enemy/enemy_3"`) {
		t.Fatalf("enemy summon should reference enemy loot table: %s", text)
	}
	if !strings.Contains(text, `"maf_enemy_skill_enemyskill_1"`) {
		t.Fatalf("enemy summon should encode skill tag: %s", text)
	}
}

func TestToEnemyFunctionLinesResolvesCustomItemEquipment(t *testing.T) {
	entry := enemies.EnemyEntry{
		ID:       "enemy_1",
		MobType:  "minecraft:zombie",
		HP:       20,
		DropMode: "replace",
		Equipment: enemies.Equipment{
			Mainhand: &enemies.EquipmentSlot{
				Kind:  "item",
				RefID: "items_3",
				Count: 1,
			},
		},
	}
	settings := ExportSettings{
		Namespace: "maf",
		Paths: ExportPaths{
			EnemyLootDir: "data/maf/loot_table/generated/enemy",
		},
	}

	lines := toEnemyFunctionLines(settings, entry, map[string]items.ItemEntry{
		"items_3": {ID: "items_3", ItemID: "minecraft:iron_sword"},
	})
	text := strings.Join(lines, "\n")
	if !strings.Contains(text, `id:"minecraft:iron_sword"`) {
		t.Fatalf("enemy summon should resolve custom item refs to minecraft item ids: %s", text)
	}
	if strings.Contains(text, `id:"items_3"`) {
		t.Fatalf("enemy summon should not emit internal item ids: %s", text)
	}
}

func TestLootTableOutputPathRejectsTraversal(t *testing.T) {
	settings := ExportSettings{OutputRoot: "/tmp/out"}
	if _, err := lootTableOutputPath(settings, "maf:loot/../escape"); err == nil {
		t.Fatalf("expected traversal table path to be rejected")
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

func TestGenerateGrimoireDebugFunctionsCreatesPerEntryFile(t *testing.T) {
	settings := ExportSettings{
		OutputRoot: t.TempDir(),
		Namespace:  "maf",
	}
	entries := []grimoire.GrimoireEntry{
		{
			ID:          "grimoire_1",
			CastID:      1,
			CastTime:    20,
			MPCost:      5,
			Title:       "Firebolt",
			Description: "Basic sample projectile spell.",
		},
	}

	count, err := generateGrimoireDebugFunctions(settings, entries)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("generated count = %d, want 1", count)
	}

	path := filepath.Join(settings.OutputRoot, "data", "maf", "function", "generated", "debug", "grimoire", "grimoire_1.mcfunction")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.TrimSpace(string(data))
	if !strings.HasPrefix(text, "give @s minecraft:written_book[") {
		t.Fatalf("debug file should use direct give command: %s", text)
	}
	if !strings.Contains(text, `custom_data={maf:{grimoire_id:"grimoire_1",spell:{castid:1,cost:5,cast:20,title:"Firebolt",description:"Basic sample projectile spell."}}}`) {
		t.Fatalf("debug file should include spell custom_data: %s", text)
	}
}
