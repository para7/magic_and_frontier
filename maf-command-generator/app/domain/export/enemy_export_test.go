package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	model "maf_command_editor/app/domain/model"
	enemyModel "maf_command_editor/app/domain/model/enemy"
	enemyskillModel "maf_command_editor/app/domain/model/enemyskill"
	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	itemModel "maf_command_editor/app/domain/model/item"
	mc "maf_command_editor/app/minecraft"
)

func TestBuildEnemySkillArtifactsAndWrite(t *testing.T) {
	master := exportMasterStub{
		enemySkills: []enemyskillModel.EnemySkill{
			{ID: "near_poison", Script: "effect give @e[distance=..2] minecraft:poison 10 2"},
		},
	}

	artifacts := BuildEnemySkillArtifacts(master, "generated/enemy/skill")
	if len(artifacts) != 1 {
		t.Fatalf("artifacts length = %d, want 1", len(artifacts))
	}
	if artifacts[0].DispatcherLine != "execute if entity @s[tag=near_poison] run function maf:generated/enemy/skill/near_poison" {
		t.Fatalf("unexpected dispatcher line: %q", artifacts[0].DispatcherLine)
	}

	dir := t.TempDir()
	if err := WriteEnemySkillArtifacts(dir, artifacts); err != nil {
		t.Fatal(err)
	}

	body, err := os.ReadFile(filepath.Join(dir, "near_poison.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "effect give @e[distance=..2] minecraft:poison 10 2\n" {
		t.Fatalf("unexpected skill body: %q", string(body))
	}

	mainBody, err := os.ReadFile(filepath.Join(dir, "main.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(mainBody)
	if !strings.Contains(text, "enemy skill dispatcher") {
		t.Fatalf("main.mcfunction should include header comment: %s", text)
	}
	if !strings.Contains(text, artifacts[0].DispatcherLine) {
		t.Fatalf("main.mcfunction should include dispatcher line: %s", text)
	}
}

func TestBuildEnemyArtifactsReplace(t *testing.T) {
	attack := 6.0
	defense := 2.0
	moveSpeed := 0.22
	master := exportMasterStub{
		items: []itemModel.Item{
			{ID: "items_1", ItemID: "minecraft:stone", NBT: "{id:\"minecraft:stone\",Count:1b}"},
		},
		grimoires: []grimoireModel.Grimoire{
			{ID: "grimoire_1", CastID: 1, CastTime: 10, MPCost: 5, Title: "Firebolt", Description: "Basic sample projectile spell."},
		},
		enemies: []enemyModel.Enemy{
			{
				ID:        "enemy_1",
				MobType:   "minecraft:zombie",
				Name:      "Sample Zombie",
				HP:        40,
				Attack:    &attack,
				Defense:   &defense,
				MoveSpeed: &moveSpeed,
				EnemySkillIDs: []string{
					"enemyskill_1",
				},
				DropMode: "replace",
				Equipment: model.Equipment{
					Mainhand: &model.EquipmentSlot{
						Kind:  "item",
						RefID: "items_1",
						Count: 1,
					},
				},
				Drops: []model.DropRef{
					{Kind: "item", RefID: "items_1", Weight: 70, CountMin: ptrFloat(1), CountMax: ptrFloat(3)},
					{Kind: "grimoire", RefID: "grimoire_1", Weight: 30, CountMin: ptrFloat(1), CountMax: ptrFloat(1)},
				},
			},
		},
	}

	artifacts, err := BuildEnemyArtifacts(master, "generated/enemy/loot", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 1 {
		t.Fatalf("artifacts length = %d, want 1", len(artifacts))
	}
	text := artifacts[0].SummonScript
	if !strings.Contains(text, "summon minecraft:zombie") {
		t.Fatalf("summon script should contain mob type: %s", text)
	}
	if !strings.Contains(text, `DeathLootTable:"maf:generated/enemy/loot/enemy_1"`) {
		t.Fatalf("summon script should contain loot table ref: %s", text)
	}
	if !strings.Contains(text, `CustomName:{text:"Sample Zombie"}`) {
		t.Fatalf("summon script should contain component-style custom name: %s", text)
	}
	if !strings.Contains(text, `"maf_enemy_skill_enemyskill_1"`) {
		t.Fatalf("summon script should contain enemy skill tag: %s", text)
	}
	if !strings.Contains(text, `id:"minecraft:stone"`) {
		t.Fatalf("summon script should resolve equipment/drop item id: %s", text)
	}

	pools, ok := artifacts[0].LootTable["pools"].([]any)
	if !ok || len(pools) != 1 {
		t.Fatalf("loot table pools = %#v, want 1 pool", artifacts[0].LootTable["pools"])
	}
}

func TestBuildEnemyArtifactsAppendMergesMinecraftLoot(t *testing.T) {
	root := t.TempDir()
	basePath, err := mc.FilePathForTable(root, "minecraft:entities/zombie")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(basePath), 0o755); err != nil {
		t.Fatal(err)
	}
	baseLoot := map[string]any{
		"type": "minecraft:entity",
		"pools": []any{
			map[string]any{
				"rolls": 1,
				"entries": []any{
					map[string]any{"type": "minecraft:item", "name": "minecraft:rotten_flesh"},
				},
			},
		},
	}
	baseData, err := json.Marshal(baseLoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(basePath, baseData, 0o644); err != nil {
		t.Fatal(err)
	}

	master := exportMasterStub{
		enemies: []enemyModel.Enemy{
			{
				ID:       "enemy_2",
				MobType:  "minecraft:zombie",
				HP:       20,
				DropMode: "append",
				Drops: []model.DropRef{
					{Kind: "minecraft_item", RefID: "minecraft:apple", Weight: 1, CountMin: ptrFloat(1), CountMax: ptrFloat(1)},
				},
			},
		},
	}

	artifacts, err := BuildEnemyArtifacts(master, "generated/enemy/loot", root)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 1 {
		t.Fatalf("artifacts length = %d, want 1", len(artifacts))
	}
	pools, ok := artifacts[0].LootTable["pools"].([]any)
	if !ok {
		t.Fatalf("loot table pools = %#v", artifacts[0].LootTable["pools"])
	}
	if len(pools) != 2 {
		t.Fatalf("append mode should merge pools: got %d", len(pools))
	}
}

func TestWriteEnemyArtifactsWritesFiles(t *testing.T) {
	enemyDir := filepath.Join(t.TempDir(), "function")
	lootDir := filepath.Join(t.TempDir(), "loot")
	artifacts := []EnemyArtifact{
		{
			ID:           "enemy_1",
			SummonScript: "summon minecraft:zombie ~ ~ ~ {}",
			LootTable: map[string]any{
				"type": "minecraft:generic",
				"pools": []any{
					map[string]any{"rolls": 1, "entries": []any{}},
				},
			},
		},
	}

	if err := WriteEnemyArtifacts(enemyDir, lootDir, artifacts); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(enemyDir, "enemy_1.mcfunction")); err != nil {
		t.Fatalf("enemy function missing: %v", err)
	}
	lootPath := filepath.Join(lootDir, "enemy_1.json")
	data, err := os.ReadFile(lootPath)
	if err != nil {
		t.Fatalf("enemy loot missing: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("enemy loot should be valid json: %v", err)
	}
}
