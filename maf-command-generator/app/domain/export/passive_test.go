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

	effects, bows, grimoires, err := BuildPassiveArtifacts(master)
	if err != nil {
		t.Fatal(err)
	}
	if len(effects) != 1 {
		t.Fatalf("effects length = %d, want 1", len(effects))
	}
	if len(bows) != 0 {
		t.Fatalf("bows length = %d, want 0", len(bows))
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
	if !strings.Contains(grimoires[0].ApplyBody, `data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1.id set value "passive_1"`) {
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

	_, _, grimoires, err := BuildPassiveArtifacts(master)
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

func TestBuildPassiveArtifactsBuildsBowEffectAndBowBody(t *testing.T) {
	master := exportMasterStub{
		passives: []passiveModel.Passive{
			{
				ID:        "bow_passive",
				Name:      "Bow Passive",
				Condition: "bow",
				Slots:     []int{1},
				Script:    []string{"say bow hit", "effect give @s glowing 1 0 true"},
				Bow:       &passiveModel.BowConfig{LifeSub: ptrInt(50)},
			},
		},
	}

	effects, bows, grimoires, err := BuildPassiveArtifacts(master)
	if err != nil {
		t.Fatal(err)
	}
	if len(effects) != 1 {
		t.Fatalf("effects length = %d, want 1", len(effects))
	}
	if len(bows) != 1 {
		t.Fatalf("bows length = %d, want 1", len(bows))
	}
	if len(grimoires) != 1 {
		t.Fatalf("grimoires length = %d, want 1", len(grimoires))
	}
	if !strings.Contains(effects[0].Body, "mafBowUsed") {
		t.Fatalf("bow effect should check mafBowUsed: %q", effects[0].Body)
	}
	if !strings.Contains(effects[0].Body, `function maf:magic/passive/tag_passive_arrow {passive_id:"bow_passive",life:1150}`) {
		t.Fatalf("bow effect should tag arrow with passive id and life: %q", effects[0].Body)
	}
	if bows[0].ID != "bow_passive" {
		t.Fatalf("unexpected bow artifact id: %#v", bows[0])
	}
	if bows[0].Body != "say bow hit\neffect give @s glowing 1 0 true" {
		t.Fatalf("unexpected bow artifact body: %q", bows[0].Body)
	}
}

func TestWritePassiveArtifactsWritesFiles(t *testing.T) {
	root := t.TempDir()
	effectDir := filepath.Join(root, "effect")
	bowDir := filepath.Join(root, "bow")
	giveDir := filepath.Join(root, "give")
	applyDir := filepath.Join(root, "apply")

	effects := []PassiveEffectFunction{
		{ID: "passive_1", Body: "say effect"},
	}
	bows := []PassiveBowFunction{
		{ID: "passive_1", Body: "say bow"},
	}
	grimoires := []PassiveGrimoireFunction{
		{FunctionID: "passive_1_slot1", GiveBody: "say give", ApplyBody: "say set slot1"},
	}

	if err := WritePassiveArtifacts(effectDir, bowDir, giveDir, applyDir, effects, bows, grimoires); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(effectDir, "passive_1.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "say effect\n" {
		t.Fatalf("unexpected passive effect body: %q", string(body))
	}
	bowBody, err := os.ReadFile(filepath.Join(bowDir, "passive_1.mcfunction"))
	if err != nil {
		t.Fatal(err)
	}
	if string(bowBody) != "say bow\n" {
		t.Fatalf("unexpected passive bow body: %q", string(bowBody))
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
			"passiveBow":     "generated/passive/bow",
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
	passiveBowPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "passive", "bow", "passive_1.mcfunction")
	if _, err := os.Stat(passiveBowPath); !os.IsNotExist(err) {
		t.Fatalf("non-bow passive should not create bow file: %v", err)
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

	defaultBowPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "passive", "bow")
	if _, err := os.Stat(defaultBowPath); !os.IsNotExist(err) {
		t.Fatalf("default passive bow dir should not be created for non-bow passive: %v", err)
	}

	legacyGivePath := filepath.Join(root, "out", "data", "maf", "function", "legacy", "passive", "grimoire", "passive_1_slot1.mcfunction")
	if _, err := os.Stat(legacyGivePath); !os.IsNotExist(err) {
		t.Fatalf("legacy passiveGrimoire path should not be used: %v", err)
	}
}

func TestExportDatapackWritesPassiveBowArtifacts(t *testing.T) {
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
		passives: []passiveModel.Passive{
			{
				ID:        "bow_passive",
				Name:      "Bow Passive",
				Condition: "bow",
				Slots:     []int{1},
				Script:    []string{"say bow hit"},
				Bow:       &passiveModel.BowConfig{LifeSub: ptrInt(10)},
			},
		},
	}

	if err := ExportDatapack(master, cfg); err != nil {
		t.Fatal(err)
	}

	passiveEffectPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "passive", "effect", "bow_passive.mcfunction")
	effectBody, err := os.ReadFile(passiveEffectPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(effectBody), "mafBowUsed") {
		t.Fatalf("bow passive effect should contain bow trigger guard: %s", string(effectBody))
	}

	passiveBowPath := filepath.Join(root, "out", "data", "maf", "function", "generated", "passive", "bow", "bow_passive.mcfunction")
	bowBody, err := os.ReadFile(passiveBowPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(bowBody) != "say bow hit\n" {
		t.Fatalf("unexpected bow passive body: %q", string(bowBody))
	}
}
