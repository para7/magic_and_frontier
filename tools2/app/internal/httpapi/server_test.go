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
	grimoireID := createJSONEntry(t, handler, "/api/grimoire", map[string]any{
		"title":    "Firebolt",
		"script":   "say fire",
		"castTime": 20,
		"mpCost":   5,
	})

	rec := request(t, handler, http.MethodGet, "/grimoire/edit?id="+grimoireID, nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `name="castid" value="1" readonly`) {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestHandlerAPIHappyPathAndSave(t *testing.T) {
	handler, root := newTestHandler(t)

	skillID := createJSONEntry(t, handler, "/api/skills", map[string]any{
		"name":        "Slash",
		"description": "Basic slash",
		"script":      "say slash",
	})

	itemID := createJSONEntry(t, handler, "/api/items", map[string]any{
		"itemId":  "minecraft:apple",
		"count":   1,
		"skillId": skillID,
	})

	grimoireID := createJSONEntry(t, handler, "/api/grimoire", map[string]any{
		"title":       "Firebolt",
		"description": "Burn",
		"script":      "say fire",
		"castTime":    20,
		"mpCost":      5,
	})

	enemySkillID := createJSONEntry(t, handler, "/api/enemy-skills", map[string]any{
		"name":        "Roar",
		"description": "Loud",
		"script":      "say roar",
	})

	createJSONEntry(t, handler, "/api/treasures", map[string]any{
		"mode":      "custom",
		"tablePath": "maf:treasure/test",
		"lootPools": []map[string]any{
			{"kind": "item", "refId": itemID, "weight": 1},
			{"kind": "grimoire", "refId": grimoireID, "weight": 1},
		},
	})

	createJSONEntry(t, handler, "/api/enemies", map[string]any{
		"mobType":       "minecraft:zombie",
		"name":          "Sample Zombie",
		"hp":            20,
		"dropMode":      "replace",
		"enemySkillIds": []string{enemySkillID},
		"drops": []map[string]any{
			{"kind": "minecraft_item", "refId": "minecraft:rotten_flesh", "weight": 1},
		},
		"equipment": map[string]any{
			"mainhand": map[string]any{
				"kind":  "minecraft_item",
				"refId": "minecraft:iron_sword",
				"count": 1,
			},
		},
	})

	rec := requestJSON(t, handler, http.MethodPost, "/api/save", map[string]any{})
	if rec.Code != http.StatusOK {
		t.Fatalf("save status = %d body=%s", rec.Code, rec.Body.String())
	}

	checkFiles := []string{
		filepath.Join(root, "out", "data", "maf", "function", "skill", skillID+".mcfunction"),
		filepath.Join(root, "out", "data", "maf", "function", "grimoire", grimoireID+".mcfunction"),
		filepath.Join(root, "out", "data", "maf", "function", "grimoire", "selectexec.mcfunction"),
		filepath.Join(root, "out", "data", "maf", "loot_table", "treasure", "test.json"),
		filepath.Join(root, "out", "data", "maf", "loot_table", "enemy", "enemy_1.json"),
	}
	for _, path := range checkFiles {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing exported file %s: %v", path, err)
		}
	}
}

