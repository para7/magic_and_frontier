package export_convert

import (
	"strings"
	"testing"

	enemyModel "maf_command_editor/app/domain/model/enemy"
)

func boolPtr(v bool) *bool {
	return &v
}

func TestEnemySummonNBTIncludesIsBabyAndPassengers(t *testing.T) {
	entry := enemyModel.Enemy{
		ID:       "drop_test",
		MobType:  "minecraft:zombie",
		Name:     "配達員",
		HP:       20,
		DropMode: "replace",
		IsBaby:   boolPtr(true),
		Passengers: []enemyModel.PassengerEntity{
			{
				MobType: "minecraft:creeper",
				Name:    "ウーバーイーツ",
				Tags:    []string{"maf_vh_checked"},
				Passengers: []enemyModel.PassengerEntity{
					{
						MobType: "minecraft:chicken",
						IsBaby:  boolPtr(true),
					},
				},
			},
		},
	}

	got := enemySummonNBT("maf:generated/enemy/loot/drop_test", entry, nil)
	want := `{Health:20f,DeathLootTable:"maf:generated/enemy/loot/drop_test",CustomName:{text:"配達員"},IsBaby:1b,Tags:["maf_enemy","maf_enemy_drop_test","maf_vh_checked"],Attributes:[{Name:generic.max_health,Base:20}],HandItems:[{},{}],HandDropChances:[0.085F,0.085F],ArmorItems:[{},{},{},{}],ArmorDropChances:[0.085F,0.085F,0.085F,0.085F],Passengers:[{id:"minecraft:creeper",CustomName:"ウーバーイーツ",Tags:["maf_vh_checked"],Passengers:[{id:"minecraft:chicken",IsBaby:1b}]}]}`
	if got != want {
		t.Fatalf("unexpected summon nbt\nwant: %s\ngot:  %s", want, got)
	}
}

func TestEnemySummonNBTOmitsFalseIsBaby(t *testing.T) {
	entry := enemyModel.Enemy{
		ID:       "drop_test",
		MobType:  "minecraft:zombie",
		HP:       20,
		DropMode: "replace",
		IsBaby:   boolPtr(false),
		Passengers: []enemyModel.PassengerEntity{
			{
				MobType: "minecraft:creeper",
				IsBaby:  boolPtr(false),
			},
		},
	}

	got := enemySummonNBT("maf:generated/enemy/loot/drop_test", entry, nil)
	if strings.Contains(got, "IsBaby:1b") {
		t.Fatalf("IsBaby:1b should be omitted when false: %s", got)
	}
	if !strings.Contains(got, `Passengers:[{id:"minecraft:creeper"}]`) {
		t.Fatalf("passengers should still be rendered when present: %s", got)
	}
}
