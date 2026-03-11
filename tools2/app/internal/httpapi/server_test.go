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

	"tools2/app/internal/config"
	"tools2/app/internal/domain/enemies"
	"tools2/app/internal/domain/enemyskills"
	"tools2/app/internal/domain/skills"
	"tools2/app/internal/domain/treasures"
	"tools2/app/internal/store"
)

func TestHandler_Health(t *testing.T) {
	handler, _ := newTestHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body map[string]bool
	decodeResponse(t, rec, &body)
	if !body["ok"] {
		t.Fatalf("body = %#v, want ok=true", body)
	}
}

func TestHandler_SSRPageLoadsAllFeatures(t *testing.T) {
	handler, _ := newTestHandler(t)

	tests := []struct {
		path      string
		text      string
		createURL string
	}{
		{path: "/items", text: "Items", createURL: "/items/new"},
		{path: "/grimoire", text: "Grimoire", createURL: "/grimoire/new"},
		{path: "/skills", text: "Skills", createURL: "/skills/new"},
		{path: "/enemy-skills", text: "Enemy Skills", createURL: "/enemy-skills/new"},
		{path: "/treasures", text: "Treasures", createURL: "/treasures/new"},
		{path: "/enemies", text: "Enemies", createURL: "/enemies/new"},
	}

	for _, tt := range tests {
		rec := request(t, handler, http.MethodGet, tt.path, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want %d", tt.path, rec.Code, http.StatusOK)
		}
		body := rec.Body.String()
		if !strings.Contains(body, tt.text) {
			t.Fatalf("%s body missing %q", tt.path, tt.text)
		}
		if !strings.Contains(body, tt.createURL) {
			t.Fatalf("%s body missing create link %q", tt.path, tt.createURL)
		}
	}
}

func TestHandler_SSRNewPageLoadsForm(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := request(t, handler, http.MethodGet, "/items/new", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `action="/items/new"`) {
		t.Fatalf("body missing new form action: %s", body)
	}
	if !strings.Contains(body, "Back to list") {
		t.Fatalf("body missing back link: %s", body)
	}
}

