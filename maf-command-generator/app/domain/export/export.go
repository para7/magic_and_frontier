package export

import (
	"path/filepath"

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
	effectSelect := filepath.Join(settings.OutputRoot, funcRoot, settings.ExportPaths.GrimoireSelectFile)
	debugDir := filepath.Join(settings.OutputRoot, funcRoot, settings.ExportPaths.GrimoireDebug)
	enemySkillLogicalDir := settings.ExportPaths.EnemySkill
	enemySkillDir := filepath.Join(settings.OutputRoot, funcRoot, enemySkillLogicalDir)
	enemyLogicalDir := settings.ExportPaths.Enemy
	enemyDir := filepath.Join(settings.OutputRoot, funcRoot, enemyLogicalDir)
	enemyLootLogicalDir := settings.ExportPaths.EnemyLoot
	enemyLootDir := filepath.Join(settings.OutputRoot, lootRoot, enemyLootLogicalDir)

	effects := BuildGrimoireArtifacts(dmas, effectLogicalDir)
	if err := WriteGrimoireArtifacts(effectDir, effectSelect, effects); err != nil {
		return err
	}
	if err := WriteGrimoireDebugArtifacts(debugDir, effects); err != nil {
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
