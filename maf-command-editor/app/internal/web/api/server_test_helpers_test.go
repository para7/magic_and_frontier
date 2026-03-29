package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"maf-command-editor/app/internal/application"
	"maf-command-editor/app/internal/config"
	"maf-command-editor/app/internal/domain/common"
	"maf-command-editor/app/internal/domain/export"
)

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
		IDCounterStatePath:     filepath.Join(root, "id_counters.json"),
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