func TestHandler_SSREditPageLoadsExistingEntry(t *testing.T) {
	handler, _ := newTestHandler(t)
	itemID := uuid("000000000301")

	rec := postForm(t, handler, "/items/new", url.Values{
		"id":                 {itemID},
		"itemId":             {"minecraft:apple"},
		"count":              {"2"},
		"customName":         {"Starter Apple"},
		"lore":               {"fresh"},
		"enchantments":       {""},
		"customNbt":          {""},
		"dropTableId":        {""},
		"unbreakable":        {""},
		"repairCost":         {""},
		"hideFlags":          {""},
		"attributeModifiers": {""},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/items")

	rec = request(t, handler, http.MethodGet, "/items/edit?id="+itemID, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `action="/items/edit"`) {
		t.Fatalf("body missing edit form action: %s", body)
	}
	if !strings.Contains(body, "Starter Apple") {
		t.Fatalf("body missing existing entry values: %s", body)
	}
}

func TestHandler_SSREditSubmitUpdatesExistingEntry(t *testing.T) {
	handler, _ := newTestHandler(t)
	itemID := uuid("000000000302")

	rec := postForm(t, handler, "/items/new", url.Values{
		"id":                 {itemID},
		"itemId":             {"minecraft:apple"},
		"count":              {"2"},
		"customName":         {"Starter Apple"},
		"lore":               {"fresh"},
		"enchantments":       {""},
		"customNbt":          {""},
		"dropTableId":        {""},
		"unbreakable":        {""},
		"repairCost":         {""},
		"hideFlags":          {""},
		"attributeModifiers": {""},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/items")

	rec = postForm(t, handler, "/items/edit", url.Values{
		"id":                 {itemID},
		"itemId":             {"minecraft:golden_apple"},
		"count":              {"4"},
		"customName":         {"Edited Apple"},
		"lore":               {"fresh"},
		"enchantments":       {""},
		"customNbt":          {""},
		"dropTableId":        {""},
		"unbreakable":        {""},
		"repairCost":         {""},
		"hideFlags":          {""},
		"attributeModifiers": {""},
	}, http.StatusSeeOther)
	assertRedirect(t, rec, "/items")

	rec = request(t, handler, http.MethodGet, "/items", nil)
	body := rec.Body.String()
	if !strings.Contains(body, "minecraft:golden_apple") {
		t.Fatalf("edit body = %s", body)
	}
	if strings.Contains(body, "minecraft:apple x2") {
		t.Fatalf("body still contains stale value: %s", body)
	}
}

func TestHandler_SSRCRUDHappyPathAndSave(t *testing.T) {
	handler, _ := newTestHandler(t)

	itemID := uuid("000000000141")
	grimoireID := uuid("000000000142")
	skillID := uuid("000000000143")
	enemySkillID := uuid("000000000144")
	treasureID := uuid("000000000145")
	enemyID := uuid("000000000146")

	itemRec := postForm(t, handler, "/items/new", url.Values{
		"id":                 {itemID},
		"itemId":             {"minecraft:apple"},
		"count":              {"2"},
		"customName":         {"Starter Apple"},
		"lore":               {"fresh"},
		"enchantments":       {""},
		"customNbt":          {""},
		"dropTableId":        {""},
		"unbreakable":        {""},
		"repairCost":         {""},
		"hideFlags":          {""},
		"attributeModifiers": {""},
	}, http.StatusSeeOther)
	assertRedirect(t, itemRec, "/items")

	grimoireRec := postForm(t, handler, "/grimoire/new", url.Values{
		"id":           {grimoireID},
		"castid":       {"2"},
		"title":        {"Apple Spell"},
		"description":  {"desc"},
		"script":       {"function maf:spell/apple"},
		"variantsText": {"6,40"},
	}, http.StatusSeeOther)
	assertRedirect(t, grimoireRec, "/grimoire")

	skillRec := postForm(t, handler, "/skills/new", url.Values{
		"id":     {skillID},
		"name":   {"Slash"},
		"script": {"say slash"},
		"itemId": {itemID},
	}, http.StatusSeeOther)
	assertRedirect(t, skillRec, "/skills")

	enemySkillRec := postForm(t, handler, "/enemy-skills/new", url.Values{
		"id":       {enemySkillID},
		"name":     {"Roar"},
		"script":   {"say roar"},
		"cooldown": {"20"},
		"trigger":  {"on_spawn"},
	}, http.StatusSeeOther)
	assertRedirect(t, enemySkillRec, "/enemy-skills")

	treasureRec := postForm(t, handler, "/treasures/new", url.Values{
		"id":            {treasureID},
		"name":          {"Starter Treasure"},
		"lootPoolsText": {"item," + itemID + ",3,1,1\ngrimoire," + grimoireID + ",1,1,1"},
	}, http.StatusSeeOther)
	assertRedirect(t, treasureRec, "/treasures")

	enemyRec := postForm(t, handler, "/enemies/new", url.Values{
		"id":            {enemyID},
		"name":          {"Zombie"},
		"hp":            {"20"},
		"dropTableId":   {treasureID},
		"enemySkillIds": {enemySkillID},
		"originX":       {"0"},
		"originY":       {"64"},
		"originZ":       {"0"},
		"distanceMin":   {"0"},
		"distanceMax":   {"32"},
	}, http.StatusSeeOther)
	assertRedirect(t, enemyRec, "/enemies")

	if rec := request(t, handler, http.MethodGet, "/items", nil); !strings.Contains(rec.Body.String(), "minecraft:apple") {
		t.Fatalf("item body = %s", rec.Body.String())
	}
	if rec := request(t, handler, http.MethodGet, "/grimoire", nil); !strings.Contains(rec.Body.String(), "Apple Spell") {
		t.Fatalf("grimoire body = %s", rec.Body.String())
	}
	if rec := request(t, handler, http.MethodGet, "/skills", nil); !strings.Contains(rec.Body.String(), "Slash") {
		t.Fatalf("skill body = %s", rec.Body.String())
	}
	if rec := request(t, handler, http.MethodGet, "/enemy-skills", nil); !strings.Contains(rec.Body.String(), "Roar") {
		t.Fatalf("enemy skill body = %s", rec.Body.String())
	}
	if rec := request(t, handler, http.MethodGet, "/treasures", nil); !strings.Contains(rec.Body.String(), "Starter Treasure") {
		t.Fatalf("treasure body = %s", rec.Body.String())
	}
	if rec := request(t, handler, http.MethodGet, "/enemies", nil); !strings.Contains(rec.Body.String(), "Zombie") {
		t.Fatalf("enemy body = %s", rec.Body.String())
	}

	saveRec := requestForm(t, handler, http.MethodPost, "/save", url.Values{
		"currentPath": {"/items"},
	})
	if saveRec.Code != http.StatusOK {
		t.Fatalf("save status = %d, want %d, body=%s", saveRec.Code, http.StatusOK, saveRec.Body.String())
	}
	if !strings.Contains(saveRec.Body.String(), "Exported") {
		t.Fatalf("save body = %s", saveRec.Body.String())
	}
	if !strings.Contains(saveRec.Body.String(), "<!doctype html>") || !strings.Contains(saveRec.Body.String(), "Save / Export") {
		t.Fatalf("save body should be full page, got %s", saveRec.Body.String())
	}

	deleteRec := requestForm(t, handler, http.MethodPost, "/enemies/"+enemyID+"/delete", url.Values{})
	assertRedirect(t, deleteRec, "/enemies")
	deleteRec = requestForm(t, handler, http.MethodPost, "/enemy-skills/"+enemySkillID+"/delete", url.Values{})
	assertRedirect(t, deleteRec, "/enemy-skills")
	deleteRec = requestForm(t, handler, http.MethodPost, "/treasures/"+treasureID+"/delete", url.Values{})
	assertRedirect(t, deleteRec, "/treasures")
	deleteRec = requestForm(t, handler, http.MethodPost, "/skills/"+skillID+"/delete", url.Values{})
	assertRedirect(t, deleteRec, "/skills")
	deleteRec = requestForm(t, handler, http.MethodPost, "/grimoire/"+grimoireID+"/delete", url.Values{})
	assertRedirect(t, deleteRec, "/grimoire")
	deleteRec = requestForm(t, handler, http.MethodPost, "/items/"+itemID+"/delete", url.Values{})
	assertRedirect(t, deleteRec, "/items")
}

func TestHandler_SaveNonHTMXReturnsCurrentFullPage(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := requestForm(t, handler, http.MethodPost, "/save", url.Values{
		"currentPath": {"/grimoire"},
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected full HTML page, got %s", body)
	}
	if !strings.Contains(body, "Grimoire") {
		t.Fatalf("expected grimoire page, got %s", body)
	}
	if strings.Contains(body, `<div class="notice`) && !strings.Contains(body, "Exported") && !strings.Contains(body, "Invalid export settings.") {
		t.Fatalf("unexpected notice content: %s", body)
	}
}

func TestHandler_SaveHTMXReturnsFragment(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/save", strings.NewReader(url.Values{
		"currentPath": {"/items"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := rec.Body.String()
	if strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected fragment, got full page: %s", body)
	}
	if !strings.Contains(body, "Exported") {
		t.Fatalf("expected export notice, got %s", body)
	}
}

func TestHandler_SaveNonHTMXFallsBackToItemsPage(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := requestForm(t, handler, http.MethodPost, "/save", url.Values{
		"currentPath": {"/not-a-screen"},
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "<!doctype html>") || !strings.Contains(body, "Items") {
		t.Fatalf("expected items full page fallback, got %s", body)
	}
}

func TestHandler_SSRValidationShowsFieldError(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := postForm(t, handler, "/skills/new", url.Values{
		"id":     {uuid("000000000201")},
		"name":   {"Slash"},
		"script": {"say slash"},
		"itemId": {uuid("000000000777")},
	}, http.StatusOK)

	if !strings.Contains(rec.Body.String(), "Referenced item") {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func assertRedirect(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusSeeOther, rec.Body.String())
	}
	if got := rec.Header().Get("Location"); got != want {
		t.Fatalf("location = %q, want %q", got, want)
	}
}

func TestHandler_CRUDHappyPathAndSave(t *testing.T) {
	handler, root := newTestHandler(t)

	itemID := uuid("000000000041")
	grimoireID := uuid("000000000042")
	skillID := uuid("000000000043")
	enemySkillID := uuid("000000000044")
	treasureID := uuid("000000000045")
	enemyID := uuid("000000000046")

	postJSON(t, handler, "/api/items", map[string]any{
		"id":                  itemID,
		"itemId":              "minecraft:apple",
		"count":               2,
		"customName":          "",
		"lore":                "",
		"enchantments":        "",
		"unbreakable":         false,
		"customModelData":     "",
		"repairCost":          "",
		"hideFlags":           "",
		"potionId":            "",
		"customPotionColor":   "",
		"customPotionEffects": "",
		"attributeModifiers":  "",
		"customNbt":           "",
	}, http.StatusOK)

	postJSON(t, handler, "/api/grimoire", map[string]any{
		"id":          grimoireID,
		"castid":      2,
		"script":      "function maf:spell/apple",
		"title":       "Apple Spell",
		"description": "desc",
		"variants":    []map[string]any{{"cast": 6, "cost": 40}},
	}, http.StatusOK)

	postJSON(t, handler, "/api/skills", map[string]any{
		"id":     skillID,
		"name":   "Slash",
		"script": "say slash",
		"itemId": itemID,
	}, http.StatusOK)

	postJSON(t, handler, "/api/enemy-skills", map[string]any{
		"id":       enemySkillID,
		"name":     "Roar",
		"script":   "say roar",
		"cooldown": 20,
		"trigger":  "on_spawn",
	}, http.StatusOK)

	postJSON(t, handler, "/api/treasures", map[string]any{
		"id":   treasureID,
		"name": "Starter Treasure",
		"lootPools": []map[string]any{
			{"kind": "item", "refId": itemID, "weight": 3},
			{"kind": "grimoire", "refId": grimoireID, "weight": 1},
		},
	}, http.StatusOK)

	postJSON(t, handler, "/api/enemies", map[string]any{
		"id":            enemyID,
		"name":          "Zombie",
		"hp":            20,
		"dropTableId":   treasureID,
		"enemySkillIds": []string{enemySkillID},
		"spawnRule": map[string]any{
			"origin":   map[string]any{"x": 0, "y": 64, "z": 0},
			"distance": map[string]any{"min": 0, "max": 32},
		},
	}, http.StatusOK)

	assertEntriesLen(t, handler, "/api/items", "items", 1)
	assertEntriesLen(t, handler, "/api/grimoire", "entries", 1)
	assertEntriesLen(t, handler, "/api/skills", "entries", 1)
	assertEntriesLen(t, handler, "/api/enemy-skills", "entries", 1)
	assertEntriesLen(t, handler, "/api/treasures", "entries", 1)
	assertEntriesLen(t, handler, "/api/enemies", "entries", 1)

	saveRec := request(t, handler, http.MethodPost, "/api/save", nil)
	if saveRec.Code != http.StatusOK {
		t.Fatalf("save status = %d, want %d, body=%s", saveRec.Code, http.StatusOK, saveRec.Body.String())
	}
	var saveBody struct {
		OK        bool `json:"ok"`
		Generated struct {
			ItemFunctions       int `json:"itemFunctions"`
			ItemLootTables      int `json:"itemLootTables"`
			SpellFunctions      int `json:"spellFunctions"`
			SpellLootTables     int `json:"spellLootTables"`
			SkillFunctions      int `json:"skillFunctions"`
			EnemySkillFunctions int `json:"enemySkillFunctions"`
			EnemyFunctions      int `json:"enemyFunctions"`
			EnemyLootTables     int `json:"enemyLootTables"`
			TreasureLootTables  int `json:"treasureLootTables"`
			TotalFiles          int `json:"totalFiles"`
		} `json:"generated"`
		OutputRoot string `json:"outputRoot"`
	}
	decodeResponse(t, saveRec, &saveBody)
	if !saveBody.OK {
		t.Fatalf("save response = %s", saveRec.Body.String())
	}
	if saveBody.Generated.ItemFunctions != 1 || saveBody.Generated.SpellLootTables != 1 || saveBody.Generated.SkillFunctions != 1 || saveBody.Generated.EnemySkillFunctions != 1 || saveBody.Generated.EnemyFunctions != 1 || saveBody.Generated.TreasureLootTables != 1 {
		t.Fatalf("unexpected generated counts: %+v", saveBody.Generated)
	}
	if _, err := os.Stat(filepath.Join(root, "out", "data", "maf", "function", "item", itemID+".mcfunction")); err != nil {
		t.Fatalf("item function missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "out", "data", "maf", "loot_table", "treasure", treasureID+".json")); err != nil {
		t.Fatalf("treasure loot missing: %v", err)
	}

	deleteOK(t, handler, "/api/enemies/"+enemyID)
	deleteOK(t, handler, "/api/enemy-skills/"+enemySkillID)
	deleteOK(t, handler, "/api/treasures/"+treasureID)
	deleteOK(t, handler, "/api/skills/"+skillID)
	deleteOK(t, handler, "/api/grimoire/"+grimoireID)
	deleteOK(t, handler, "/api/items/"+itemID)
}

func TestHandler_SkillsRejectUnknownItem(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := postJSON(t, handler, "/api/skills", map[string]any{
		"id":     uuid("000000000001"),
		"name":   "Slash",
		"script": "say slash",
		"itemId": uuid("000000000777"),
	}, http.StatusBadRequest)

	var body map[string]any
	decodeResponse(t, rec, &body)
	fieldErrors := body["fieldErrors"].(map[string]any)
	if !strings.Contains(fieldErrors["itemId"].(string), "Referenced item") {
		t.Fatalf("fieldErrors = %#v", fieldErrors)
	}
}

func TestHandler_EnemiesRejectUnknownEnemySkill(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := postJSON(t, handler, "/api/enemies", map[string]any{
		"id":            uuid("000000000002"),
		"name":          "Zombie",
		"hp":            20,
		"dropTableId":   "treasure-a",
		"enemySkillIds": []string{uuid("000000000888")},
		"spawnRule": map[string]any{
			"origin":   map[string]any{"x": 0, "y": 64, "z": 0},
			"distance": map[string]any{"min": 0, "max": 32},
		},
	}, http.StatusBadRequest)

	var body map[string]any
	decodeResponse(t, rec, &body)
	fieldErrors := body["fieldErrors"].(map[string]any)
	if !strings.Contains(fieldErrors["enemySkillIds.0"].(string), "Referenced enemy skill") {
		t.Fatalf("fieldErrors = %#v", fieldErrors)
	}
}

func TestHandler_EnemiesRejectInvalidDistance(t *testing.T) {
	handler, _ := newTestHandler(t)
	enemySkillID := uuid("000000000003")

	postJSON(t, handler, "/api/enemy-skills", map[string]any{
		"id":     enemySkillID,
		"name":   "Roar",
		"script": "say roar",
	}, http.StatusOK)

	rec := postJSON(t, handler, "/api/enemies", map[string]any{
		"id":            uuid("000000000004"),
		"name":          "Zombie",
		"hp":            20,
		"dropTableId":   "treasure-a",
		"enemySkillIds": []string{enemySkillID},
		"spawnRule": map[string]any{
			"origin":   map[string]any{"x": 0, "y": 64, "z": 0},
			"distance": map[string]any{"min": 50, "max": 10},
		},
	}, http.StatusBadRequest)

	var body map[string]any
	decodeResponse(t, rec, &body)
	fieldErrors := body["fieldErrors"].(map[string]any)
	if !strings.Contains(fieldErrors["spawnRule.distance.min"].(string), "Must be <=") {
		t.Fatalf("fieldErrors = %#v", fieldErrors)
	}
}

func TestHandler_DeleteEnemySkillRejectsReference(t *testing.T) {
	handler, _ := newTestHandler(t)
	itemID := uuid("000000000011")
	enemySkillID := uuid("000000000012")
	treasureID := uuid("000000000013")

	postJSON(t, handler, "/api/items", map[string]any{
		"id":                  itemID,
		"itemId":              "minecraft:apple",
		"count":               1,
		"customName":          "",
		"lore":                "",
		"enchantments":        "",
		"unbreakable":         false,
		"customModelData":     "",
		"repairCost":          "",
		"hideFlags":           "",
		"potionId":            "",
		"customPotionColor":   "",
		"customPotionEffects": "",
		"attributeModifiers":  "",
		"customNbt":           "",
	}, http.StatusOK)
	postJSON(t, handler, "/api/enemy-skills", map[string]any{
		"id":     enemySkillID,
		"name":   "Roar",
		"script": "say roar",
	}, http.StatusOK)
	postJSON(t, handler, "/api/treasures", map[string]any{
		"id":   treasureID,
		"name": "Starter Treasure",
		"lootPools": []map[string]any{
			{"kind": "item", "refId": itemID, "weight": 1},
		},
	}, http.StatusOK)
	postJSON(t, handler, "/api/enemies", map[string]any{
		"id":            uuid("000000000014"),
		"name":          "Zombie",
		"hp":            20,
		"dropTableId":   treasureID,
		"enemySkillIds": []string{enemySkillID},
		"spawnRule": map[string]any{
			"origin":   map[string]any{"x": 0, "y": 64, "z": 0},
			"distance": map[string]any{"min": 0, "max": 32},
		},
	}, http.StatusOK)

	rec := request(t, handler, http.MethodDelete, "/api/enemy-skills/"+enemySkillID, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	var body map[string]any
	decodeResponse(t, rec, &body)
	if body["code"] != "REFERENCE_ERROR" {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandler_SaveRejectsInvalidConfig(t *testing.T) {
	root := t.TempDir()
	cfg := config.Config{
		Port:                8787,
		ItemStatePath:       filepath.Join(root, "item-state.json"),
		GrimoireStatePath:   filepath.Join(root, "grimoire-state.json"),
		SkillStatePath:      filepath.Join(root, "skill-state.json"),
		EnemySkillStatePath: filepath.Join(root, "enemy-skill-state.json"),
		EnemyStatePath:      filepath.Join(root, "enemy-state.json"),
		TreasureStatePath:   filepath.Join(root, "treasure-state.json"),
		ExportSettingsPath:  filepath.Join(root, "broken-export-settings.json"),
	}
	if err := os.WriteFile(cfg.ExportSettingsPath, []byte("{\"outputRoot\":1}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(cfg, Dependencies{
		ItemRepo:       store.NewItemStateRepository(cfg.ItemStatePath),
		GrimoireRepo:   store.NewGrimoireStateRepository(cfg.GrimoireStatePath),
		SkillRepo:      store.NewEntryStateRepository[skills.SkillEntry](cfg.SkillStatePath),
		EnemySkillRepo: store.NewEntryStateRepository[enemyskills.EnemySkillEntry](cfg.EnemySkillStatePath),
		EnemyRepo:      store.NewEntryStateRepository[enemies.EnemyEntry](cfg.EnemyStatePath),
		TreasureRepo:   store.NewEntryStateRepository[treasures.TreasureEntry](cfg.TreasureStatePath),
		Now:            fixedNow,
	})

	rec := request(t, handler, http.MethodPost, "/api/save", nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	var body map[string]any
	decodeResponse(t, rec, &body)
	if body["code"] != "INVALID_CONFIG" {
		t.Fatalf("body = %#v", body)
	}
}

func newTestHandler(t *testing.T) (http.Handler, string) {
	t.Helper()
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export-settings.json")
	templatePath := filepath.Join(root, "pack-template.mcmeta")
	if err := os.WriteFile(templatePath, []byte("{\"pack\":{\"pack_format\":61,\"description\":\"test\"}}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	settings := map[string]any{
		"outputRoot":       "./out",
		"namespace":        "maf",
		"templatePackPath": "./pack-template.mcmeta",
		"paths": map[string]any{
			"itemFunctionDir":       "data/maf/function/item",
			"itemLootDir":           "data/maf/loot_table/item",
			"spellFunctionDir":      "data/maf/function/grimoire",
			"spellLootDir":          "data/maf/loot_table/grimoire",
			"skillFunctionDir":      "data/maf/function/skill",
			"enemySkillFunctionDir": "data/maf/function/enemy_skill",
			"enemyFunctionDir":      "data/maf/function/enemy/spawn",
			"enemyLootDir":          "data/maf/loot_table/enemy",
			"treasureLootDir":       "data/maf/loot_table/treasure",
			"debugFunctionDir":      "data/maf/function/debug/give",
			"minecraftTagDir":       "data/minecraft/tags/function",
		},
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Config{
		Port:                8787,
		ItemStatePath:       filepath.Join(root, "item-state.json"),
		GrimoireStatePath:   filepath.Join(root, "grimoire-state.json"),
		SkillStatePath:      filepath.Join(root, "skill-state.json"),
		EnemySkillStatePath: filepath.Join(root, "enemy-skill-state.json"),
		EnemyStatePath:      filepath.Join(root, "enemy-state.json"),
		TreasureStatePath:   filepath.Join(root, "treasure-state.json"),
		ExportSettingsPath:  settingsPath,
	}
	handler := NewHandler(cfg, Dependencies{
		ItemRepo:       store.NewItemStateRepository(cfg.ItemStatePath),
		GrimoireRepo:   store.NewGrimoireStateRepository(cfg.GrimoireStatePath),
		SkillRepo:      store.NewEntryStateRepository[skills.SkillEntry](cfg.SkillStatePath),
		EnemySkillRepo: store.NewEntryStateRepository[enemyskills.EnemySkillEntry](cfg.EnemySkillStatePath),
		EnemyRepo:      store.NewEntryStateRepository[enemies.EnemyEntry](cfg.EnemyStatePath),
		TreasureRepo:   store.NewEntryStateRepository[treasures.TreasureEntry](cfg.TreasureStatePath),
		Now:            fixedNow,
	})
	return handler, root
}

func fixedNow() time.Time {
	return time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC)
}

func postJSON(t *testing.T, handler http.Handler, path string, body any, wantStatus int) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	rec := request(t, handler, http.MethodPost, path, data)
	if rec.Code != wantStatus {
		t.Fatalf("%s status = %d, want %d, body=%s", path, rec.Code, wantStatus, rec.Body.String())
	}
	return rec
}

func request(t *testing.T, handler http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func requestForm(t *testing.T, handler http.Handler, method, path string, values url.Values) *httptest.ResponseRecorder {
	t.Helper()
	body := values.Encode()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func postForm(t *testing.T, handler http.Handler, path string, values url.Values, wantStatus int) *httptest.ResponseRecorder {
	t.Helper()
	rec := requestForm(t, handler, http.MethodPost, path, values)
	if rec.Code != wantStatus {
		t.Fatalf("%s status = %d, want %d, body=%s", path, rec.Code, wantStatus, rec.Body.String())
	}
	return rec
}

func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder, dest any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), dest); err != nil {
		t.Fatalf("decode %q: %v", rec.Body.String(), err)
	}
}

func assertEntriesLen(t *testing.T, handler http.Handler, path, field string, want int) {
	t.Helper()
	rec := request(t, handler, http.MethodGet, path, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("%s status = %d, want %d", path, rec.Code, http.StatusOK)
	}
	var body map[string]any
	decodeResponse(t, rec, &body)
	entries, ok := body[field].([]any)
	if !ok || len(entries) != want {
		t.Fatalf("%s body = %#v, want %s len %d", path, body, field, want)
	}
}

func deleteOK(t *testing.T, handler http.Handler, path string) {
	t.Helper()
	rec := request(t, handler, http.MethodDelete, path, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("%s status = %d, want %d, body=%s", path, rec.Code, http.StatusOK, rec.Body.String())
	}
}

func uuid(suffix string) string {
	return "00000000-0000-4000-8000-" + suffix
}
