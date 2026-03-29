package export

import (
	dmaster "maf-command-editor/app/internal/domain/master"
)

type MasterExportParams struct {
	ExportSettingsPath     string
	MinecraftLootTableRoot string
}

func ExportDatapackFromMaster(master dmaster.DBMaster, params MasterExportParams) SaveDataResponse {
	if err := ValidateSettings(params.ExportSettingsPath); err != nil {
		return SaveDataResponse{
			OK:      false,
			Code:    "INVALID_CONFIG",
			Message: "Invalid export settings.",
			Details: err.Error(),
		}
	}
	validation := master.ValidateSavedAll()
	if !validation.OK {
		return SaveDataResponse{
			OK:      false,
			Code:    "VALIDATION_FAILED",
			Message: "Savedata validation failed.",
			Details: validation.String(),
		}
	}
	return ExportDatapack(ExportParams{
		Items:                  master.Items().ListAll(),
		Grimoires:              master.Grimoires().ListAll(),
		Skills:                 master.Skills().ListAll(),
		EnemySkills:            master.EnemySkills().ListAll(),
		Enemies:                master.Enemies().ListAll(),
		SpawnTables:            master.SpawnTables().ListAll(),
		Treasures:              master.Treasures().ListAll(),
		LootTables:             master.LootTables().ListAll(),
		ExportSettingsPath:     params.ExportSettingsPath,
		MinecraftLootTableRoot: params.MinecraftLootTableRoot,
	})
}
