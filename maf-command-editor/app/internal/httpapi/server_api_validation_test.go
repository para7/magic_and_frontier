package httpapi

import (
	"net/http"
	"testing"

	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
)

func TestHandlerAPITreasureRejectsDuplicateTablePath(t *testing.T) {
	handler, _ := newTestHandler(t)

	createJSONEntry(t, handler, "/api/treasures", treasures.SaveInput{
		ID:        "treasure_1",
		TablePath: "minecraft:chests/simple_dungeon",
		LootPools: []treasures.DropRef{{Kind: "minecraft_item", RefID: "minecraft:apple", Weight: 1}},
	})

	rec := requestJSON(t, handler, http.MethodPost, "/api/treasures", treasures.SaveInput{
		ID:        "treasure_2",
		TablePath: "minecraft:chests/simple_dungeon",
		LootPools: []treasures.DropRef{{Kind: "minecraft_item", RefID: "minecraft:stick", Weight: 1}},
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var body common.SaveResult[treasures.TreasureEntry]
	decodeResponse(t, rec, &body)
	if body.FieldErrors["tablePath"] == "" {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandlerAPITreasureRejectsMissingVanillaSource(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := requestJSON(t, handler, http.MethodPost, "/api/treasures", treasures.SaveInput{
		ID:        "treasure_1",
		TablePath: "minecraft:chests/missing_entry",
		LootPools: []treasures.DropRef{{Kind: "minecraft_item", RefID: "minecraft:apple", Weight: 1}},
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandlerAPIGrimoireUsesServerManagedCastID(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := requestJSON(t, handler, http.MethodPost, "/api/grimoire", grimoire.SaveInput{
		ID:       "grimoire_999",
		CastID:   999,
		Title:    "Firebolt",
		Script:   "say fire",
		CastTime: 20,
		MPCost:   5,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var body common.SaveResult[grimoire.GrimoireEntry]
	decodeResponse(t, rec, &body)
	if body.Entry == nil {
		t.Fatalf("body = %#v", body)
	}
	if body.Entry.ID != "grimoire_999" {
		t.Fatalf("entry = %#v", body.Entry)
	}
	if body.Entry.CastID != 1 {
		t.Fatalf("entry = %#v", body.Entry)
	}

	rec = requestJSON(t, handler, http.MethodPost, "/api/grimoire", grimoire.SaveInput{
		ID:       "grimoire_999",
		CastID:   555,
		Title:    "Firebolt v2",
		Script:   "say fire2",
		CastTime: 30,
		MPCost:   9,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	body = common.SaveResult[grimoire.GrimoireEntry]{}
	decodeResponse(t, rec, &body)
	if body.Entry != nil {
		t.Fatalf("body = %#v", body)
	}
	if body.FieldErrors["id"] == "" {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandlerAPISkillRejectsDuplicateID(t *testing.T) {
	handler, _ := newTestHandler(t)

	createJSONEntry(t, handler, "/api/skills", skills.SaveInput{
		ID:     "skill_1",
		Name:   "Slash",
		Script: "say slash",
	})

	rec := requestJSON(t, handler, http.MethodPost, "/api/skills", skills.SaveInput{
		ID:     "skill_1",
		Name:   "Slash v2",
		Script: "say slash2",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}

	var body common.SaveResult[skills.SkillEntry]
	decodeResponse(t, rec, &body)
	if body.FieldErrors["id"] == "" {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandlerAPISkillDeleteRejectsReferencedItem(t *testing.T) {
	handler, _ := newTestHandler(t)

	skillID := createJSONEntry(t, handler, "/api/skills", skills.SaveInput{
		ID:          "skill_1",
		Name:        "Slash",
		Description: "Basic slash",
		Script:      "say slash",
	})
	createJSONEntry(t, handler, "/api/items", items.SaveInput{
		ID:      "items_1",
		ItemID:  "minecraft:apple",
		SkillID: skillID,
	})

	rec := request(t, handler, http.MethodDelete, "/api/skills/"+skillID, nil, "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}

	var body common.DeleteResult
	decodeResponse(t, rec, &body)
	if body.Code != "REFERENCE_ERROR" {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandlerAPISpawnTableRejectsOverlap(t *testing.T) {
	handler, _ := newTestHandler(t)

	enemyID := createJSONEntry(t, handler, "/api/enemies", enemies.SaveInput{
		ID:       "enemy_1",
		MobType:  "minecraft:zombie",
		Name:     "Zombie",
		HP:       20,
		DropMode: "replace",
		Drops:    []enemies.DropRef{{Kind: "minecraft_item", RefID: "minecraft:rotten_flesh", Weight: 1}},
	})

	createJSONEntry(t, handler, "/api/spawn-tables", spawntables.SaveInput{
		ID:            "spawntable_1",
		SourceMobType: "minecraft:zombie",
		Dimension:     "minecraft:overworld",
		MinX:          0,
		MaxX:          100,
		MinY:          -64,
		MaxY:          320,
		MinZ:          0,
		MaxZ:          100,
		BaseMobWeight: 8000,
		Replacements:  []spawntables.ReplacementEntry{{EnemyID: enemyID, Weight: 2000}},
	})

	rec := requestJSON(t, handler, http.MethodPost, "/api/spawn-tables", spawntables.SaveInput{
		ID:            "spawntable_2",
		SourceMobType: "minecraft:zombie",
		Dimension:     "minecraft:overworld",
		MinX:          50,
		MaxX:          150,
		MinY:          -64,
		MaxY:          320,
		MinZ:          50,
		MaxZ:          150,
		BaseMobWeight: 9000,
		Replacements:  []spawntables.ReplacementEntry{{EnemyID: enemyID, Weight: 1000}},
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}

	var body common.SaveResult[spawntables.SpawnTableEntry]
	decodeResponse(t, rec, &body)
	if body.FieldErrors["range"] == "" {
		t.Fatalf("body = %#v", body)
	}
}
