package treasure

import model "maf_command_editor/app/domain/model"

type Treasure struct {
	ID        string          `json:"id"        validate:"trimmed_required"`
	TablePath string          `json:"tablePath" validate:"trimmed_required"`
	LootPools []model.DropRef `json:"lootPools" validate:"min=1,dive"`
}
