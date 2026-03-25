package treasures

import "tools2/app/internal/domain/entity"

type DropRef = entity.DropRef

type TreasureEntry struct {
	ID        string    `json:"id"`
	TablePath string    `json:"tablePath"`
	LootPools []DropRef `json:"lootPools"`
	UpdatedAt string    `json:"updatedAt"`
}

type SaveInput struct {
	ID        string    `json:"id"`
	TablePath string    `json:"tablePath" validate:"trimmed_required,trimmed_min=1,trimmed_max=200"`
	LootPools []DropRef `json:"lootPools" validate:"min=1,dive"`
}
