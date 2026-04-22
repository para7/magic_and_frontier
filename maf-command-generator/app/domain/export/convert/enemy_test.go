package export_convert

import (
	"encoding/json"
	"testing"

	model "maf_command_editor/app/domain/model"
	enemyModel "maf_command_editor/app/domain/model/enemy"
	itemModel "maf_command_editor/app/domain/model/item"
)

func floatPtr(v float64) *float64 {
	return &v
}

func TestEnemySummonNBTMergesMinecraftAndMaf(t *testing.T) {
	entry := enemyModel.Enemy{
		ID:      "drop_test",
		MobType: "minecraft:zombie",
		Minecraft: map[string]any{
			"Health":     json.Number("20"),
			"CustomName": `{"text":"配達員"}`,
			"Tags":       []any{"base_tag"},
			"equipment": map[string]any{
				"head": map[string]any{"id": "minecraft:carved_pumpkin", "count": json.Number("1")},
			},
			"drop_chances": map[string]any{
				"head": json.Number("0.2"),
			},
		},
		Maf: enemyModel.EnemyMaf{
			Equipment: model.Equipment{
				Mainhand: &model.EquipmentSlot{Kind: "item", RefID: "items_1", Count: 1},
				Head:     &model.EquipmentSlot{Kind: "minecraft_item", RefID: "minecraft:diamond_helmet", Count: 1, DropChance: floatPtr(0.5)},
			},
			EnemySkillIDs: []string{"near_poison"},
			DropMode:      "replace",
		},
		Passengers: []enemyModel.Passenger{
			{
				MobType: "minecraft:creeper",
				Minecraft: map[string]any{
					"CustomName": `"ウーバーイーツ"`,
				},
				Maf: enemyModel.PassengerMaf{
					Tags:          []string{"maf_vh_checked"},
					EnemySkillIDs: []string{"creeper_instant_explode"},
				},
			},
		},
	}

	itemsByID := map[string]itemModel.Item{
		"items_1": {ItemID: "minecraft:stone"},
	}

	got := enemySummonNBT("maf:generated/enemy/loot/drop_test", entry, itemsByID)
	want := `{CustomName:"{\"text\":\"配達員\"}",DeathLootTable:"maf:generated/enemy/loot/drop_test",Health:20,Passengers:[{CustomName:"\"ウーバーイーツ\"",Tags:["maf_vh_checked","EnemySkill","creeper_instant_explode","maf_enemy_skill_creeper_instant_explode"],id:"minecraft:creeper"}],Tags:["base_tag","maf_enemy","maf_enemy_drop_test","maf_vh_checked","EnemySkill","near_poison","maf_enemy_skill_near_poison"],drop_chances:{head:0.5,mainhand:0.085},equipment:{head:{count:1,id:"minecraft:diamond_helmet"},mainhand:{count:1,id:"minecraft:stone"}}}`
	if got != want {
		t.Fatalf("unexpected summon nbt\nwant: %s\ngot:  %s", want, got)
	}
}

func TestEnemySummonNBTMergesExistingPassengers(t *testing.T) {
	entry := enemyModel.Enemy{
		ID:      "drop_test",
		MobType: "minecraft:zombie",
		Minecraft: map[string]any{
			"Passengers": []any{
				map[string]any{"id": "minecraft:chicken"},
			},
		},
		Maf: enemyModel.EnemyMaf{
			DropMode: "replace",
		},
		Passengers: []enemyModel.Passenger{
			{MobType: "minecraft:creeper"},
		},
	}

	got := enemySummonNBT("maf:generated/enemy/loot/drop_test", entry, nil)
	want := `{DeathLootTable:"maf:generated/enemy/loot/drop_test",Passengers:[{id:"minecraft:chicken"},{id:"minecraft:creeper"}],Tags:["maf_enemy","maf_enemy_drop_test","maf_vh_checked"]}`
	if got != want {
		t.Fatalf("unexpected merged passengers\nwant: %s\ngot:  %s", want, got)
	}
}

func TestEnemySummonNBTMafMainhandOverridesMinecraftMainhand(t *testing.T) {
	entry := enemyModel.Enemy{
		ID:      "drop_test",
		MobType: "minecraft:zombie",
		Minecraft: map[string]any{
			"equipment": map[string]any{
				"mainhand": map[string]any{
					"id":    "minecraft:wooden_sword",
					"count": json.Number("1"),
				},
			},
			"drop_chances": map[string]any{
				"mainhand": json.Number("0.2"),
			},
		},
		Maf: enemyModel.EnemyMaf{
			Equipment: model.Equipment{
				Mainhand: &model.EquipmentSlot{
					Kind:       "minecraft_item",
					RefID:      "minecraft:diamond_sword",
					Count:      2,
					DropChance: floatPtr(0.9),
				},
			},
			DropMode: "replace",
		},
	}

	got := enemySummonNBT("maf:generated/enemy/loot/drop_test", entry, nil)
	want := `{DeathLootTable:"maf:generated/enemy/loot/drop_test",Tags:["maf_enemy","maf_enemy_drop_test","maf_vh_checked"],drop_chances:{mainhand:0.9},equipment:{mainhand:{count:2,id:"minecraft:diamond_sword"}}}`
	if got != want {
		t.Fatalf("unexpected mainhand precedence\nwant: %s\ngot:  %s", want, got)
	}
}
