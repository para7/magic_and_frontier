package export

import (
	"path/filepath"
	"strings"

	config "maf_command_editor/app/files"
)

const funcRoot = "data/maf/function"
const lootRoot = "data/maf/loot_table"

func ExportDatapack(dmas DBMaster, mafconfig config.MafConfig) error {
	settings, err := config.LoadExportSettings(mafconfig.ExportSettingsPath)
	if err != nil {
		return err
	}

	effectLogicalDir := settings.ExportPaths.GrimoireEffect
	effectDir := filepath.Join(settings.OutputRoot, funcRoot, effectLogicalDir)
	debugDir := filepath.Join(settings.OutputRoot, funcRoot, settings.ExportPaths.GrimoireDebug)
	passiveEffectLogicalDir := normalizePathOrDefault(settings.ExportPaths.PassiveEffect, "generated/passive/effect")
	passiveEffectDir := filepath.Join(settings.OutputRoot, funcRoot, passiveEffectLogicalDir)
	passiveGiveLogicalDir := normalizePathOrDefault(settings.ExportPaths.PassiveGive, "generated/passive/give")
	passiveGiveDir := filepath.Join(settings.OutputRoot, funcRoot, passiveGiveLogicalDir)
	passiveApplyLogicalDir := normalizePathOrDefault(settings.ExportPaths.PassiveApply, "generated/passive/apply")
	passiveApplyDir := filepath.Join(settings.OutputRoot, funcRoot, passiveApplyLogicalDir)
	enemySkillLogicalDir := settings.ExportPaths.EnemySkill
	enemySkillDir := filepath.Join(settings.OutputRoot, funcRoot, enemySkillLogicalDir)
	enemyLogicalDir := settings.ExportPaths.Enemy
	enemyDir := filepath.Join(settings.OutputRoot, funcRoot, enemyLogicalDir)
	enemyLootLogicalDir := settings.ExportPaths.EnemyLoot
	enemyLootDir := filepath.Join(settings.OutputRoot, lootRoot, enemyLootLogicalDir)

	effects := BuildGrimoireArtifacts(dmas)
	passiveEffects, passiveGrimoires, err := BuildPassiveArtifacts(dmas)
	if err != nil {
		return err
	}
	if err := WriteGrimoireArtifacts(effectDir, effects); err != nil {
		return err
	}
	if err := WriteGrimoireDebugArtifacts(debugDir, effects); err != nil {
		return err
	}
	grimoireDir := filepath.Dir(effectDir)
	if err := removeFileIfExists(filepath.Join(grimoireDir, "selectexec.mcfunction")); err != nil {
		return err
	}
	if err := removeFileIfExists(filepath.Join(grimoireDir, "setup_effect_ref_map.mcfunction")); err != nil {
		return err
	}
	if err := WritePassiveArtifacts(passiveEffectDir, passiveGiveDir, passiveApplyDir, passiveEffects, passiveGrimoires); err != nil {
		return err
	}
	skills := BuildEnemySkillArtifacts(dmas, enemySkillLogicalDir)
	if err := WriteEnemySkillArtifacts(enemySkillDir, skills); err != nil {
		return err
	}

	enemies, err := BuildEnemyArtifacts(dmas, enemyLootLogicalDir, mafconfig.MinecraftLootTableRoot)
	if err != nil {
		return err
	}
	if err := WriteEnemyArtifacts(enemyDir, enemyLootDir, enemies); err != nil {
		return err
	}
	return nil
}

func normalizePathOrDefault(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
