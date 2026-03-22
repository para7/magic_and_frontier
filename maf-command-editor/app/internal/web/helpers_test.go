package web

import "testing"

func TestParseTreasurePools(t *testing.T) {
	t.Run("parses valid lines", func(t *testing.T) {
		errs := map[string]string{}
		pools := parseTreasurePools(errs, "item,item_1,1,1,2\ngrimoire,grimoire_1,0.5,,")
		if len(errs) != 0 {
			t.Fatalf("errs = %#v", errs)
		}
		if len(pools) != 2 {
			t.Fatalf("len(pools) = %d", len(pools))
		}
		if pools[0].Kind != "item" || pools[0].RefID != "item_1" || pools[0].Weight != 1 {
			t.Fatalf("pools[0] = %#v", pools[0])
		}
		if pools[0].CountMin == nil || *pools[0].CountMin != 1 {
			t.Fatalf("pools[0].CountMin = %#v", pools[0].CountMin)
		}
		if pools[0].CountMax == nil || *pools[0].CountMax != 2 {
			t.Fatalf("pools[0].CountMax = %#v", pools[0].CountMax)
		}
		if pools[1].CountMin != nil || pools[1].CountMax != nil {
			t.Fatalf("pools[1] = %#v", pools[1])
		}
	})

	t.Run("rejects non numeric count", func(t *testing.T) {
		errs := map[string]string{}
		pools := parseTreasurePools(errs, "item,item_1,1,nope,2")
		if pools != nil {
			t.Fatalf("pools = %#v", pools)
		}
		if errs["lootPoolsText"] != "Count values must be numeric when provided." {
			t.Fatalf("errs = %#v", errs)
		}
	})
}

func TestParseEquipment(t *testing.T) {
	t.Run("parses valid slots", func(t *testing.T) {
		errs := map[string]string{}
		equipment := parseEquipment(errs, "mainhand,minecraft_item,minecraft:iron_sword,1,0.25\nhead,item,item_1,1,")
		if len(errs) != 0 {
			t.Fatalf("errs = %#v", errs)
		}
		if equipment.Mainhand == nil || equipment.Mainhand.RefID != "minecraft:iron_sword" {
			t.Fatalf("equipment = %#v", equipment)
		}
		if equipment.Mainhand.DropChance == nil || *equipment.Mainhand.DropChance != 0.25 {
			t.Fatalf("equipment.Mainhand = %#v", equipment.Mainhand)
		}
		if equipment.Head == nil || equipment.Head.Kind != "item" || equipment.Head.RefID != "item_1" {
			t.Fatalf("equipment.Head = %#v", equipment.Head)
		}
	})

	t.Run("rejects invalid slot", func(t *testing.T) {
		errs := map[string]string{}
		equipment := parseEquipment(errs, "wing,item,item_1,1,")
		if equipment.Mainhand != nil || equipment.Head != nil {
			t.Fatalf("equipment = %#v", equipment)
		}
		if errs["equipmentText"] != "Equipment slot must be one of mainhand,offhand,head,chest,legs,feet." {
			t.Fatalf("errs = %#v", errs)
		}
	})
}
