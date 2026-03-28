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
	spellEffectPath := settings.ExportPaths.SpellEffect
	writeDir := filepath.Join(settings.OutputRoot, spellEffectPath)
	artifacts := BuildGrimoireArtifacts(dmas, spellEffectPath)
	return WriteGrimoireArtifacts(writeDir, artifacts)
}
