package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	grimoireModel "maf_command_editor/app/domain/model/grimoire"
	passiveModel "maf_command_editor/app/domain/model/passive"
	config "maf_command_editor/app/files"
)

func TestBuildPassiveArtifactsBuildsEffectAndSlotGrimoire(t *testing.T) {
	master := exportMasterStub{
		passives: []passiveModel.Passive{
			{
				ID:        "passive_1",
				Name:      "Quickstep",
				Condition: "always",
				Slots:     []int{2, 1},
				Script:    []string{"say passive"},
			},
		},
	}

	effects, grimoires, err := BuildPassiveArtifacts(master)
	if err != nil {
		t.Fatal(err)
	}
	if len(effects) != 1 {
		t.Fatalf("effects length = %d, want 1", len(effects))
	}
	if effects[0].ID != "passive_1" || effects[0].Body != "say passive" {
		t.Fatalf("unexpected effects[0]: %#v", effects[0])
	}
	if len(grimoires) != 2 {
		t.Fatalf("grimoires length = %d, want 2", len(grimoires))
	}
	if grimoires[0].FunctionID != "passive_1_slot1" {
		t.Fatalf("unexpected first slot artifact: %#v", grimoires[0])
	}
	if grimoires[1].FunctionID != "passive_1_slot2" {
		t.Fatalf("unexpected second slot artifact: %#v", grimoires[1])
	}
	if !strings.Contains(grimoires[0].GiveBody, `give @p minecraft:book[`) {
		t.Fatalf("unexpected slot1 give body: %q", grimoires[0].GiveBody)
	}
	if !strings.Contains(grimoires[0].ApplyBody, `data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].passive.slot1.id set value "passive_1"`) {
		t.Fatalf("unexpected slot1 apply id body: %q", grimoires[0].ApplyBody)
	}
	if !strings.Contains(grimoires[0].ApplyBody, "function #oh_my_dat:please") {
		t.Fatalf("unexpected slot1 apply function call: %q", grimoires[0].ApplyBody)
	}
	if !strings.Contains(grimoires[0].ApplyBody, `tellraw @s [{"text":"[slot1]に[Quickstep]を設定しました"}]`) {
		t.Fatalf("unexpected slot1 apply message: %q", grimoires[0].ApplyBody)
	}
}

func TestBuildPassiveArtifactsUsesIDWhenNameIsBlank(t *testing.T) {
	master := exportMasterStub{
		passives: []passiveModel.Passive{
			{
				ID:        "passive_1",
				Name:      "   ",
				Condition: "always",
				Slots:     []int{1},
				Script:    []string{"say passive"},
			},
		},
	}

	_, grimoires, err := BuildPassiveArtifacts(master)
	if err != nil {
		t.Fatal(err)
	}
	if len(grimoires) != 1 {
		t.Fatalf("grimoires length = %d, want 1", len(grimoires))
	}
	if !strings.Contains(grimoires[0].ApplyBody, `tellraw @s [{"text":"[slot1]に[passive_1]を設定しました"}]`) {
		t.Fatalf("unexpected apply message fallback: %q", grimoires[0].ApplyBody)
	}
}

