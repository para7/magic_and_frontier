package export

import (
	"path/filepath"

	config "maf_command_editor/app/files"
)

func ExportDatapack(dmas DBMaster, mafconfig config.MafConfig) error {
	settings, err := config.LoadExportSettings(mafconfig.ExportSettingsPath)
	if err != nil {
		return err
	}

	effectRelDir := settings.ExportPaths.GrimoireEffect
	effectDir := filepath.Join(settings.OutputRoot, effectRelDir)
	effectSelect := filepath.Join(settings.OutputRoot, settings.ExportPaths.GrimoireSelectFile)

	effects := BuildGrimoireArtifacts(dmas, effectRelDir)
	return WriteGrimoireArtifacts(effectDir, effectSelect, effects)
}