func TestHandlerAPITreasureRejectsDuplicateCustomPath(t *testing.T) {
	handler, _ := newTestHandler(t)

	createJSONEntry(t, handler, "/api/treasures", map[string]any{
		"mode":      "custom",
		"tablePath": "maf:treasure/test",
		"lootPools": []map[string]any{{"kind": "minecraft_item", "refId": "minecraft:apple", "weight": 1}},
	})

	rec := requestJSON(t, handler, http.MethodPost, "/api/treasures", map[string]any{
		"mode":      "custom",
		"tablePath": "maf:treasure/test",
		"lootPools": []map[string]any{{"kind": "minecraft_item", "refId": "minecraft:stick", "weight": 1}},
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	decodeResponse(t, rec, &body)
	fieldErrors := body["fieldErrors"].(map[string]any)
	if fieldErrors["tablePath"] == nil {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandlerAPITreasureRejectsTraversalPath(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := requestJSON(t, handler, http.MethodPost, "/api/treasures", map[string]any{
		"mode":      "custom",
		"tablePath": "maf:loot/../escape",
		"lootPools": []map[string]any{{"kind": "minecraft_item", "refId": "minecraft:apple", "weight": 1}},
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandlerAPIGrimoireUsesServerManagedIdentity(t *testing.T) {
	handler, _ := newTestHandler(t)

	rec := requestJSON(t, handler, http.MethodPost, "/api/grimoire", map[string]any{
		"id":       "grimoire_999",
		"castid":   999,
		"title":    "Firebolt",
		"script":   "say fire",
		"castTime": 20,
		"mpCost":   5,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	decodeResponse(t, rec, &body)
	entry := body["entry"].(map[string]any)
	if entry["id"] != "grimoire_1" {
		t.Fatalf("entry = %#v", entry)
	}
	if entry["castid"].(float64) != 1 {
		t.Fatalf("entry = %#v", entry)
	}

	rec = requestJSON(t, handler, http.MethodPost, "/api/grimoire", map[string]any{
		"id":       "grimoire_1",
		"castid":   555,
		"title":    "Firebolt v2",
		"script":   "say fire2",
		"castTime": 30,
		"mpCost":   9,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	decodeResponse(t, rec, &body)
	entry = body["entry"].(map[string]any)
	if entry["castid"].(float64) != 1 {
		t.Fatalf("entry = %#v", entry)
	}
	if entry["title"] != "Firebolt v2" {
		t.Fatalf("entry = %#v", entry)
	}
}

func TestHandlerAPISkillDeleteRejectsReferencedItem(t *testing.T) {
	handler, _ := newTestHandler(t)

	skillID := createJSONEntry(t, handler, "/api/skills", map[string]any{
		"name":        "Slash",
		"description": "Basic slash",
		"script":      "say slash",
	})
	createJSONEntry(t, handler, "/api/items", map[string]any{
		"itemId":  "minecraft:apple",
		"count":   1,
		"skillId": skillID,
	})

	rec := request(t, handler, http.MethodDelete, "/api/skills/"+skillID, nil, "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	decodeResponse(t, rec, &body)
	if body["code"] != "REFERENCE_ERROR" {
		t.Fatalf("body = %#v", body)
	}
}

func createJSONEntry(t *testing.T, handler http.Handler, path string, payload map[string]any) string {
	t.Helper()
	rec := requestJSON(t, handler, http.MethodPost, path, payload)
	if rec.Code != http.StatusOK {
		t.Fatalf("%s status = %d body=%s", path, rec.Code, rec.Body.String())
	}
	var body map[string]any
	decodeResponse(t, rec, &body)
	entry := body["entry"].(map[string]any)
	return entry["id"].(string)
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
	writeJSONFile(t, settingsPath, settings)

	cfg := config.Config{
		Port:                8787,
		ItemStatePath:       filepath.Join(root, "item-state.json"),
		GrimoireStatePath:   filepath.Join(root, "grimoire-state.json"),
		SkillStatePath:      filepath.Join(root, "skill-state.json"),
		EnemySkillStatePath: filepath.Join(root, "enemy-skill-state.json"),
		EnemyStatePath:      filepath.Join(root, "enemy-state.json"),
		TreasureStatePath:   filepath.Join(root, "treasure-state.json"),
		IDCounterStatePath:  filepath.Join(root, "id-counters.json"),
		ExportSettingsPath:  settingsPath,
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

func requestJSON(t *testing.T, handler http.Handler, method, path string, payload any) *httptest.ResponseRecorder {
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

func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder, dest any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), dest); err != nil {
		t.Fatalf("decode response: %v body=%s", err, rec.Body.String())
	}
}

func writeJSONFile(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}
