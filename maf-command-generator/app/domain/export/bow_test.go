package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	bowModel "maf_command_editor/app/domain/model/bow"
	config "maf_command_editor/app/files"
)

func TestBuildBowArtifactsBuildsAllOutputs(t *testing.T) {
	master := exportMasterStub{
		bows: []bowModel.BowPassive{
			{
				ID:           "test_full",
				LifeSub:      ptrInt(100),
				ScriptHit:    []string{"say hit"},
				ScriptFired:  []string{"say fired"},
				ScriptFlying: []string{"say flying"},
				ScriptGround: []string{"say ground"},
			},
		},
	}

	effects, hits, flyings, grounds, err := BuildBowArtifacts(master)
	if err != nil {
		t.Fatal(err)
	}
	if len(effects) != 1 || len(hits) != 1 || len(flyings) != 1 || len(grounds) != 1 {
		t.Fatalf("unexpected artifact counts: %d %d %d %d", len(effects), len(hits), len(flyings), len(grounds))
	}
	if effects[0].ID != "bow_test_full" {
		t.Fatalf("unexpected effect id: %#v", effects[0])
	}
	if !strings.Contains(effects[0].Body, `function maf:magic/bow/tag_bow_arrow {bow_id:"test_full",life:1100}`) {
		t.Fatalf("effect should tag bow arrow: %q", effects[0].Body)
	}
	if !strings.Contains(effects[0].Body, "function maf:magic/bow/prepare_hit_arrow") {
		t.Fatalf("effect should prepare hit arrow: %q", effects[0].Body)
	}
	if !strings.Contains(effects[0].Body, "tag @s add flying") {
		t.Fatalf("effect should add flying tag: %q", effects[0].Body)
	}
	if !strings.Contains(effects[0].Body, "tag @s add ground") {
		t.Fatalf("effect should add ground tag: %q", effects[0].Body)
	}
	if !strings.Contains(effects[0].Body, "say fired") {
		t.Fatalf("effect should append fired script: %q", effects[0].Body)
	}
	if hits[0].Body != "say hit" {
		t.Fatalf("unexpected hit body: %q", hits[0].Body)
	}
	if flyings[0].Body != "say flying" {
		t.Fatalf("unexpected flying body: %q", flyings[0].Body)
	}
	if grounds[0].Body != "say ground" {
		t.Fatalf("unexpected ground body: %q", grounds[0].Body)
	}
}

func TestWriteBowArtifactsWritesFiles(t *testing.T) {
	root := t.TempDir()

	if err := WriteBowArtifacts(
		filepath.Join(root, "effect"),
		filepath.Join(root, "bow"),
		filepath.Join(root, "flying"),
		filepath.Join(root, "ground"),
		[]bowModel.BowPassive{{ID: "test", ScriptHit: []string{"say hit"}, ScriptFlying: []string{"say flying"}, ScriptGround: []string{"say ground"}}},
		[]BowEffectFunction{{ID: "bow_test", Body: "say effect"}},
		[]BowHitFunction{{ID: "test", Body: "say hit"}},
		[]BowFlyingFunction{{ID: "test", Body: "say flying"}},
		[]BowGroundFunction{{ID: "test", Body: "say ground"}},
	); err != nil {
		t.Fatal(err)
	}

	assertBody := func(path, want string) {
		t.Helper()
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != want+"\n" {
			t.Fatalf("unexpected body for %s: %q", path, string(body))
		}
	}

	assertBody(filepath.Join(root, "effect", "bow_test.mcfunction"), "say effect")
	assertBody(filepath.Join(root, "bow", "test.mcfunction"), "say hit")
	assertBody(filepath.Join(root, "flying", "test_flying.mcfunction"), "say flying")
	assertBody(filepath.Join(root, "ground", "test_ground.mcfunction"), "say ground")
}

