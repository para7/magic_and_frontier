package export

import (
	"path/filepath"

	config "maf_command_editor/app/files"
)

const funcRoot = "data/maf/function"

func ExportDatapack(dmas DBMaster, mafconfig config.MafConfig) error {
	settings, err := config.LoadExportSettings(mafconfig.ExportSettingsPath)
	if err != nil {
		return err
	}

	effectLogicalDir := settings.ExportPaths.GrimoireEffect
	effectDir := filepath.Join(settings.OutputRoot, funcRoot, effectLogicalDir)
	effectSelect := filepath.Join(settings.OutputRoot, funcRoot, settings.ExportPaths.GrimoireSelectFile)
	debugDir := filepath.Join(settings.OutputRoot, funcRoot, settings.ExportPaths.GrimoireDebug)

	effects := BuildGrimoireArtifacts(dmas, effectLogicalDir)
	if err := WriteGrimoireArtifacts(effectDir, effectSelect, effects); err != nil {
		return err
	}
	if err := WriteGrimoireDebugArtifacts(debugDir, effects); err != nil {
		return err
	}
	return nil
}
