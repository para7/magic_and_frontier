package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"tools2/app/internal/application"
	"tools2/app/internal/config"
	"tools2/app/internal/domain/common"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/grimoire"
	"tools2/app/internal/domain/items"
	"tools2/app/internal/domain/loottables"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/spawntables"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/export"
)

func TestHandlerHealth(t *testing.T) {
	handler, _ := newTestHandler(t)
	rec := request(t, handler, http.MethodGet, "/health", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body map[string]bool
	decodeResponse(t, rec, &body)
	if !body["ok"] {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandlerSSRSkillCreate(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := request(t, handler, http.MethodGet, "/skills/new", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `action="/skills/new"`) {
		t.Fatalf("body = %s", rec.Body.String())
	}

	rec = postForm(t, handler, "/skills/new", url.Values{
		"id":          {"skill_1"},
		"name":        {"Slash"},
		"description": {"Basic skill"},
		"script":      {"say slash"},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/skills")

	rec = request(t, handler, http.MethodGet, "/skills", nil, "")
	if !strings.Contains(rec.Body.String(), "Slash") {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestHandlerSSRGrimoireEditShowsReadonlyCastID(t *testing.T) {
	handler, _ := newTestHandler(t)
	grimoireID := createJSONEntry(t, handler, "/api/grimoire", grimoire.SaveInput{
		ID:       "grimoire_1",
		CastID:   1,
		Title:    "Firebolt",
		Script:   "say fire",
		CastTime: 20,
		MPCost:   5,
	})

	rec := request(t, handler, http.MethodGet, "/grimoire/edit?id="+grimoireID, nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `name="castid" value="1" readonly`) {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestHandlerSSRGrimoireNewShowsEditableCastID(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := request(t, handler, http.MethodGet, "/grimoire/new", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `name="castid" value=""`) {
		t.Fatalf("body = %s", body)
	}
	if strings.Contains(body, `name="castid" value="" readonly`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, "1以上の整数") {
		t.Fatalf("body = %s", body)
	}
}

func TestHandlerSSRGrimoireCreateRejectsNonNumericCastID(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := postForm(t, handler, "/grimoire/new", url.Values{
		"id":       {"grimoire_1"},
		"castid":   {"abc"},
		"castTime": {"20"},
		"mpCost":   {"5"},
		"title":    {"Firebolt"},
		"script":   {"say fire"},
	}, http.StatusOK)
	body := rec.Body.String()
	if !strings.Contains(body, "Must be a number.") {
		t.Fatalf("body = %s", body)
	}
}

func TestHandlerSSRGrimoireCreateDuplicateIDDoesNotShowCastIDValidationError(t *testing.T) {
	handler, _ := newTestHandler(t)
	createJSONEntry(t, handler, "/api/grimoire", grimoire.SaveInput{
		ID:       "grimoire_1",
		CastID:   1,
		Title:    "Existing",
		Script:   "say existing",
		CastTime: 20,
		MPCost:   5,
	})

	rec := postForm(t, handler, "/grimoire/new", url.Values{
		"id":       {"grimoire_1"},
		"castid":   {"999"},
		"castTime": {"20"},
		"mpCost":   {"5"},
		"title":    {"Duplicate"},
		"script":   {"say dup"},
	}, http.StatusOK)
	body := rec.Body.String()
	if !strings.Contains(body, "この ID は既に使用されています。") {
		t.Fatalf("body = %s", body)
	}
	if strings.Contains(body, "Must satisfy gte 1.") {
		t.Fatalf("body = %s", body)
	}
	if strings.Contains(body, "Must be a number.") {
		t.Fatalf("body = %s", body)
	}
}

func TestHandlerSSRTreasureCreateValidationKeepsIDEditable(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := postForm(t, handler, "/treasures/edit", url.Values{
		"id":            {"bad"},
		"tablePath":     {"minecraft:chests/simple_dungeon"},
		"lootPoolsText": {"minecraft_item,minecraft:apple,1,,"},
	}, http.StatusOK)
	body := rec.Body.String()
	if !strings.Contains(body, `name="id" value="bad"`) {
		t.Fatalf("body = %s", body)
	}
	if strings.Contains(body, `name="id" value="bad" readonly`) {
		t.Fatalf("body = %s", body)
	}
}

func TestHandlerSSRItemsListIncludesClientControlsAndReturnTo(t *testing.T) {
	handler, _ := newTestHandler(t)
	skillID := createJSONEntry(t, handler, "/api/skills", skills.SaveInput{
		ID:     "skill_1",
		Name:   "Slash",
		Script: "say slash",
	})
	itemID := createJSONEntry(t, handler, "/api/items", items.SaveInput{
		ID:      "items_1",
		ItemID:  "minecraft:apple",
		SkillID: skillID,
	})

	rec := request(t, handler, http.MethodGet, "/items?q=apple&page=2", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `data-list-root`) || !strings.Contains(body, `data-list-search`) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `href="/items/new?returnTo=%2Fitems%3Fq%3Dapple%26page%3D2"`) {
		t.Fatalf("body = %s", body)
	}
	wantEdit := `href="/items/edit?id=` + itemID + `&amp;returnTo=%2Fitems%3Fq%3Dapple%26page%3D2"`
	if !strings.Contains(body, wantEdit) {
		t.Fatalf("body = %s", body)
	}
	if !strings.Contains(body, `data-sort-item_id="minecraft:apple"`) {
		t.Fatalf("body = %s", body)
	}
}

func TestHandlerSSRSkillEditRespectsReturnToOnSaveAndFallback(t *testing.T) {
	handler, _ := newTestHandler(t)
	skillID := createJSONEntry(t, handler, "/api/skills", skills.SaveInput{
		ID:          "skill_1",
		Name:        "Slash",
		Description: "Basic slash",
		Script:      "say slash",
	})

	rec := request(t, handler, http.MethodGet, "/skills/edit?id="+skillID+"&returnTo=%2Fskills%3Fq%3Dslash%26page%3D2", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `href="/skills?q=slash&amp;page=2">Back to list</a>`) {
		t.Fatalf("body = %s", rec.Body.String())
	}

	rec = postForm(t, handler, "/skills/edit", url.Values{
		"id":          {skillID},
		"name":        {"Slash v2"},
		"description": {"Updated"},
		"script":      {"say slash2"},
		"returnTo":    {"/skills?q=slash&page=2"},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/skills?q=slash&page=2")

	rec = postForm(t, handler, "/skills/edit", url.Values{
		"id":          {skillID},
		"name":        {"Slash v3"},
		"description": {"Updated again"},
		"script":      {"say slash3"},
		"returnTo":    {"https://evil.example"},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/skills")
}

func TestHandlerSSRItemReturnToRejectsNonListPaths(t *testing.T) {
	handler, _ := newTestHandler(t)
	itemID := createJSONEntry(t, handler, "/api/items", items.SaveInput{
		ID:     "items_1",
		ItemID: "minecraft:apple",
	})

	rec := request(t, handler, http.MethodGet, "/items/edit?id="+itemID+"&returnTo=%2Fsave", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `href="/items">Back to list</a>`) {
		t.Fatalf("body = %s", rec.Body.String())
	}

	rec = postForm(t, handler, "/items/edit", url.Values{
		"id":       {itemID},
		"itemId":   {"minecraft:apple"},
		"returnTo": {"/items/edit?id=" + itemID},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/items")
}

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
		CastID:      100,
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

func TestHandlerAPIGrimoireUsesClientProvidedCastID(t *testing.T) {
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
	if body.Entry.CastID != 999 {
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

func TestHandlerAPIGrimoireRejectsDuplicateCastID(t *testing.T) {
	handler, _ := newTestHandler(t)

	createJSONEntry(t, handler, "/api/grimoire", grimoire.SaveInput{
		ID:       "grimoire_1",
		CastID:   7,
		Title:    "Firebolt",
		Script:   "say fire",
		CastTime: 20,
		MPCost:   5,
	})

	rec := requestJSON(t, handler, http.MethodPost, "/api/grimoire", grimoire.SaveInput{
		ID:       "grimoire_2",
		CastID:   7,
		Title:    "Icebolt",
		Script:   "say ice",
		CastTime: 20,
		MPCost:   5,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var body common.SaveResult[grimoire.GrimoireEntry]
	decodeResponse(t, rec, &body)
	if body.FieldErrors["castid"] == "" {
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

type entryIDOnly struct {
	ID string `json:"id"`
}

func createJSONEntry[T any](t *testing.T, handler http.Handler, path string, payload T) string {
	t.Helper()
	rec := requestJSON(t, handler, http.MethodPost, path, payload)
	if rec.Code != http.StatusOK {
		t.Fatalf("%s status = %d body=%s", path, rec.Code, rec.Body.String())
	}
	var body common.SaveResult[entryIDOnly]
	decodeResponse(t, rec, &body)
	if body.Entry == nil {
		t.Fatalf("%s response missing entry: %#v", path, body)
	}
	return body.Entry.ID
}

func newTestHandler(t *testing.T) (http.Handler, string) {
	t.Helper()
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	templatePath := filepath.Join(root, "pack-template.mcmeta")
	if err := os.WriteFile(templatePath, []byte("{\"pack\":{\"pack_format\":61,\"description\":\"test\"}}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	settings := export.ExportSettings{
		OutputRoot:       "./out",
		Namespace:        "maf",
		TemplatePackPath: "./pack-template.mcmeta",
		Paths: export.ExportPaths{
			ItemFunctionDir:       "data/maf/function/generated/item",
			ItemLootDir:           "data/maf/loot_table/generated/item",
			SpellFunctionDir:      "data/maf/function/generated/grimoire",
			SpellLootDir:          "data/maf/loot_table/generated/grimoire",
			SkillFunctionDir:      "data/maf/function/generated/skill",
			EnemySkillFunctionDir: "data/maf/function/generated/enemy_skill",
			EnemyFunctionDir:      "data/maf/function/generated/enemy",
			EnemyLootDir:          "data/maf/loot_table/generated/enemy",
			TreasureLootDir:       "data/maf/loot_table/generated/treasure",
			LoottableLootDir:      "data/maf/loot_table/generated/loottable",
			DebugFunctionDir:      "data/maf/function/debug/give",
			MinecraftTagDir:       "data/minecraft/tags/function",
		},
	}
	writeJSONFile(t, settingsPath, settings)

	cfg := config.Config{
		Port:                   8787,
		ItemStatePath:          filepath.Join(root, "item.json"),
		GrimoireStatePath:      filepath.Join(root, "grimoire.json"),
		SkillStatePath:         filepath.Join(root, "skill.json"),
		EnemySkillStatePath:    filepath.Join(root, "enemy_skill.json"),
		EnemyStatePath:         filepath.Join(root, "enemy.json"),
		SpawnTableStatePath:    filepath.Join(root, "spawn_table.json"),
		TreasureStatePath:      filepath.Join(root, "treasure.json"),
		LootTablesStatePath:    filepath.Join(root, "loottables.json"),
		ExportSettingsPath:     settingsPath,
		MinecraftLootTableRoot: writeMinecraftLootTableRoot(t, root),
	}
	handler := NewHandler(cfg, application.Dependencies{
		Now: func() time.Time { return time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC) },
	})
	return handler, root
}

func request(t *testing.T, handler http.Handler, method, path string, body *bytes.Reader, contentType string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		reader = body
	}
	req := httptest.NewRequest(method, path, reader)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func requestJSON[T any](t *testing.T, handler http.Handler, method, path string, payload T) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	return request(t, handler, method, path, bytes.NewReader(data), "application/json")
}

func postForm(t *testing.T, handler http.Handler, path string, values url.Values, wantStatus int) *httptest.ResponseRecorder {
	t.Helper()
	rec := request(t, handler, http.MethodPost, path, bytes.NewReader([]byte(values.Encode())), "application/x-www-form-urlencoded")
	if rec.Code != wantStatus {
		t.Fatalf("%s status = %d, want %d, body=%s", path, rec.Code, wantStatus, rec.Body.String())
	}
	return rec
}

func assertRedirect(t *testing.T, rec *httptest.ResponseRecorder, wantLocation string) {
	t.Helper()
	if got := rec.Header().Get("Location"); got != wantLocation {
		t.Fatalf("location = %q, want %q", got, wantLocation)
	}
}

func decodeResponse[T any](t *testing.T, rec *httptest.ResponseRecorder, dest *T) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), dest); err != nil {
		t.Fatalf("decode response: %v body=%s", err, rec.Body.String())
	}
}

func writeJSONFile[T any](t *testing.T, path string, value T) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeMinecraftLootTableRoot(t *testing.T, root string) string {
	t.Helper()
	dir := filepath.Join(root, "minecraft", "1.21.11", "loot_table", "chests")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "simple_dungeon.json"), []byte("{\"type\":\"minecraft:generic\",\"pools\":[]}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return filepath.Join(root, "minecraft", "1.21.11", "loot_table")
}