func TestBuildBowArtifactsSkipsEmptyOptionalScripts(t *testing.T) {
	master := exportMasterStub{
		bows: []bowModel.BowPassive{
			{
				ID:      "test_hitless",
				LifeSub: ptrInt(100),
			},
		},
	}

	effects, hits, flyings, grounds, err := BuildBowArtifacts(master)
	if err != nil {
		t.Fatal(err)
	}
	if len(effects) != 1 {
		t.Fatalf("unexpected effect count: %d", len(effects))
	}
	if len(hits) != 0 || len(flyings) != 0 || len(grounds) != 0 {
		t.Fatalf("optional outputs should be skipped when scripts are empty: %d %d %d", len(hits), len(flyings), len(grounds))
	}
	if strings.Contains(effects[0].Body, "function maf:magic/bow/prepare_hit_arrow") {
		t.Fatalf("effect should not prepare hit arrow when script_hit is empty: %q", effects[0].Body)
	}
	if strings.Contains(effects[0].Body, "tag @s add flying") {
		t.Fatalf("effect should not add flying tag when script_flying is empty: %q", effects[0].Body)
	}
	if strings.Contains(effects[0].Body, "tag @s add ground") {
		t.Fatalf("effect should not add ground tag when script_ground is empty: %q", effects[0].Body)
	}
}

func TestWriteBowArtifactsRemovesFilesForEmptyOptionalScripts(t *testing.T) {
	root := t.TempDir()
	bowDir := filepath.Join(root, "bow")
	flyingDir := filepath.Join(root, "flying")
	groundDir := filepath.Join(root, "ground")

	if err := writeFunctionFile(filepath.Join(bowDir, "test.mcfunction"), ""); err != nil {
		t.Fatal(err)
	}
	if err := writeFunctionFile(filepath.Join(flyingDir, "test_flying.mcfunction"), ""); err != nil {
		t.Fatal(err)
	}
	if err := writeFunctionFile(filepath.Join(groundDir, "test_ground.mcfunction"), ""); err != nil {
		t.Fatal(err)
	}

	if err := WriteBowArtifacts(
		filepath.Join(root, "effect"),
		bowDir,
		flyingDir,
		groundDir,
		[]bowModel.BowPassive{{ID: "test"}},
		[]BowEffectFunction{{ID: "bow_test", Body: "say effect"}},
		nil,
		nil,
		nil,
	); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{
		filepath.Join(bowDir, "test.mcfunction"),
		filepath.Join(flyingDir, "test_flying.mcfunction"),
		filepath.Join(groundDir, "test_ground.mcfunction"),
	} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be removed, got err=%v", path, err)
		}
	}
}

func TestExportDatapackWritesBowArtifacts(t *testing.T) {
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	settings := map[string]any{
		"outputRoot": filepath.Join(root, "out"),
		"exportPaths": map[string]any{
			"grimoireEffect": "generated/grimoire/effect",
			"grimoireDebug":  "generated/grimoire/give",
			"passiveEffect":  "generated/passive/effect",
			"passiveBow":     "generated/passive/bow",
			"passiveGive":    "generated/passive/give",
			"passiveApply":   "generated/passive/apply",
			"bowFlying":      "generated/bow/flying",
			"bowGround":      "generated/bow/ground",
			"enemy":          "generated/enemy/spawn",
			"enemySkill":     "generated/enemy/skill",
			"enemyLoot":      "generated/enemy/loot",
		},
	}
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := config.LoadConfig()
	cfg.ExportSettingsPath = settingsPath
	cfg.MinecraftLootTableRoot = filepath.Join(root, "minecraft", "loot_table")

	master := exportMasterStub{
		bows: []bowModel.BowPassive{
			{
				ID:           "test_full",
				LifeSub:      ptrInt(100),
				ScriptHit:    []string{"say hit"},
				ScriptFired:  []string{"say fired"},
				ScriptFlying: []string{"say flying"},
				ScriptGround: []string{"say ground"},
			},
		},
	}

	if err := ExportDatapack(master, cfg); err != nil {
		t.Fatal(err)
	}

	checkContains := func(path, want string) {
		t.Helper()
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), want) {
			t.Fatalf("%s should contain %q: %s", path, want, string(body))
		}
	}

	checkContains(filepath.Join(root, "out", "data", "maf", "function", "generated", "passive", "effect", "bow_test_full.mcfunction"), "say fired")
	checkContains(filepath.Join(root, "out", "data", "maf", "function", "generated", "passive", "bow", "test_full.mcfunction"), "say hit")
	checkContains(filepath.Join(root, "out", "data", "maf", "function", "generated", "bow", "flying", "test_full_flying.mcfunction"), "say flying")
	checkContains(filepath.Join(root, "out", "data", "maf", "function", "generated", "bow", "ground", "test_full_ground.mcfunction"), "say ground")
}