func TestWritePassiveArtifactsWritesFiles(t *testing.T) {
	root := t.TempDir()
	effectDir := filepath.Join(root, "effect")
	giveDir := filepath.Join(root, "give")
	applyDir := filepath.Join(root, "apply")

	effects := []PassiveEffectFunction{
		{ID: "passive_1", Body: "say effect"},
	}
	grimoires := []PassiveGrimoireFunction{
		{FunctionID: "passive_1_slot1", GiveBody: "say give", ApplyBody: "say set slot1"},
	}

	if err := WritePassiveArtifacts(effectDir, giveDir, applyDir, effects, grimoires); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(effectDir, "passive_1.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "say effect\n" {
		t.Fatalf("unexpected passive effect body: %q", string(body))
	}
	setBody, err := os.ReadFile(filepath.Join(giveDir, "passive_1_slot1.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if string(setBody) != "say give\n" {
		t.Fatalf("unexpected passive grimoire give body: %q", string(setBody))
	}
	applyBody, err := os.ReadFile(filepath.Join(applyDir, "passive_1_slot1.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if string(applyBody) != "say set slot1\n" {
		t.Fatalf("unexpected passive grimoire apply body: %q", string(applyBody))
	}
}

func TestExportDatapackWritesPassiveArtifacts(t *testing.T) {
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	settings := map[string]any{
		"outputRoot": filepath.Join(root, "out"),
		"exportPaths": map[string]any{
			"grimoireEffect": "generated/grimoire/effect",
			"grimoireDebug":  "generated/grimoire/give",
			"passiveEffect":  "generated/passive/effect",
			"passiveGive":    "generated/passive/give",
			"passiveApply":   "generated/passive/apply",
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
		grimoires: []grimoireModel.Grimoire{
			{ID: "fire", Script: []string{"say fire"}, Title: "Fire"},
		},
		passives: []passiveModel.Passive{
			{ID: "passive_1", Name: "Quickstep", Condition: "always", Slots: []int{1}, Script: []string{"say passive"}},
		},
	}

	if err := ExportDatapack(master, cfg); err != nil {
		t.Fatal(err)
	}

	passiveEffectPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "passive", "effect", "passive_1.mcfunction")
	if _, err := os.Stat(passiveEffectPath); err != nil {
		t.Fatalf("missing passive effect file: %v", err)
	}
	passiveGrimoirePath := filepath.Join(root, "out", "data", "maf", "function", "generated", "passive", "give", "passive_1_slot1.mcfunction")
	if _, err := os.Stat(passiveGrimoirePath); err != nil {
		t.Fatalf("missing passive grimoire file: %v", err)
	}
	passiveApplyPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "passive", "apply", "passive_1_slot1.mcfunction")
	if _, err := os.Stat(passiveApplyPath); err != nil {
		t.Fatalf("missing passive apply file: %v", err)
	}
	applyBody, err := os.ReadFile(passiveApplyPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(applyBody), "function #oh_my_dat:please") {
		t.Fatalf("passive apply should contain oh_my_dat function call: %s", string(applyBody))
	}
	if !strings.Contains(string(applyBody), `tellraw @s [{"text":"[slot1]に[Quickstep]を設定しました"}]`) {
		t.Fatalf("passive apply should contain success message: %s", string(applyBody))
	}

	grimoireEffectPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "grimoire", "effect", "fire.mcfunction")
	if _, err := os.Stat(grimoireEffectPath); err != nil {
		t.Fatalf("missing grimoire effect file: %v", err)
	}
	selectExecPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "grimoire", "selectexec.mcfunction")
	if _, err := os.Stat(selectExecPath); !os.IsNotExist(err) {
		t.Fatalf("legacy selectexec should not be created: %v", err)
	}
	grimoireSetupPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "grimoire", "setup_effect_ref_map.mcfunction")
	if _, err := os.Stat(grimoireSetupPath); !os.IsNotExist(err) {
		t.Fatalf("legacy setup map should not be created: %v", err)
	}
}

func TestExportDatapackUsesDefaultPassiveGivePathAndIgnoresLegacyField(t *testing.T) {
	root := t.TempDir()
	settingsPath := filepath.Join(root, "export_settings.json")
	settings := map[string]any{
		"outputRoot": filepath.Join(root, "out"),
		"exportPaths": map[string]any{
			"grimoireEffect":  "generated/grimoire/effect",
			"grimoireDebug":   "generated/grimoire/give",
			"passiveEffect":   "generated/passive/effect",
			"passiveApply":    "generated/passive/apply",
			"passiveGrimoire": "legacy/passive/grimoire",
			"enemy":           "generated/enemy/spawn",
			"enemySkill":      "generated/enemy/skill",
			"enemyLoot":       "generated/enemy/loot",
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
		passives: []passiveModel.Passive{
			{ID: "passive_1", Name: "Quickstep", Condition: "always", Slots: []int{1}, Script: []string{"say passive"}},
		},
	}

	if err := ExportDatapack(master, cfg); err != nil {
		t.Fatal(err)
	}

	defaultGivePath := filepath.Join(root, "out", "data", "maf", "function", "generated", "passive", "give", "passive_1_slot1.mcfunction")
	if _, err := os.Stat(defaultGivePath); err != nil {
		t.Fatalf("default passive give file should exist: %v", err)
	}

	legacyGivePath := filepath.Join(root, "out", "data", "maf", "function", "legacy", "passive", "grimoire", "passive_1_slot1.mcfunction")
	if _, err := os.Stat(legacyGivePath); !os.IsNotExist(err) {
		t.Fatalf("legacy passiveGrimoire path should not be used: %v", err)
	}
}
