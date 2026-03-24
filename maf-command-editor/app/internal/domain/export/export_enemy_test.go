package export

import (
	"strings"
	"testing"

	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/items"
)

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
	if !strings.Contains(text, `"EnemySkill"`) || !strings.Contains(text, `"enemyskill_1"`) {
		t.Fatalf("enemy summon should encode EnemySkill tags: %s", text)
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
