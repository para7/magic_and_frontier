package loottable

import model "maf_command_editor/app/domain/model"

type LootTable struct {
	ID        string          `json:"id"        validate:"trimmed_required"`
	Memo      string          `json:"memo"      validate:"trimmed_max=400"`
	LootPools []model.DropRef `json:"lootPools" validate:"min=1,dive"`
	UpdatedAt string          `json:"updatedAt"`
}
