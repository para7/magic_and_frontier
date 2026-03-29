package api

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"maf-command-editor/app/internal/domain/entity/enemies"
	"maf-command-editor/app/internal/domain/entity/enemyskills"
	"maf-command-editor/app/internal/domain/entity/grimoire"
	"maf-command-editor/app/internal/domain/entity/items"
	"maf-command-editor/app/internal/domain/entity/loottables"
	"maf-command-editor/app/internal/domain/entity/skills"
	"maf-command-editor/app/internal/domain/entity/treasures"
)

func TestHandlerAPIHappyPathAndSave(t *testing.T) {
	handler, root := newTestHandler(t)

	skillID := createJSONEntry(t, handler, "/api/skills", skills.SaveInput{
		ID:          "skill_1",
		Name:        "Slash",
		Description: "Basic slash",
		Script:      "say slash",
	})

	itemID := createJSONEntry(t, handler, "/api/items", items.SaveInput{
		ID:      "items_1",
		ItemID:  "minecraft:apple",
		SkillID: skillID,
	})

	grimoireID := createJSONEntry(t, handler, "/api/grimoire", grimoire.SaveInput{
		ID:          "grimoire_1",
		Title:       "Firebolt",
		Description: "Burn",
		Script:      "say fire",
		CastTime:    20,
		MPCost:      5,
	})

	enemySkillID := createJSONEntry(t, handler, "/api/enemy-skills", enemyskills.SaveInput{
		ID:          "enemyskill_1",
		Name:        "Roar",
		Description: "Loud",
		Script:      "say roar",
	})

	createJSONEntry(t, handler, "/api/treasures", treasures.SaveInput{
		ID:        "treasure_1",
		TablePath: "minecraft:chests/simple_dungeon",
		LootPools: []treasures.DropRef{
			{Kind: "item", RefID: itemID, Weight: 1},
			{Kind: "grimoire", RefID: grimoireID, Weight: 1},
		},
	})

	createJSONEntry(t, handler, "/api/loottables", loottables.SaveInput{
		ID: "loottable_1",
		LootPools: []treasures.DropRef{
			{Kind: "item", RefID: itemID, Weight: 1},
			{Kind: "grimoire", RefID: grimoireID, Weight: 1},
		},
	})

	createJSONEntry(t, handler, "/api/enemies", enemies.SaveInput{
		ID:            "enemy_1",
		MobType:       "minecraft:zombie",
		Name:          "Sample Zombie",
		HP:            20,
		DropMode:      "replace",
		EnemySkillIDs: []string{enemySkillID},
		Drops: []enemies.DropRef{
			{Kind: "minecraft_item", RefID: "minecraft:rotten_flesh", Weight: 1},
		},
		Equipment: enemies.Equipment{
			Mainhand: &enemies.EquipmentSlot{
				Kind:  "minecraft_item",
				RefID: "minecraft:iron_sword",
				Count: 1,
			},
		},
	})

	rec := requestJSON(t, handler, http.MethodPost, "/api/save", struct{}{})
	if rec.Code != http.StatusOK {
		t.Fatalf("save status = %d body=%s", rec.Code, rec.Body.String())
	}

	checkFiles := []string{
		filepath.Join(root, "out", "data", "maf", "function", "generated", "skill", skillID+".mcfunction"),
		filepath.Join(root, "out", "data", "maf", "function", "generated", "grimoire", grimoireID+".mcfunction"),
		filepath.Join(root, "out", "data", "maf", "function", "generated", "grimoire", "selectexec.mcfunction"),
		filepath.Join(root, "out", "data", "maf", "function", "generated", "debug", "grimoire", grimoireID+".mcfunction"),
		filepath.Join(root, "out", "data", "minecraft", "loot_table", "chests", "simple_dungeon.json"),
		filepath.Join(root, "out", "data", "maf", "loot_table", "generated", "loottable", "loottable_1.json"),
		filepath.Join(root, "out", "data", "maf", "loot_table", "generated", "enemy", "enemy_1.json"),
	}
	for _, path := range checkFiles {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing exported file %s: %v", path, err)
		}
	}
}
