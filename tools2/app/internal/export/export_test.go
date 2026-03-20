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
